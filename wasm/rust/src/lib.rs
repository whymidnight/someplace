use crate::candymachine::*;
use wasm_bindgen::prelude::wasm_bindgen;
use wasm_bindgen::JsValue;
use wasm_bindgen_futures::future_to_promise;
extern crate base64;

// When the `wee_alloc` feature is enabled, use `wee_alloc` as the global
// allocator.
#[cfg(feature = "wee_alloc")]
#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

pub mod candymachine;

#[wasm_bindgen]
pub async fn fetch_candies(candy_machine_id: String) -> JsValue {
    let candies = future_to_promise(get_candies(candy_machine_id));
    let result = wasm_bindgen_futures::JsFuture::from(candies).await.unwrap();

    result
}
