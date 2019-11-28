extern crate flate2;
extern crate tar;
extern crate hyper;
extern crate reqwest;

use std::error::Error;
use serde::{Deserialize, Deserializer};
use serde::de::{self, Visitor};
use serde_json;
use std::fmt;
use std::env;
// use std::io::copy;
// use std::fs::File;
// use std::path::Path;
use tar::Archive;
use flate2::read::GzDecoder;
use std::io::Read;
// use hyper::header::{Headers, Authorization, Basic};
// use reqwest::header::{HeaderMap, Authorization, Basic};

// const REGISTRY_BASE:  &str = "https://registry-1.docker.io";
// const AUTH_BASE:      &str = "https://auth.docker.io";
// const AUTH_SERVICE:   &str = "registry.docker.io";

const REGISTRY_BASE:  &str = "https://us.gcr.io";
const AUTH_BASE:      &str = "https://us.gcr.io";
// const AUTH_SERVICE:   &str = "registry.docker.io";
const IMAGE:          &str = "freshly-docker/address-service";
const TAG:            &str = "master";
const MATCHER:        &str = "COPY";

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
    struct LayerCommandVisitor;

    impl<'de> Visitor<'de> for LayerCommandVisitor {
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

    deserializer.deserialize_string(LayerCommandVisitor)
}

fn get_token() -> Result<AuthResponse, Box<dyn Error>> {
    // This is the auth format for docker hub.
    // Unsure where the distinctions between these two are specified.
    // let request_url = format!(
    //     "{auth_base}/token?service={auth_service}&scope=repository:{image}:pull",
    //     auth_base = AUTH_BASE,
    //     auth_service = "registry.docker.io",
    //     image = IMAGE,
    // );

    // GCP does not use the service query string param.
    // It also needs basic auth, which makes sense i guess?
    let request_url = format!(
        "{auth_base}/v2/token?scope=repository:{image}:pull",
        auth_base = AUTH_BASE,
        image = IMAGE,
    );

    let client = reqwest::Client::new();
    let mut response = client.get(&request_url)
        .basic_auth("_token", Some(env::var("GCLOUD_TOKEN")?))
        .send()?
        .error_for_status()
        .unwrap();

    let auth_response: AuthResponse = response.json()?;

    println!("TOKEN IS !!!! : ");
    println!("{}", auth_response.token);

    Ok(auth_response)
}

pub fn find_layer() -> Result<String, Box<dyn Error>> {
    download_layer()
}

pub fn download_layer() -> Result<String, Box<dyn Error>> {
    let token = get_token()?.token;
    let layer = get_layer_sha().unwrap();

    let request_url = format!(
        "{registry_base}/v2/{image}/blobs/{layer}",
        image = IMAGE,
        registry_base = REGISTRY_BASE,
        layer = layer,
    );

    let client = reqwest::Client::new();

    let response = client.get(&request_url)
        .header("Authorization", format!("Bearer {}", token))
        .send()?;

    let tar = GzDecoder::new(response);
    let mut archive = Archive::new(tar);

    for file in archive.entries().unwrap() {
        let mut file = file.unwrap();

        // files implement the Read trait
        let mut s = String::new();
        file.read_to_string(&mut s).unwrap();
        println!("{}", s);
    }

    Ok(layer)
}

fn get_layer_sha() -> Result<String, Box<dyn Error>> {
    let token = get_token()?.token;

    let request_url = format!(
        "{registry_base}/v2/{image}/manifests/{tag}",
        image = IMAGE,
        registry_base = REGISTRY_BASE,
        tag = TAG,
    );

    let client = reqwest::Client::new();

    // println!("{}", token);

    let mut response = client.get(&request_url)
        .header("Authorization", format!("Bearer {}", token))
        // .header("Accept", "application/vnd.docker.container.image.v1+json")
        .send()?
        .error_for_status()
        .unwrap();

    println!("{}", response.text().unwrap());

    let parsed: Manifest = response.json()?;
    let index = parsed.history.into_iter().position(|x| x.command.contains(MATCHER)).unwrap();
    let layer = parsed.fs_layers[index].blob_sum.clone();

    Ok(layer)
}
