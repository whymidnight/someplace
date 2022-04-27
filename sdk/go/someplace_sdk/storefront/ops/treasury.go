package ops

import (
	"fmt"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func Treasure(oracle solana.PrivateKey, mint solana.PublicKey) {
	treasuryAuthority, _ := storefront.GetTreasuryAuthority(oracle.PublicKey())
	treasuryTokenAccount, _ := storefront.GetTreasuryTokenAccount(oracle.PublicKey())
	oracleTokenAccount, _ := utils.GetTokenWallet(oracle.PublicKey(), mint)

	treasuryIx := someplace.NewInitTreasuryInstructionBuilder().
		SetAdornment("fedcoin").
		SetOracleAccount(oracle.PublicKey()).
		SetOracleTokenAccountAccount(oracleTokenAccount).
		SetRentAccount(solana.SysVarRentPubkey).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTokenProgramAccount(solana.TokenProgramID).
		SetTreasuryAuthorityAccount(treasuryAuthority).
		SetTreasuryTokenAccountAccount(treasuryTokenAccount).
		SetTreasuryTokenMintAccount(mint)

	if e := treasuryIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	utils.SendTx(
		"sell",
		append(make([]solana.Instruction, 0), treasuryIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}

func TreasureWhitelistCandyMachine(oracle solana.PrivateKey, candyMachine solana.PublicKey) {
	treasuryAuthority, _ := storefront.GetTreasuryAuthority(oracle.PublicKey())
	candyMachineCreator, _, _ := storefront.GetCandyMachineCreator(candyMachine)
	treasuryWhitelist, _ := storefront.GetTreasuryWhitelist(oracle.PublicKey(), treasuryAuthority, candyMachineCreator)

	treasuryIx := someplace.NewAddWhitelistedCmInstructionBuilder().
		SetCandyMachine(candyMachine).
		SetCandyMachineCreator(candyMachineCreator).
		SetOracleAccount(oracle.PublicKey()).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTreasuryAuthorityAccount(treasuryAuthority).
		SetTreasuryWhitelistAccount(treasuryWhitelist)

	utils.SendTx(
		"sell",
		append(make([]solana.Instruction, 0), treasuryIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}

func TreasureVerifyCM(oracle, candyMachine solana.PublicKey) {
	treasuryAuthority, _ := storefront.GetTreasuryAuthority(oracle)
	treasuryWhitelist, _ := storefront.GetTreasuryWhitelist(oracle, treasuryAuthority, candyMachine)

	storefront.GetTreasuryWhitelistData(treasuryWhitelist)
}

func TreasureVerify(oracle solana.PublicKey) {
	treasuryAuthority, _ := storefront.GetTreasuryAuthority(oracle)

	storefront.GetTreasuryAuthorityData(treasuryAuthority)
}
