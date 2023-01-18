use anyhow::Result;
use spin_sdk::http::Response;

pub(crate) fn internal_server_error(err: String) -> Result<Response> {
    Ok(http::Response::builder()
        .status(http::StatusCode::INTERNAL_SERVER_ERROR)
        .header(http::header::CONTENT_TYPE, "text/plain")
        .body(Some(err.into()))?)
}

pub(crate) fn ok(payload: String) -> Result<Response> {
    Ok(http::Response::builder()
        .status(http::StatusCode::OK)
        .header(http::header::CONTENT_TYPE, "application/json")
        .body(Some(payload.into()))?)
}

pub(crate) fn method_not_allowed() -> Result<Response> {
    quick_response(http::StatusCode::METHOD_NOT_ALLOWED)
}

pub(crate) fn bad_request() -> Result<Response> {
    quick_response(http::StatusCode::BAD_REQUEST)
}

pub(crate) fn not_found() -> Result<Response> {
    quick_response(http::StatusCode::NOT_FOUND)
}

pub(crate) fn no_content() -> Result<Response> {
    quick_response(http::StatusCode::NO_CONTENT)
}

fn quick_response(s: http::StatusCode) -> Result<Response> {
    Ok(http::Response::builder().status(s).body(None)?)
}

pub(crate) fn get_params_from_route(route: &str) -> Vec<String> {
    route
        .split('/')
        .flat_map(|s| if s == "" { None } else { Some(s.to_string()) })
        .collect::<Vec<String>>()
}

pub(crate) fn get_last_param_from_route(route: &str) -> Option<String> {
    get_params_from_route(route).last().cloned()
}
