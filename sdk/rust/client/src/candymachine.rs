use crate::structs::Client;
use anyhow::Result;
use someplace::helper_fns::*;
use someplace::structs::ConfigLine;
use std::cell::RefCell;

/// fn returns accumulation of all `ConfigLine`s from a candy machine.
pub async fn fetch_candies(client: Client, cm_id: &str) -> Result<Vec<ConfigLine>> {
    let mut candies: Vec<ConfigLine> = Vec::new();
    /*

    let data = client.get_account(cm_id.to_string()).await?;
    let data_ref = RefCell::new(data.as_slice());
    let configlines = get_config_count_ref(&data_ref.borrow_mut())?;
    println!("{}", configlines);
    for i in 0..configlines {
        let configline = get_config_lines(data_ref.borrow_mut(), i)?;
        candies.push(configline);
    }
    */
    Ok(candies)
}
