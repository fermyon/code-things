use anyhow::{anyhow, Result, Ok};
use bytes::Bytes;
use http::HeaderMap;
use serde::{Deserialize, Serialize};

use spin_sdk::mysql::{self, ParameterValue, Decode};

use crate::utils::get_last_param_from_route;

fn as_param<'a>(value: &'a Option<String>) -> Option<ParameterValue<'a>> {
    match value {
        Some(value) => Some(ParameterValue::Str(value.as_str())),
        None => None
    }
}

fn as_nullable_param<'a>(value: &'a Option<String>) -> ParameterValue<'a> {
    match as_param(value) {
        Some(value) => value,
        None => ParameterValue::DbNull,
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub(crate) struct Profile {
    pub id: Option<String>,
    pub handle: String,
    pub avatar: Option<String>,
}

impl Profile {
    pub(crate) fn from_path(headers: &HeaderMap) -> Result<Self> {
        let header = headers.get("spin-path-info").ok_or(anyhow!("Error: Failed to discover path"))?;
        let path = header.to_str()?;
        match get_last_param_from_route(path) {
            Some(handle) => Ok(Profile {
                id: None,
                handle: handle,
                avatar: None,
            }),
            None => Err(anyhow!("Failed to parse handle from path")),
        }
    }

    pub(crate) fn from_bytes(b: &Bytes) -> Result<Self> {
        Ok(serde_json::from_slice(&b)?)
    }

    fn from_row(row: &spin_sdk::mysql::Row) -> Result<Self> {
        let id = String::decode(&row[0]).ok();
        let handle = String::decode(&row[1])?;
        let avatar = String::decode(&row[2]).ok();
        Ok(Profile {
            id,
            handle,
            avatar,
        })
    }

    pub(crate) fn insert(&self, db_url: &str) -> Result<()> {
        let params = vec![
            as_param(&self.id).ok_or(anyhow!("The id field is currently required for insert"))?,
            ParameterValue::Str(&self.handle),
            match as_param(&self.avatar) {
                Some(p) => p,
                None => ParameterValue::DbNull,
            }
        ];
        mysql::execute(db_url, "INSERT INTO profiles (id, handle, avatar) VALUES (?, ?, ?)", &params)?;
        Ok(())
    }

    pub(crate) fn get_by_handle(handle: &str, db_url: &str) -> Result<Profile> {
        let params = vec![ParameterValue::Str(handle)];
        let row_set = mysql::query(db_url, "SELECT id, handle, avatar from profiles WHERE handle = ?", &params)?;
        match row_set.rows.first() {
            Some(row) => Profile::from_row(row),
            None => Err(anyhow!("Profile not found for handle '{:?}'", handle))
        }
    }

    pub(crate) fn update(&self, db_url: &str) -> Result<()> {
        match &self.id {
            Some(id) => {
                let params = vec![
                    ParameterValue::Str(&self.handle),
                    as_nullable_param(&self.avatar),
                    ParameterValue::Str(id.as_str()),
                ];
                mysql::execute(db_url, "UPDATE profiles SET handle=?, avatar=? WHERE id=?", &params)?
            },
            None => {
                let params = vec![
                    as_nullable_param(&self.avatar),
                    ParameterValue::Str(self.handle.as_str())
                ];
                mysql::execute(db_url, "UPDATE profiles SET avatar=? WHERE handle=?", &params)?
            }
        }
        Ok(())
    }

    pub(crate) fn delete(&self, db_url: &str) -> Result<()> {
        match &self.id {
            Some(id) => {
                let params = vec![
                    ParameterValue::Str(id.as_str())
                ];
                mysql::execute(db_url, "DELETE FROM profiles WHERE id=?", &params)?
            },
            None => {
                let params = vec![
                    ParameterValue::Str(self.handle.as_str())
                ];
                mysql::execute(db_url, "DELETE FROM profiles WHERE handle=?", &params)?
            }
        }
        Ok(())
    }
}
