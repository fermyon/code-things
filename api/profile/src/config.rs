use spin_sdk::config::get;

const KEY_DB_URL: &str = "db_url";

#[derive(Debug)]
pub(crate) struct Config {
    pub db_url: String,
}

impl Config {
    pub(crate) fn get() -> Config {
        Config {
            db_url: get(KEY_DB_URL).unwrap(),
        }
    }
}