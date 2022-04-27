module creaturez.nft/operatorCLI/v2

go 1.16

replace creaturez.nft/someplace => ../sdk/go/someplace_sdk

replace creaturez.nft/utils => ../sdk/go/utils

require (
	creaturez.nft/someplace v0.0.0-00010101000000-000000000000
	creaturez.nft/utils v0.0.0-00010101000000-000000000000
	github.com/gagliardetto/solana-go v1.4.0
	github.com/go-gota/gota v0.12.0
	gopkg.in/yaml.v2 v2.4.0
)
