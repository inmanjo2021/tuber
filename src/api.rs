use std::error::Error;
use std::env;

pub fn make_request(path: &str) -> Result<reqwest::Response, Box<dyn Error>> {
    let host = env::var("K8S_HOST").unwrap();

    let request_url = format!(
        "{host}/{path}",
        host = host,
        path = path,
    );

    let mut response = reqwest::get(&request_url)?;

    Ok(response)
}

pub fn get_deployment(name: &str) -> Result<reqwest::Response, Box<dyn Error>> {
    let response = make_request(
        &format!(
            "apis/apps/v1/namespaces/default/deployments/{name}",
            name = name,
        )
    );

    Ok(response.unwrap())
    // if ()
}
