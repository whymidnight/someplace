package ops

import (
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func StartQuest(oracle, initializer, pixelBallzMint solana.PublicKey, questIndex uint64) *questing.Instruction {
	questor, _ := quests.GetQuestorAccount(initializer)
	questee, _ := quests.GetQuesteeAccount(pixelBallzMint)

	questPda, _ := quests.GetQuest(oracle, questIndex)
	questData := quests.GetQuestData(questPda)
	questAccount, _ := quests.GetQuestAccount(questor, questee, questPda)
	questDeposit, _ := quests.GetQuestDepositTokenAccount(questee, questPda)

	pixelBallzTokenAddress, _ := utils.GetTokenWallet(initializer, pixelBallzMint)

	startQuestIx := questing.NewStartQuestInstructionBuilder().
		SetDepositTokenAccountAccount(questDeposit).
		SetInitializerAccount(initializer).
		SetPixelballzMintAccount(pixelBallzMint).
		SetPixelballzTokenAccountAccount(pixelBallzTokenAddress).
		SetQuestAccAccount(questAccount).
		SetQuestAccount(questPda).
		SetQuesteeAccount(questee).
		SetQuestorAccount(questor).
		SetRentAccount(solana.SysVarRentPubkey).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTokenProgramAccount(solana.TokenProgramID)

	if questData.Tender != nil && questData.TenderSplits != nil {
		tenderTokenAccount, _ := utils.GetTokenWallet(initializer, questData.Tender.MintAddress)
		startQuestIx.Append(&solana.AccountMeta{tenderTokenAccount, true, false})
		for _, tenderSplit := range *questData.TenderSplits {
			startQuestIx.Append(&solana.AccountMeta{tenderSplit.TokenAddress, true, false})
		}
	}

	if e := startQuestIx.Validate(); e != nil {
		fmt.Println(e.Error())
		return nil
	}

	return startQuestIx.Build()
}
