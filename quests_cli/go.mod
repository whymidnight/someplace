module creaturez.nft/questingCLI/v2

go 1.16

replace creaturez.nft/questing => ../sdk/go/questing
replace creaturez.nft/someplace => ../sdk/go/someplace_sdk

replace creaturez.nft/utils => ../sdk/go/utils

require (
	creaturez.nft/questing v0.0.0
	creaturez.nft/someplace v0.0.0
	creaturez.nft/utils v0.0.0-00010101000000-000000000000
	github.com/gagliardetto/solana-go v1.4.0
	github.com/go-gota/gota v0.12.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
