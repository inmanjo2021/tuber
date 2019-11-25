use std::error::Error;
use serde::{Deserialize, Deserializer};
use serde::de::{self, Visitor};
use serde_json;

const REGISTRY_BASE:  &str = "https://registry-1.docker.io";
const AUTH_BASE:      &str = "https://auth.docker.io";
const AUTH_SERVICE:   &str = "registry.docker.io";
const IMAGE:          &str = "library/node";
const TAG:            &str = "latest";

#[derive(Deserialize)]
struct AuthResponse {
    pub token: String,
}

#[derive(Deserialize)]
#[serde(rename_all = "camelCase")]
struct Layer {
    blob_sum: String,
}

#[derive(Deserialize)]
struct History {
    #[serde(deserialize_with = "deserialize_layer_command")]
    #[serde(rename(deserialize = "v1Compatibility"))]
    command: String,
}

#[derive(Deserialize)]
struct HistoryData {
    container_config: ContainerConfig,
}

#[derive(Deserialize)]
struct ContainerConfig {
    #[serde(rename(deserialize = "Cmd"))]
    cmd: Vec<String>,
}

#[derive(Deserialize)]
#[serde(rename_all = "camelCase")]
struct Manifest {
    fs_layers: Vec<Layer>,
    history:   Vec<History>,
}

fn deserialize_layer_command<'de, D>(deserializer: D) -> Result<String, D::Error>
where
    D: Deserializer<'de>,
{
    use std::fmt;

    struct TimeVisitor;

    impl<'de> Visitor<'de> for TimeVisitor {
        type Value = String;

        fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
            formatter.write_str("an integer or string containing an integer")
        }

        #[inline]
        fn visit_str<E>(self, value: &str) -> Result<Self::Value, E>
        where
            E: de::Error,
        {
            let deserialized: HistoryData = serde_json::from_str(&value).unwrap();

            Ok(deserialized.container_config.cmd.join(" ").clone())
        }
    }

    deserializer.deserialize_string(TimeVisitor)
}

pub fn get_layers() -> Result<(), Box<dyn Error>> {
    let token = get_token()?.token;

    let request_url = format!(
        "{registry_base}/v2/{image}/manifests/{tag}",
        image = IMAGE,
        registry_base = REGISTRY_BASE,
        tag = TAG,
    );

    let client = reqwest::Client::new();

    let mut response = client.get(&request_url)
        .header("Authorization", format!("Bearer {}", token))
        .send()?;

    // print!("{}", response.text()?);
    let parsed: Manifest = response.json()?;
    let index = parsed.history.into_iter().position(|x| x.command.contains("COPY")).unwrap();
    let layer = parsed.fs_layers[index].blob_sum.clone();

    print!("{}", layer);

    Ok(())
}

fn get_token() -> Result<AuthResponse, Box<dyn Error>> {
    let request_url = format!(
        "{auth_base}/token?service={auth_service}&scope=repository:{image}:pull",
        auth_base = AUTH_BASE,
        auth_service = AUTH_SERVICE,
        image = IMAGE,
    );

    let mut response = reqwest::get(&request_url)?;
    let auth_response: AuthResponse = response.json()?;

    Ok(auth_response)
}
