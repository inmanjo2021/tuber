#[macro_use]
extern crate serde;
extern crate serde_derive;
extern crate reqwest;

mod api;

use std::error::Error;
// use std::result::Result;
// use reqwest::Error;
use dotenv::dotenv;
use std::env;
use serde::Deserialize;
// use api;


fn main() -> Result<(), Box<dyn Error>> {
    println!("Hello, world!");
    dotenv().ok();

    let response = api::get_deployment(&"component-library");
    println!("{:?}", response?.status());

    Ok(())
}
