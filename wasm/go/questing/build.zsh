#!/bin/zsh

GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o questing.wasm
cp questing.wasm ~/creatures/Triptych-Labs-Official-Site/public/questing.wasm

#scp someplace.wasm ddigiacomo@10.145:/noshit/triptych_labs/homepage/public
#scp someplace.wasm ddigiacomo@10.145:/noshit/triptych_labs/homepage/public/marketplace

#scp someplace.wasm ddigiacomo@10.145:/noshit/cdn
