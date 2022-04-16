#!/bin/zsh

wasm-bindgen --target no-modules ./target/wasm32-unknown-unknown/debug/hello.wasm --out-dir ./build
