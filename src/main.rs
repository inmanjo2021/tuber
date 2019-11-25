extern crate serde;
extern crate serde_derive;
extern crate reqwest;

mod api;
mod fetch_layer;

use std::error::Error;
// use std::result::Result;
// use reqwest::Error;
use dotenv::dotenv;
// use std::env;
// use serde::Deserialize;
// use api;


fn main() -> Result<(), Box<dyn Error>> {
    dotenv().ok();

    // let response = api::get_deployment(&"component-library");
    // println!("{}", response.token);
    // println!("{:?}", response.text());
    fetch_layer::get_layers()?;

    Ok(())
}
