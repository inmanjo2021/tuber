extern crate serde;
extern crate serde_derive;
extern crate reqwest;

mod api;
mod fetch_layer;

use std::error::Error;
use dotenv::dotenv;

fn main() -> Result<(), Box<dyn Error>> {
    dotenv().ok();
    fetch_layer::find_layer()?;
    let res = fetch_layer::find_layer()?;

    for yaml in res.yamls {
        println!("###############################################");
        println!("{}", yaml.filename);
        println!("###############################################");
        println!("{}", yaml.content);
    }

    Ok(())
}
