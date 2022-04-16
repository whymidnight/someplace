#!/bin/zsh

./build.zsh;
./bind.zsh;
# scp build/hello_bg.wasm ddigiacomo@10.145:/noshit/triptych_labs/homepage/public/marketplace;
# scp build/hello_bg.wasm ddigiacomo@10.145:/noshit/triptych_labs/homepage/public/mint;
scp build/hello_bg.wasm ddigiacomo@10.145:/noshit/triptych_labs/homepage/public;
scp build/hello.js ddigiacomo@10.145:/noshit/triptych_labs/homepage/public

# scp build/hello_bg.wasm ddigiacomo@10.145:/noshit/triptych_labs/homepage/build/marketplace;
# scp build/hello_bg.wasm ddigiacomo@10.145:/noshit/triptych_labs/homepage/build/mint;
# scp build/hello_bg.wasm ddigiacomo@10.145:/noshit/triptych_labs/homepage/build;
# scp build/hello.js ddigiacomo@10.145:/noshit/triptych_labs/homepage/build
