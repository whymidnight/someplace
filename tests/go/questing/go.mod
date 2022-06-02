module creaturez.nft/questing_tests/v2

go 1.16

require (
	creaturez.nft/questing v0.0.0
	creaturez.nft/someplace v0.0.0
	creaturez.nft/utils v0.0.0-00010101000000-000000000000
	github.com/btcsuite/btcutil v1.0.2
	github.com/gagliardetto/binary v0.6.1
	github.com/gagliardetto/metaplex-go v0.2.1
	github.com/gagliardetto/solana-go v1.4.0
	github.com/mr-tron/base58 v1.2.0 // indirect
)

replace creaturez.nft/someplace => ../../../sdk/go/someplace_sdk

replace creaturez.nft/questing => ../../../sdk/go/questing

replace creaturez.nft/utils => ../../../sdk/go/utils
