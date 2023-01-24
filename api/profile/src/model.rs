use std::collections::HashMap;

use anyhow::{anyhow, Result, Ok};
use bytes::Bytes;
use http::HeaderMap;
use serde::{Deserialize, Serialize};

use spin_sdk::pg::{self as db, Decode, ParameterValue, Row};

use crate::utils::{get_last_param_from_route, get_column_lookup};

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

    fn from_row(row: &Row, columns: &HashMap<&str, usize>) -> Result<Self> {
        let id = String::decode(&row[columns["id"]]).ok();
        let handle = String::decode(&row[columns["handle"]])?;
        let avatar = String::decode(&row[columns["avatar"]]).ok();
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
        db::execute(db_url, "INSERT INTO profiles (id, handle, avatar) VALUES ($1, $2, $3)", &params)?;
        Ok(())
    }

    pub(crate) fn get_by_handle(handle: &str, db_url: &str) -> Result<Profile> {
        let params = vec![ParameterValue::Str(handle)];
        let row_set = db::query(db_url, "SELECT id, handle, avatar from profiles WHERE handle = $1", &params)?;

        let columns = get_column_lookup(&row_set.columns);

        match row_set.rows.first() {
            Some(row) => Profile::from_row(row, &columns),
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
                db::execute(db_url, "UPDATE profiles SET handle=$1, avatar=$2 WHERE id=$3", &params)?;
            },
            None => {
                let params = vec![
                    as_nullable_param(&self.avatar),
                    ParameterValue::Str(self.handle.as_str())
                ];
                db::execute(db_url, "UPDATE profiles SET avatar=$1 WHERE handle=$2", &params)?;
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
                db::execute(db_url, "DELETE FROM profiles WHERE id=$1", &params)?;
            },
            None => {
                let params = vec![
                    ParameterValue::Str(self.handle.as_str())
                ];
                db::execute(db_url, "DELETE FROM profiles WHERE handle=$1", &params)?;
            }
        }
        Ok(())
    }
}
