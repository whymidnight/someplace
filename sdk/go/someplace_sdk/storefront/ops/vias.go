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

func DoRng(oracle, initializer, pixelBallzMint, questPda, questor, questee, rewardMint solana.PublicKey) *someplace.Instruction {
	// rng after quest end
	batches, _ := storefront.GetBatches(oracle)
	batchesData := storefront.GetBatchesData(batches)

	viaMap, _ := storefront.GetViaMapping(oracle, rewardMint)
	viaMapData := storefront.GetViaMappingData(viaMap)

	via, viaBump := storefront.GetVia(batchesData.Oracle, viaMapData.ViasIndex)

	rewardTicket, _ := storefront.GetRewardTicket(via, questPda, questee, initializer)
	rewardTicketData := storefront.GetRewardTicketData(rewardTicket)
	rewardTokenAccount, _ := utils.GetTokenWallet(initializer, viaMapData.TokenMint)

	if rewardTicketData == nil {
		rngRewardIndiceIx := someplace.NewRngNftAfterQuestInstructionBuilder().
			SetBatchesAccount(batches).
			SetInitializerAccount(initializer).
			SetQuestAccount(questPda).
			SetQuesteeAccount(questee).
			SetRewardTicketAccount(rewardTicket).
			SetRewardTokenAccountAccount(rewardTokenAccount).
			SetSlotHashesAccount(solana.MustPublicKeyFromBase58("SysvarS1otHashes111111111111111111111111111")).
			SetSystemProgramAccount(solana.SystemProgramID).
			SetViaAccount(via).
			SetViaBump(viaBump).
			SetViaMapAccount(viaMap)

		if err := rngRewardIndiceIx.Validate(); err != nil {
			panic(err)
		}

		// TODO IMPLEMENT MAX OF `MIN_BATCH_RNG`
		for i := range make([]int, batchesData.Counter) {
			batchReceipt, _ := storefront.GetBatchReceipt(oracle, uint64(i))
			batchReceiptData := storefront.GetBatchReceiptData(batchReceipt)
			batchAccount := batchReceiptData.BatchAccount
			batchCardinalitiesReport, _ := storefront.GetBatchCardinalitiesReport(batchAccount)

			rngRewardIndiceIx.Append(&solana.AccountMeta{
				PublicKey:  batchCardinalitiesReport,
				IsWritable: false,
				IsSigner:   false,
			})
		}
		return rngRewardIndiceIx.Build()
	} else {
		rngRewardIndiceIx := someplace.NewRecycleRngNftAfterQuestInstructionBuilder().
			SetBatchesAccount(batches).
			SetInitializerAccount(initializer).
			SetQuestAccount(questPda).
			SetRewardTicketAccount(rewardTicket).
			SetRewardTokenAccountAccount(rewardTokenAccount).
			SetSlotHashesAccount(solana.MustPublicKeyFromBase58("SysvarS1otHashes111111111111111111111111111")).
			SetSystemProgramAccount(solana.SystemProgramID).
			SetViaAccount(via).
			SetViaBump(viaBump).
			SetViaMapAccount(viaMap)

		if err := rngRewardIndiceIx.Validate(); err != nil {
			panic(err)
		}

		// TODO IMPLEMENT MAX OF `MIN_BATCH_RNG`
		for i := range make([]int, batchesData.Counter) {
			batchReceipt, _ := storefront.GetBatchReceipt(oracle, uint64(i))
			batchReceiptData := storefront.GetBatchReceiptData(batchReceipt)
			batchAccount := batchReceiptData.BatchAccount
			batchCardinalitiesReport, _ := storefront.GetBatchCardinalitiesReport(batchAccount)

			rngRewardIndiceIx.Append(&solana.AccountMeta{
				PublicKey:  batchCardinalitiesReport,
				IsWritable: false,
				IsSigner:   false,
			})
		}
		return rngRewardIndiceIx.Build()
	}

}

func MintViaRarityTicket(oracle, initializer, questPda, questee, rewardMint solana.PublicKey, offset uint64) solana.Instruction {

	viaMap, _ := storefront.GetViaMapping(oracle, rewardMint)
	viaMapData := storefront.GetViaMappingData(viaMap)
	via, _ := storefront.GetVia(oracle, viaMapData.ViasIndex)
	viaData := storefront.GetViaData(via)

	rewardTicket, rewardTicketBump := storefront.GetRewardTicket(via, questPda, questee, initializer)
	rewardTicketData := storefront.GetRewardTicketData(rewardTicket)
	if rewardTicketData == nil {
		panic("null reward ticket data")
	}

	candyMachineAddress := rewardTicketData.BatchAccount

	mint, _, _ := storefront.GetMint(oracle, via, viaData.Mints+offset)
	mintAta, _ := utils.GetTokenWallet(initializer, mint)
	mintViaHash, _, _ := storefront.GetMintHashVia(oracle, viaData.TokenMint, viaData.Mints+offset)

	batchCardinalitiesReport, _ := storefront.GetBatchCardinalitiesReport(candyMachineAddress)

	metadataAddress, err := utils.GetMetadata(mint)
	if err != nil {
		panic(err)
	}
	masterEdition, err := utils.GetMasterEdition(mint)
	if err != nil {
		panic(err)
	}
	candyMachineCreator, creatorBump, err := storefront.GetCandyMachineCreator(candyMachineAddress)
	if err != nil {
		panic(err)
	}
	rewardTokenAccount, _ := utils.GetTokenWallet(initializer, viaMapData.TokenMint)

	mintIx := someplace.NewMintNftViaInstructionBuilder().
		SetClockAccount(solana.SysVarClockPubkey).
		SetCreatorBump(creatorBump).
		SetMintAccount(mint).
		SetViaAccount(via).
		SetMintAtaAccount(mintAta).
		SetCandyMachineAccount(rewardTicketData.BatchAccount).
		SetMasterEditionAccount(masterEdition).
		SetCandyMachineCreatorAccount(candyMachineCreator).
		SetBatchCardinalitiesReportAccount(batchCardinalitiesReport).
		SetMetadataAccount(metadataAddress).
		SetRewardTicketAccount(rewardTicket).
		SetSlotHashesAccount(solana.MustPublicKeyFromBase58("SysvarS1otHashes111111111111111111111111111")).
		SetInstructionSysvarAccountAccount(solana.SysVarInstructionsPubkey).
		SetMintHashAccount(mintViaHash).
		SetOracleAccount(oracle).
		SetPayerAccount(initializer).
		SetRentAccount(solana.SysVarRentPubkey).
		SetRewardTicketBump(rewardTicketBump).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetRewardTokenAccountAccount(rewardTokenAccount).
		SetRewardTokenMintAccountAccount(viaMapData.TokenMint).
		SetTokenProgramAccount(solana.TokenProgramID).
		SetTokenMetadataProgramAccount(solana.TokenMetadataProgramID).
		SetAssociatedTokenProgramAccount(solana.SPLAssociatedTokenAccountProgramID)

	err = mintIx.Validate()
	if err != nil {
		panic(err)
	}

	return mintIx.Build()

}

type ViaHash struct {
	MintAddress solana.PublicKey
	MintHash    someplace.MintHash
}

func GetViasHashMap(oracle solana.PublicKey) map[solana.PublicKey][]ViaHash {
	viasHashMap := make(map[solana.PublicKey][]ViaHash)

	viasPda, _ := storefront.GetVias(oracle)
	viasData := storefront.GetViasData(viasPda)
	for i := range make([]int, viasData.Vias) {
		via, _ := storefront.GetVia(oracle, uint64(i))
		viaData := storefront.GetViaData(via)
		mints := make([]ViaHash, 0)
		for ii := range make([]int, viaData.Mints) {
			viaMint, _, _ := storefront.GetMint(oracle, via, uint64(ii))
			mintViaHash, _, _ := storefront.GetMintHashVia(oracle, viaData.TokenMint, uint64(ii))
			mintViaHashData := storefront.GetMintHashData(mintViaHash)
			if mintViaHashData == nil {
				continue
			}

			mints = append(mints, ViaHash{
				MintAddress: viaMint,
				MintHash:    *mintViaHashData,
			})
		}
		viasHashMap[via] = mints
	}

	return viasHashMap
}
