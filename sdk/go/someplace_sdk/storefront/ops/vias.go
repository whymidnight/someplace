package ops

import (
	"fmt"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func EnableVias(oracle solana.PrivateKey) {
	treasuryAuthority, _ := storefront.GetTreasuryAuthority(oracle.PublicKey())
	vias, _ := storefront.GetVias(oracle.PublicKey())

	enableViasIx := someplace.NewEnableViasInstructionBuilder().
		SetOracleAccount(oracle.PublicKey()).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTreasuryAuthorityAccount(treasuryAuthority).
		SetViasAccount(vias)

	if e := enableViasIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	utils.SendTx(
		"sell",
		append(make([]solana.Instruction, 0), enableViasIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}

func EnableViasForRarityTokens(oracle solana.PublicKey, vias []someplace.ViaMint) []solana.Instruction {
	viaIxs := make([]solana.Instruction, 0)

	viasPda, _ := storefront.GetVias(oracle)
	viasData := storefront.GetViasData(viasPda)
    fmt.Println(viasPda, viasData)
	for i, via := range vias {
		treasuryAuthority, _ := storefront.GetTreasuryAuthority(oracle)
		viaPda, _ := storefront.GetVia(oracle, viasData.Vias+uint64(i))
		viaMapping, _ := storefront.GetViaMapping(oracle, via.MintAddress)

		enableRarityTokenIx := someplace.NewEnableViaRarityTokenMintingInstructionBuilder().
			SetOracleAccount(oracle).
			SetRarity(via.Rarity).
			SetRarityTokenMintAccount(via.MintAddress).
			SetSystemProgramAccount(solana.SystemProgramID).
			SetTreasuryAuthorityAccount(treasuryAuthority).
			SetViaAccount(viaPda).
			SetViaMappingAccount(viaMapping).
			SetViasAccount(viasPda)

		if e := enableRarityTokenIx.Validate(); e != nil {
			fmt.Println(e.Error())
			panic("...")
		}
		viaIxs = append(viaIxs, enableRarityTokenIx.Build())
	}

	return viaIxs
}
