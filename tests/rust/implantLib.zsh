#!/bin/zsh

# build then copy the binary to `operator_cli`
cargo build
cp ./target/debug/someplace_rusty ./../../operator_cli/libs
