use spin_sdk::config::get;

const KEY_DB_URL: &str = "db_url";

#[derive(Debug)]
pub(crate) struct Config {
    pub db_url: String,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            db_url: get(KEY_DB_URL).expect("Missing config item 'db_url'"),
        }
    }
}