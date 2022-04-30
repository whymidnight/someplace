package storefront

import (
	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront/ops"
	"github.com/gagliardetto/solana-go"
)

// AmmendStorefrontSplits will configure the storefront to spread distribution of tenders during mint to specified holders. Up to 10 splits may be defined.
func AmmendStorefrontSplits(oracle solana.PrivateKey, splits []someplace.Split) {

	ops.AmmendStorefrontSplits(oracle, splits)

}
