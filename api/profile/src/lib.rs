mod config;
mod model;

use anyhow::{anyhow, Result};
use bytes::Bytes;
use config::Config;
use model::Profile;
use spin_sdk::{
    http::{Request, Response},
    http_component,
};

enum Api {
    Create(model::Profile),
    ReadByHandle(String),
    Update(model::Profile),
    Delete(model::Profile),
    BadRequest,
    NotFound,
    MethodNotAllowed,
}

#[http_component]
fn profile_api(req: Request) -> Result<Response> {
    let cfg = Config::default();

    match api_from_request(&req) {
        Api::BadRequest => bad_request(),
        Api::MethodNotAllowed => method_not_allowed(),
        Api::Create(model) => handle_create(&cfg.db_url, model),
        Api::Update(model) => handle_update(&cfg.db_url, model),
        Api::ReadByHandle(handle) => handle_read_by_handle(&cfg.db_url, handle),
        Api::Delete(handle) => handle_delete_by_handle(&cfg.db_url, handle),
        _ => not_found(),
    }
}

fn api_from_request(req: &Request) -> Api {
    let auth_header = req
        .headers()
        .get(http::header::AUTHORIZATION)
        .and_then(|x| x.to_str().ok())
        .unwrap_or("None");
    println!("AUTHORIZATION: {auth_header:?}");

    let method = req.method().to_owned();
    let profile = match method {
        http::Method::GET | http::Method::DELETE => Profile::from_path(&req.headers()),
        http::Method::PUT | http::Method::POST => {
            Profile::from_bytes(req.body().as_ref().unwrap_or(&Bytes::new()))
        }
        _ => Err(anyhow!("Unsupported Http Method")),
    };

    match (method, profile) {
        (http::Method::POST, Ok(profile)) => Api::Create(profile),
        (http::Method::POST, Err(_)) => Api::BadRequest,
        (http::Method::GET, Ok(profile)) => Api::ReadByHandle(profile.handle),
        (http::Method::GET, Err(_)) => Api::NotFound,
        (http::Method::PUT, Ok(profile)) => Api::Update(profile),
        (http::Method::PUT, Err(_)) => Api::BadRequest,
        (http::Method::DELETE, Ok(profile)) => Api::Delete(profile),
        (http::Method::DELETE, Err(_)) => Api::NotFound,
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

fn handle_read_by_handle(db_url: &str, handle: String) -> Result<Response> {
    match Profile::get_by_handle(handle.as_str(), &db_url) {
        Ok(model) => ok(serde_json::to_string(&model)?),
        Err(_) => not_found(),
    }
}

fn handle_update(db_url: &str, model: Profile) -> Result<Response> {
    model.update(&db_url)?;
    handle_read_by_handle(&db_url, model.handle)
}

fn handle_delete_by_handle(db_url: &str, model: Profile) -> Result<Response> {
    match model.delete(&db_url) {
        Ok(_) => no_content(),
        Err(_) => internal_server_error(String::from("Error while deleting profile")),
    }
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

fn bad_request() -> Result<Response> {
    quick_response(http::StatusCode::BAD_REQUEST)
}

fn not_found() -> Result<Response> {
    quick_response(http::StatusCode::NOT_FOUND)
}

fn no_content() -> Result<Response> {
    quick_response(http::StatusCode::NO_CONTENT)
}

fn quick_response(s: http::StatusCode) -> Result<Response> {
    Ok(http::Response::builder().status(s).body(None)?)
}
