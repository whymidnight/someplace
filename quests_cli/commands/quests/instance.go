package quests

import (
	"creaturez.nft/questing/quests/ops"
	storefront_ops "creaturez.nft/someplace/storefront/ops"
	"github.com/gagliardetto/solana-go"
)

// Instance will execute required operations that configure a storefront for `oracle`.
func Instance(oracle solana.PrivateKey) {
	/*
	   Required ops include,

	   enable batching for max(u64) candy machines
	   initialize a treasury for storefront under `oracle`
	*/
	ops.EnableQuests(oracle)
	storefront_ops.EnableVias(oracle)
}
