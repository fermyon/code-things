mod auth;
mod config;
mod model;

use anyhow::{anyhow, Context, Result};
use bytes::Bytes;
use config::Config;
use jwt_simple::{
    claims::{JWTClaims, NoCustomClaims},
    prelude::VerificationOptions,
};
use model::Profile;
use spin_sdk::{
    http::{Request, Response},
    http_component,
    key_value::Store,
};

enum Api {
    Create(model::Profile),
    ReadById(String),
    Update(model::Profile),
    DeleteById(String),
    MethodNotAllowed,
    NotFound,
}

#[http_component]
fn profile_api(req: Request) -> Result<Response> {
    let store = spin_sdk::key_value::Store::open_default()?;
    let cfg = Config::try_get(&store)?;

    // parse the profile from the request
    let method = req.method();
    let profile = match parse_profile(method, &req) {
        Ok(profile) => profile,
        Err(e) => return bad_request(e),
    };

    // guard against unauthenticated requests
    let claims = match claims_from_request(&cfg, &req, &profile.id, &store) {
        Ok(claims) if claims.subject.is_some() => claims,
        Ok(_) => return forbidden("Token is missing 'sub'.".to_string()),
        Err(e) => return forbidden(e.to_string()),
    };

    // add the subject to the profile
    let profile = profile.with_id(claims.subject);

    // match api action to handler
    match api_from_profile(method, profile) {
        Api::Create(profile) => handle_create(&cfg.db_url, profile),
        Api::Update(profile) => handle_update(&cfg.db_url, profile),
        Api::ReadById(id) => handle_read_by_id(&cfg.db_url, id),
        Api::DeleteById(id) => handle_delete_by_id(&cfg.db_url, id),
        Api::MethodNotAllowed => method_not_allowed(),
        Api::NotFound => not_found(),
    }
}

fn claims_from_request(
    cfg: &Config,
    req: &Request,
    subject: &Option<String>,
    store: &Store,
) -> Result<JWTClaims<NoCustomClaims>> {
    let keys = auth::JsonWebKeySet::get(cfg.jwks_url.to_owned(), store)
        .context(format!("Failed to retrieve JWKS from {:?}", cfg.jwks_url))?;

    let token = get_access_token(req.headers()).ok_or(anyhow!(
        "Failed to get access token from Authorization header"
    ))?;

    let options = VerificationOptions {
        max_validity: cfg.auth_max_validity,
        allowed_audiences: Some(cfg.auth_audiences.to_owned()),
        allowed_issuers: Some(cfg.auth_issuers.to_owned()),
        required_subject: subject.to_owned(),
        ..Default::default()
    };

    println!("[DEBUG] {:#?}", options);

    let claims = keys
        .verify(token, Some(options))
        .context("Failed to verify token")?;

    Ok(claims)
}

fn parse_profile(method: &http::Method, req: &Request) -> Result<Profile> {
    // parse the data model from body or url
    let profile = match method {
        &http::Method::GET | &http::Method::DELETE => Profile::from_path(&req.headers()),
        &http::Method::PUT | &http::Method::POST => {
            Profile::from_bytes(req.body().as_ref().unwrap_or(&Bytes::new()))
        }
        _ => Err(anyhow!("Unsupported Http Method")),
    }?;
    Ok(profile)
}

fn api_from_profile(method: &http::Method, profile: Profile) -> Api {
    match (method, profile) {
        (&http::Method::POST, profile) => Api::Create(profile),
        (&http::Method::GET, profile) if profile.id.is_some() => Api::ReadById(profile.id.unwrap()),
        (&http::Method::GET, _) => Api::NotFound,
        (&http::Method::PUT, profile) => Api::Update(profile),
        (&http::Method::DELETE, profile) if profile.id.is_some() => {
            Api::DeleteById(profile.id.unwrap())
        }
        (&http::Method::DELETE, _) => Api::NotFound,
        _ => Api::MethodNotAllowed,
    }
}

fn handle_create(db_url: &str, model: Profile) -> Result<Response> {
    model.insert(db_url)?;
    Ok(http::Response::builder()
        .status(http::StatusCode::CREATED)
        .header(
            http::header::LOCATION,
            format!("/api/profile/{}", model.handle),
        )
        .body(None)?)
}

fn handle_read_by_id(db_url: &str, id: String) -> Result<Response> {
    match Profile::get_by_id(id.as_str(), &db_url) {
        Ok(model) => ok(serde_json::to_string(&model)?),
        Err(_) => not_found(),
    }
}

fn handle_update(db_url: &str, model: Profile) -> Result<Response> {
    model.update(&db_url)?;
    handle_read_by_id(&db_url, model.handle)
}

fn handle_delete_by_id(db_url: &str, id: String) -> Result<Response> {
    match Profile::delete_by_id(&id, &db_url) {
        Ok(_) => no_content(),
        Err(_) => internal_server_error(String::from("Error while deleting profile")),
    }
}

fn get_access_token(headers: &http::HeaderMap) -> Option<&str> {
    headers
        .get("Authorization")?
        .to_str()
        .unwrap()
        .strip_prefix("Bearer ")
}

fn internal_server_error(err: String) -> Result<Response> {
    Ok(http::Response::builder()
        .status(http::StatusCode::INTERNAL_SERVER_ERROR)
        .header(http::header::CONTENT_TYPE, "text/plain")
        .body(Some(err.into()))?)
}

fn ok(payload: String) -> Result<Response> {
    Ok(http::Response::builder()
        .status(http::StatusCode::OK)
        .header(http::header::CONTENT_TYPE, "application/json")
        .body(Some(payload.into()))?)
}

fn method_not_allowed() -> Result<Response> {
    quick_response(http::StatusCode::METHOD_NOT_ALLOWED)
}

fn bad_request(err: anyhow::Error) -> Result<Response> {
    Ok(http::Response::builder()
        .status(http::StatusCode::BAD_REQUEST)
        .body(Some(err.to_string().into()))?)
}

fn not_found() -> Result<Response> {
    quick_response(http::StatusCode::NOT_FOUND)
}

fn no_content() -> Result<Response> {
    quick_response(http::StatusCode::NO_CONTENT)
}

fn forbidden(reason: String) -> Result<Response> {
    Ok(http::Response::builder()
        .status(http::StatusCode::FORBIDDEN)
        .header(http::header::CONTENT_TYPE, "text/plain")
        .body(Some(reason.into()))?)
}

fn quick_response(s: http::StatusCode) -> Result<Response> {
    Ok(http::Response::builder().status(s).body(None)?)
}
