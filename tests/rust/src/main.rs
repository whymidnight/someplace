use anyhow::Result;
use clap::Parser;
use sdk::candymachine::fetch_candies;
use sdk::constants::URL;
use sdk::structs::{Client, Cluster};

#[derive(Parser, Debug)]
pub struct Opts {
    #[clap(long)]
    candy_machine: String,
}

// This example assumes a local validator is running with the programs
// deployed at the addresses given by the CLI args.
fn main() -> Result<()> {
    println!("Starting test...");
    let opts = Opts::parse();

    let (rpc, wss) = URL;
    let client = Client::init(Cluster {
        rpc: rpc.to_string(),
        wss: wss.to_string(),
    });
    let candies = fetch_candies(&client, &opts.candy_machine)?;
    println!("{}", serde_json::to_string(&candies)?);
    Ok(())
}
