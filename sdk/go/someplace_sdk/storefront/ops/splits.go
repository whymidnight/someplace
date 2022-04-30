package ops

import (
	"fmt"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func AmmendStorefrontSplits(oracle solana.PrivateKey, splits []someplace.Split) {
	treasuryAuthority, _ := storefront.GetTreasuryAuthority(oracle.PublicKey())

	splitsIx := someplace.NewAmmendStorefrontSplitsInstructionBuilder().
		SetOracleAccount(oracle.PublicKey()).
		SetStorefrontSplits(splits).
		SetTreasuryAuthorityAccount(treasuryAuthority)

	if e := splitsIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	utils.SendTx(
		"sell",
		append(make([]solana.Instruction, 0), splitsIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}

