use std::collections::HashSet;

use jwt_simple::prelude::Duration;
use spin_sdk::config::get;

const KEY_DB_URL: &str = "db_url";
const KEY_AUTH_DOMAIN: &str = "auth_domain";
const KEY_AUTH_AUDIENCE: &str = "auth_audience";
const KEY_AUTH_MAX_VALIDITY_SECS: &str = "auth_max_validity_secs";

#[derive(Debug)]
pub(crate) struct Config {
    pub db_url: String,
    pub auth_audiences: HashSet<String>,
    pub auth_issuers: HashSet<String>,
    pub auth_max_validity: Option<Duration>,
    pub jwks_url: String,
}

impl Default for Config {
    fn default() -> Self {
        let db_url = get(KEY_DB_URL).expect("Missing config item 'db_url'");
        let auth_domain = get(KEY_AUTH_DOMAIN).expect("Missing config item 'auth_domain'");
        let auth_max_validity: Option<Duration> = get(KEY_AUTH_MAX_VALIDITY_SECS)
            .ok()
            .map(|s| {
                s.parse::<u64>()
                    .expect("Value provided must parse into an integer")
            })
            .map(Duration::from_secs);

        let auth_audiences = HashSet::from([
            get(KEY_AUTH_AUDIENCE).expect("Missing configuration item 'auth_audience'")
        ]);
        let auth_issuers = HashSet::from([format!("https://{0}/", auth_domain)]);
        let jwks_url = format!("https://{0}/.well-known/jwks.json", auth_domain);

        Self {
            db_url,
            auth_audiences,
            auth_issuers,
            auth_max_validity,
            jwks_url,
        }
    }
}
