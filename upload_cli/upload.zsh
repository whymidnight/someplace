#!/bin/zsh

npx ts-node-script ./src/candy-machine-v2-cli.ts upload \
    -e testnet \
    -k ./oracle.key \
    -cp ./config.json \
    -c example \
    -r https://api.testnet.solana.com \
    ./assets
