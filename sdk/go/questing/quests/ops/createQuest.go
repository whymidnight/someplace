package ops

import (
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"github.com/gagliardetto/solana-go"
)

func CreateQuest(oracle solana.PublicKey, questData questing.Quest) (solana.Instruction, uint64) {
	questsPda, _ := quests.GetQuests(oracle)
	questsData := quests.GetQuestsData(questsPda)
	quest, _ := quests.GetQuest(oracle, questsData.Quests)
	createQuestIx := questing.NewCreateQuestInstructionBuilder().
		SetDuration(questData.Duration).
		SetName(questData.Name).
		SetOracleAccount(oracle).
		SetQuestAccount(quest).
		SetQuestIndex(questsData.Quests).
		SetQuestsAccount(questsPda).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetWlCandyMachines(questData.WlCandyMachines).
		SetXp(questData.Xp).
		SetEnabled(true).
		SetRequiredLevel(questData.RequiredLevel)

	if questData.Tender != nil {
		createQuestIx.SetTender(*questData.Tender)
		createQuestIx.SetTenderSplits(*questData.TenderSplits)
	}

	if e := createQuestIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	return createQuestIx.Build(), questsData.Quests
}
