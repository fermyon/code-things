use std::collections::HashSet;

use anyhow::{Context, Result};
use jwt_simple::prelude::Duration;
use spin_sdk::{config, key_value};

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

impl Config {
    fn try_get_value(key: &str, store: &key_value::Store) -> Result<String> {
        // first try to get the value from key-value store
        store
            .get(key)
            .map(|b| String::from_utf8(b).unwrap())
            // then try to get the value from the config file
            .or_else(|_| config::get(key))
            // then try to get the value from the environment
            .or_else(|_| std::env::var(key))
            .context(format!(
                "Failed to get configuration value for key '{}'",
                key
            ))
    }

    pub(crate) fn try_get(store: &key_value::Store) -> Result<Self> {
        let db_url = Self::try_get_value(KEY_DB_URL, store)?;
        let auth_domain = Self::try_get_value(KEY_AUTH_DOMAIN, store)?;
        let auth_max_validity: Option<Duration> =
            Self::try_get_value(KEY_AUTH_MAX_VALIDITY_SECS, store)
                .and_then(|s| {
                    Ok(s.parse::<u64>().context(format!(
                        "Value provided for '{}' must parse into an integer",
                        KEY_AUTH_MAX_VALIDITY_SECS
                    ))?)
                })
                .map(Duration::from_secs)
                .ok();

        let auth_audiences = HashSet::from([Self::try_get_value(KEY_AUTH_AUDIENCE, store)?]);
        let auth_issuers = HashSet::from([format!("https://{}/", auth_domain)]);
        let jwks_url = format!("https://{}/.well-known/jwks.json", auth_domain);

        Ok(Self {
            db_url,
            auth_audiences,
            auth_issuers,
            auth_max_validity,
            jwks_url,
        })
    }
}
