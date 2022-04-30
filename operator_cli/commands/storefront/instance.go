package storefront

import (
	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront/ops"
	"github.com/gagliardetto/solana-go"
)

// Instance will execute required operations that configure a storefront for `oracle`.
func Instance(oracle solana.PrivateKey, mint solana.PublicKey, splits []someplace.Split) {
	/*
	   Required ops include,

	   enable batching for max(u64) candy machines
	   initialize a treasury for storefront under `oracle`
	*/
	ops.EnableBatches(oracle)
	ops.Treasure(oracle, mint)
	ops.AmmendStorefrontSplits(oracle, splits)
}
