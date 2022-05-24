#!/bin/zsh

npx ts-node-script ./src/candy-machine-v2-cli.ts upload \
    -e devnet \
    -k ./oracle.key \
    -cp ./config.json \
    -c example \
    -r https://psytrbhymqlkfrhudd.dev.genesysgo.net:8899/ \
    ./assets
