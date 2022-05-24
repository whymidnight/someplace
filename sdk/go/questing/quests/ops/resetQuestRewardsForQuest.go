package ops

import (
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"github.com/gagliardetto/solana-go"
)

func ResetQuestRewardsForQuest(oracle solana.PublicKey, questIndex uint64) solana.Instruction {
	quest, questBump := quests.GetQuest(oracle, questIndex)
	resetQuestIndexIx := questing.NewResetQuestRewardsInstructionBuilder().
		SetOracleAccount(oracle).
		SetQuestAccount(quest).
		SetQuestBump(questBump).
		SetQuestIndex(questIndex)

	if e := resetQuestIndexIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	return resetQuestIndexIx.Build()
}
