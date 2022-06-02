package ops

import (
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func EndQuest(oracle, initializer, pixelBallzMint solana.PublicKey, questIndex uint64) *questing.Instruction {
	questor, _ := quests.GetQuestorAccount(initializer)
	questee, _ := quests.GetQuesteeAccount(pixelBallzMint)

	questPda, questPdaBump := quests.GetQuest(oracle, questIndex)
	questData := quests.GetQuestData(questPda)
	questAccount, _ := quests.GetQuestAccount(questor, questee, questPda)
	questDeposit, questDepositBump := quests.GetQuestDepositTokenAccount(questee, questPda)

	pixelBallzTokenAddress, _ := utils.GetTokenWallet(initializer, pixelBallzMint)

	questQuesteeReceipt, _ := quests.GetQuestQuesteeReceiptAccount(questor, questee, questPda)

	endQuestIx := questing.NewEndQuestInstructionBuilder().
		SetAssociatedTokenProgramAccount(solana.SPLAssociatedTokenAccountProgramID).
		SetDepositTokenAccountAccount(questDeposit).
		SetDepositTokenAccountBump(questDepositBump).
		SetInitializerAccount(initializer).
		SetOracleAccount(oracle).
		SetPixelballzMintAccount(pixelBallzMint).
		SetPixelballzTokenAccountAccount(pixelBallzTokenAddress).
		SetQuestAccAccount(questAccount).
		SetQuestAccount(questPda).
		SetQuesteeAccount(questee).
		SetQuestorAccount(questor).
		SetQuestQuesteeReceiptAccount(questQuesteeReceipt).
		SetRentAccount(solana.SysVarRentPubkey).
		SetQuestBump(questPdaBump).
		SetSlotHashesAccount(solana.MustPublicKeyFromBase58("SysvarS1otHashes111111111111111111111111111")).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTokenProgramAccount(solana.TokenProgramID)

	for _, reward := range questData.Rewards {
		endQuestIx.Append(&solana.AccountMeta{PublicKey: reward.MintAddress, IsWritable: true, IsSigner: false})
	}
	for _, reward := range questData.Rewards {
		rewardAta, _ := utils.GetTokenWallet(initializer, reward.MintAddress)
		endQuestIx.Append(&solana.AccountMeta{PublicKey: rewardAta, IsWritable: true, IsSigner: false})
	}

	if e := endQuestIx.Validate(); e != nil {
		fmt.Println(e.Error())
		return nil
	}

	return endQuestIx.Build()
}
