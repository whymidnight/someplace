package ops

import (
	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"github.com/gagliardetto/solana-go"
)

func VerifyList(oracle, batch solana.PublicKey, index uint64) {
	listing, _ := storefront.GetListing(oracle, batch, index)

	storefront.GetListingData(listing)
}

func List(oracle solana.PublicKey, batch solana.PublicKey, index, lifecycleStart, price uint64) solana.Instruction {
	treasuryAuthority, _ := storefront.GetTreasuryAuthority(oracle)

	listing, _ := storefront.GetListing(oracle, batch, index)
	listIx := someplace.NewCreateListingInstructionBuilder().
		SetBatchAccount(batch).
		SetConfigIndex(index).
		SetLifecycleStart(lifecycleStart).
		SetListingAccount(listing).
		SetOracleAccount(oracle).
		SetPrice(price).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTreasuryAuthorityAccount(treasuryAuthority)

	if err := listIx.Validate(); err != nil {
		panic(err)
	}

	return listIx.Build()

}

func Modify(oracle solana.PublicKey, batch solana.PublicKey, index, lifecycleStart, price uint64, isListed bool) solana.Instruction {
	treasuryAuthority, _ := storefront.GetTreasuryAuthority(oracle)

	listing, _ := storefront.GetListing(oracle, batch, index)
	listIx := someplace.NewModifyListingInstructionBuilder().
		SetBatchAccount(batch).
		SetConfigIndex(index).
		SetLifecycleStart(lifecycleStart).
		SetListingAccount(listing).
		SetOracleAccount(oracle).
		SetIsListed(isListed).
		SetPrice(price).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTreasuryAuthorityAccount(treasuryAuthority)

	if err := listIx.Validate(); err != nil {
		panic(err)
	}

	return listIx.Build()

}
