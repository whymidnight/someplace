package ops

import (
	"creaturez.nft/someplace"
	"creaturez.nft/someplace/marketplace"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func MarketCreate(oracle solana.PrivateKey, mint, marketUid solana.PublicKey) {
	marketAuthority, _ := marketplace.GetMarketAuthority(oracle.PublicKey(), marketUid)

	listIx := someplace.NewInitMarketInstructionBuilder().
		SetMarketAuthorityAccount(marketAuthority).
		SetMarketMintAccount(mint).
		SetMarketUid(marketUid).
		SetName("market test").
		SetOracleAccount(oracle.PublicKey()).
		SetSystemProgramAccount(solana.SystemProgramID)

	if err := listIx.Validate(); err != nil {
		panic(err)
	}

	utils.SendTx(
		"list",
		append(make([]solana.Instruction, 0), listIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)

}

func VerifyMarketCreate(oracle, marketUid solana.PublicKey) {
	marketAuthority, _ := marketplace.GetMarketAuthority(oracle, marketUid)
	marketplace.GetMarketAuthorityData(marketAuthority)
}

