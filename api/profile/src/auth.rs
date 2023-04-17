use chrono::{TimeZone, Utc};

use anyhow::{bail, Result};
use base64::{alphabet, engine, Engine as _};
use jwt_simple::prelude::*;
use serde::{Deserialize, Serialize};
use spin_sdk::outbound_http;

// base64 decoding should support URL safe with no padding and allow trailing bits for JWT tokens
const BASE64_CONFIG: engine::GeneralPurposeConfig = engine::GeneralPurposeConfig::new()
    .with_decode_allow_trailing_bits(true)
    .with_decode_padding_mode(engine::DecodePaddingMode::RequireNone);
const BASE64_ENGINE: engine::GeneralPurpose =
    engine::GeneralPurpose::new(&alphabet::URL_SAFE, BASE64_CONFIG);

#[derive(Serialize, Deserialize, Debug)]
pub(crate) struct JsonWebKey {
    #[serde(rename = "alg")]
    algorithm: String,
    #[serde(rename = "kty")]
    key_type: String,
    #[serde(rename = "use")]
    public_key_use: String,
    #[serde(rename = "n")]
    modulus: String,
    #[serde(rename = "e")]
    exponent: String,
    #[serde(rename = "kid")]
    identifier: String,
    #[serde(rename = "x5t")]
    thumbprint: String,
    #[serde(rename = "x5c")]
    chain: Vec<String>,
}

impl JsonWebKey {
    //TODO: cache the public key after it's been computed
    pub fn to_rsa256_public_key(self) -> Result<RS256PublicKey> {
        let n = BASE64_ENGINE.decode(self.modulus)?;
        let e = BASE64_ENGINE.decode(self.exponent)?;
        Ok(RS256PublicKey::from_components(&n, &e)?.with_key_id(self.identifier.as_str()))
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub(crate) struct JsonWebKeySet {
    keys: Vec<JsonWebKey>,
}

fn get_cached_jwks(store: &spin_sdk::key_value::Store) -> Option<bytes::Bytes> {
    let expiry = store.get("jwks_ttl").ok()?;
    let expiry = expiry.try_into().ok()?;
    let expiry = Utc
        .timestamp_millis_opt(i64::from_le_bytes(expiry))
        .single()?;
    if expiry <= Utc::now() {
        return None;
    }
    let jwks = store.get("jwks").ok()?;
    Some(bytes::Bytes::from(jwks))
}

fn set_cached_jwks(store: &spin_sdk::key_value::Store, jwks: bytes::Bytes) -> Result<()> {
    let expiry = Utc::now() + chrono::Duration::minutes(5);
    let expiry = expiry.timestamp_millis().to_le_bytes();
    store.set("jwks_ttl", expiry)?;
    store.set("jwks", jwks)?;
    Ok(())
}

impl JsonWebKeySet {
    pub fn get(url: String, store: &spin_sdk::key_value::Store) -> Result<Self> {
        let jwks_bytes = match get_cached_jwks(&store) {
            Some(jwks) => jwks,
            None => {
                let res = outbound_http::send_request(
                    http::Request::builder().method("GET").uri(url).body(None)?,
                )?;
                let res_body = match res.body().as_ref() {
                    Some(bytes) => bytes.slice(..),
                    None => bytes::Bytes::default(),
                };
                set_cached_jwks(&store, res_body.to_owned())?;
                res_body
            }
        };
        Ok(serde_json::from_slice::<JsonWebKeySet>(&jwks_bytes)?)
    }

    pub fn verify(
        self,
        token: &str,
        options: Option<VerificationOptions>,
    ) -> Result<JWTClaims<NoCustomClaims>> {
        for key in self.keys {
            let key = key.to_rsa256_public_key()?;

            // add a required key id verification to options
            let options = options.clone().map(|o| VerificationOptions {
                // ensure the token is validated by this key specifically
                required_key_id: key.key_id().to_owned(),
                ..o
            });

            let claims = key.verify_token::<NoCustomClaims>(token, options);
            if claims.is_ok() {
                return claims;
            }
        }
        bail!("No key in the set was able to verify the token.")
    }
}
