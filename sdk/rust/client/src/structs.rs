use serde::{Deserialize, Serialize};

pub struct Cluster {
    pub rpc: String,
    pub wss: String,
}

pub struct Client {
    pub cfg: Config,
    pub client: surf::Client,
}

pub struct Config {
    pub cluster: Cluster,
}

#[derive(Serialize, Deserialize)]
pub struct JSONRPCRequest {
    pub jsonrpc: String,
    pub id: u8,
    pub method: String,
    pub params: JSONRPCRequestParams,
}

#[derive(Serialize, Deserialize)]
pub struct JSONRPCRequestParams(pub String, pub JSONRPCRequestParamsEncoding);

#[derive(Serialize, Deserialize)]
pub struct JSONRPCRequestParamsEncoding {
    pub encoding: String,
}

#[derive(Serialize, Deserialize)]
pub struct JSONRPCResponse {
    pub jsonrpc: String,
    pub result: JSONRPCResponseResult,
}

#[derive(Serialize, Deserialize)]
pub struct JSONRPCResponseResult {
    pub value: Option<JSONRPCResponseValue>,
}

#[derive(Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct JSONRPCResponseValue {
    pub data: Vec<String>,
    pub executable: bool,
    pub lamports: u64,
    pub owner: String,
    pub rent_epoch: u64,
}
