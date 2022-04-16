use sdk::constants::URL;
use sdk::structs::{Client, Cluster};
use someplace::helper_fns::*;
use someplace::structs::ConfigLine;
use std::cell::RefCell;
use wasm_bindgen::JsValue;
use web_sys::console;

pub async fn get_candies(candy_machine_id: String) -> Result<JsValue, JsValue> {
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

    console::log_1(&JsValue::from_str(
        format!("{} {}", "debug", candy_machine_id).as_str(),
    ));
    let data = client.get_account(candy_machine_id).await.unwrap();
    let result = data.result;
    if result.value.is_some() {
        let encoded = result.value.unwrap().data[0].clone();
        let decoded = base64::decode(encoded).unwrap();

        let data_ref = RefCell::new(decoded.as_slice());
        let configlines = get_config_count_ref(&data_ref.borrow_mut()).unwrap();
        console::log_1(&JsValue::from_str(format!("{}", configlines).as_str()));

        for i in 0..configlines {
            let configline = get_config_lines(data_ref.borrow_mut(), i).unwrap();
            candies.push(configline);
        }
    }
    Ok(JsValue::from_str(
        format!("{}", serde_json::to_string(&candies).unwrap()).as_str(),
    ))
}
