extern crate flate2;
extern crate tar;
extern crate hyper;
extern crate reqwest;

use std::error::Error;
use serde::{Deserialize};
// use serde::de::{self, Visitor};
// use serde_json;
// use std::fmt;
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

const MEGABYTE: u32 = 1_000_000;
const MAX_SIZE: u32 = MEGABYTE * 1;

#[derive(Deserialize, Debug)]
struct AuthResponse {
    pub token: String,
}

#[derive(Deserialize)]
#[serde(rename_all = "camelCase")]
struct Layer {
    digest: String,
    size:   u32,
}

#[derive(Deserialize)]
#[serde(rename_all = "camelCase")]
struct Manifest {
    layers: Vec<Layer>,
}

#[derive(Debug)]
pub struct DownloadLayerResponse {
    should_continue: bool,
    pub yamls: Vec<Yaml>,
}

#[derive(Debug)]
pub struct Yaml {
    pub content: String,
    pub filename: String,
}

fn get_token() -> Result<AuthResponse, Box<dyn Error>> {
    let request_url = format!(
        "{auth_base}/v2/token?scope=repository:{image}:pull",
        auth_base = env::var("AUTH_BASE")?,
        image = env::var("IMAGE_NAME")?,
    );

    let client = reqwest::Client::new();
    let mut response = client.get(&request_url)
        .basic_auth("_token", Some(env::var("GCLOUD_TOKEN")?))
        .send()?
        .error_for_status()
        .unwrap();

    let auth_response: AuthResponse = response.json()?;
    // println!("{:?}", auth_response);

    Ok(auth_response)
}

fn get_layers() -> Result<Vec<Layer>, Box<dyn Error>> {
    let token = get_token()?.token;

    let request_url = format!(
        "{registry_base}/v2/{image}/manifests/{tag}",
        image = env::var("IMAGE_NAME")?,
        registry_base = env::var("REGISTRY_BASE")?,
        tag = env::var("IMAGE_TAG")?,
    );

    let client = reqwest::Client::new();

    // println!("{}", token);

    let mut response = client.get(&request_url)
        .header("Authorization", format!("Bearer {}", token))
        // .header("Accept", "application/vnd.docker.container.image.v2+json")
        .send()?
        .error_for_status()
        .unwrap();

    let parsed: Manifest = response.json()?;

    Ok(parsed.layers)
}

pub fn find_layer<'a>() -> Result<DownloadLayerResponse, Box<dyn Error>> {
    let layers = get_layers()?;

    for layer in layers.into_iter().rev() {
        if layer.size > MAX_SIZE {
            println!("Layer to large, skipping...");
            continue;
        }

        let res = download_layer(&layer).unwrap();

        if !res.should_continue {
            return Ok(res);
        }
    };

    panic!("No tuber layer found")
}

fn download_layer(layer: &Layer) -> Result<DownloadLayerResponse, Box<dyn Error>> {
    let token = get_token()?.token;
    let layer = &layer.digest;

    let request_url = format!(
        "{registry_base}/v2/{image}/blobs/{layer}",
        image = env::var("IMAGE_NAME")?,
        registry_base = env::var("REGISTRY_BASE")?,
        layer = layer,
    );

    let client = reqwest::Client::new();

    let response = client.get(&request_url)
        .header("Authorization", format!("Bearer {}", token))
        .send()?;

    let tar = GzDecoder::new(response);
    let mut archive = Archive::new(tar);
    let mut should_continue = true;
    let mut yamls: Vec<Yaml> = vec![];

    for file in archive.entries().unwrap() {
        let mut file = file.unwrap();
        let mut yaml = String::new();

        // this read happens here, before the early exits because `rustc --explain E0502`
        file.read_to_string(&mut yaml).unwrap_or(0);

        let path = file.path().unwrap().clone();

        if !path.starts_with(".tuber") {
            break;
        }

        let ext = path.extension().unwrap_or(::std::ffi::OsStr::new(""));
        let filename = path.file_name().unwrap();

        should_continue = false;

        if ext != "yaml" {
            continue;
        }

        yamls.push(Yaml {
            content:  yaml,
            filename: filename.to_string_lossy().to_string(),
        });
    }

    Ok(DownloadLayerResponse {
        should_continue: should_continue,
        yamls:           yamls,
    })
}
