module creaturez.nft/wasm/v2

go 1.16

require (
	creaturez.nft/someplace v0.0.0
	creaturez.nft/questing v0.0.0
	creaturez.nft/utils v0.0.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gagliardetto/binary v0.6.1
	github.com/gagliardetto/gofuzz v1.2.2 // indirect
	github.com/gagliardetto/metaplex-go v0.2.1
	github.com/gagliardetto/solana-go v1.4.0
	github.com/gagliardetto/treeout v0.1.4 // indirect
)

replace creaturez.nft/someplace => ../../../sdk/go/someplace_sdk
replace creaturez.nft/questing => ../../../sdk/go/questing

replace creaturez.nft/utils => ../../../sdk/go/utils
