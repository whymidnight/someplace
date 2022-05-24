use anyhow::Result;
use clap::Parser;
use futures::executor::block_on;
use sdk::constants::URL;
use sdk::structs::{Client, Cluster};
use someplace::helper_fns::*;
use someplace::structs::ConfigLine;
use std::cell::RefCell;

extern crate base64;

#[derive(Parser, Debug)]
pub struct Opts {
    #[clap(long)]
    candy_machine: String,
    #[clap(long)]
    endpoint: String,
}

pub async fn get_candies() {
    let opts = Opts::parse();
    let mut candies: Vec<ConfigLine> = Vec::new();
    let (rpc, wss) = (
        format!("{}", opts.endpoint),
        format!("wss://{}", opts.endpoint),
    );
    let surf_client = surf::client();
    let client = Client::init(
        Cluster {
            rpc: rpc.to_string(),
            wss: wss.to_string(),
        },
        surf_client,
    );

    let data = client
        .get_account(opts.candy_machine.to_string())
        .await
        .unwrap();
    let result = data.result;
    if result.value.is_some() {
        let encoded = result.value.unwrap().data[0].clone();
        let decoded = base64::decode(encoded).unwrap();

        let data_ref = RefCell::new(decoded.as_slice());
        let configlines = get_config_count_ref(&data_ref.borrow_mut()).unwrap();

        for i in 0..configlines {
            let configline = get_config_lines(data_ref.borrow_mut(), i).unwrap();
            candies.push(configline);
        }
    }
    println!("{}", serde_json::to_string(&candies).unwrap());
}

pub async fn get_candies_of_cardinality(cardinality: String) {
    let opts = Opts::parse();
    let mut candies: Vec<ConfigLine> = Vec::new();
    let (rpc, wss) = URL;
    let surf_client = surf::client();
    let client = Client::init(
        Cluster {
            rpc: rpc.to_string(),
            wss: wss.to_string(),
        },
        surf_client,
    );

    let data = client
        .get_account(opts.candy_machine.to_string())
        .await
        .unwrap();
    let result = data.result;
    if result.value.is_some() {
        let mut assets: Vec<ConfigLine> = Vec::new();
        let encoded = result.value.unwrap().data[0].clone();
        let decoded = base64::decode(encoded).unwrap();

        let data_ref = RefCell::new(decoded.as_slice());
        let configlines = get_config_count_ref(&data_ref.borrow_mut()).unwrap();

        for i in 0..configlines {
            let configline = get_config_lines(data_ref.borrow_mut(), i).unwrap();
            assets.push(configline);
        }

        candies = assets
            .into_iter()
            .filter(|asset| asset.cardinality == cardinality)
            .collect();
    }
    println!("{}", serde_json::to_string(&candies).unwrap());
}

pub async fn get_candy_of_index(index: usize) {
    let opts = Opts::parse();
    let mut candies: Vec<ConfigLine> = Vec::new();
    let (rpc, wss) = URL;
    let surf_client = surf::client();
    let client = Client::init(
        Cluster {
            rpc: rpc.to_string(),
            wss: wss.to_string(),
        },
        surf_client,
    );

    let data = client
        .get_account(opts.candy_machine.to_string())
        .await
        .unwrap();
    let result = data.result;
    if result.value.is_some() {
        let mut assets: Vec<ConfigLine> = Vec::new();
        let encoded = result.value.unwrap().data[0].clone();
        let decoded = base64::decode(encoded).unwrap();

        let data_ref = RefCell::new(decoded.as_slice());

        let configline = get_config_lines(data_ref.borrow_mut(), index).unwrap();
        assets.push(configline);

        candies = assets;
    }
    println!("{}", serde_json::to_string(&candies).unwrap());
}

fn main() -> Result<()> {
    let future = get_candies();
    block_on(future);

    Ok(())
}
