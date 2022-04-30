#!/bin/zsh

npx ts-node-script ./src/candy-machine-v2-cli.ts upload \
    -e devnet \
    -k ./oracle.key \
    -cp ./config.json \
    -c example \
    -r https://api.devnet.solana.com \
    ./assets
