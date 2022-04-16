use crate::structs::{
    Client, Cluster, Config, JSONRPCRequest, JSONRPCRequestParams, JSONRPCRequestParamsEncoding,
    JSONRPCResponse,
};
pub use anchor_lang;
use anyhow::Result;

impl Client {
    pub fn init(cluster: Cluster, client: surf::Client) -> Self {
        Self {
            cfg: Config { cluster },
            client,
        }
    }

    pub async fn get_account(&self, account_address: String) -> Result<JSONRPCResponse> {
        let body = serde_json::to_string(&JSONRPCRequest::get_account_info(account_address))?;
        let mut res = self
            .client
            .post(self.cfg.cluster.rpc.to_string())
            .content_type("application/json")
            .body_string(body)
            .await
            .unwrap();
        let data = res.body_bytes().await.unwrap();
        let rpc_response: JSONRPCResponse = serde_json::from_slice(data.as_slice()).unwrap();

        Ok(rpc_response)
    }
}

impl JSONRPCRequest {
    pub fn get_account_info(address: String) -> JSONRPCRequest {
        JSONRPCRequest {
            jsonrpc: "2.0".to_string(),
            id: 1,
            method: "getAccountInfo".to_string(),
            params: JSONRPCRequestParams(
                address,
                JSONRPCRequestParamsEncoding {
                    encoding: "jsonParsed".to_string(),
                },
            ),
        }
    }
}
