package ops

import (
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func CreateQuest(oracle solana.PrivateKey, questData questing.Quest) {
	questsPda, _ := quests.GetQuests(oracle.PublicKey())
	questsData := quests.GetQuestsData(questsPda)
	quest, _ := quests.GetQuest(oracle.PublicKey(), questsData.Quests)
	createQuestIx := questing.NewCreateQuestInstructionBuilder().
		SetDuration(questData.Duration).
		SetOracleAccount(questData.Oracle).
		SetQuestAccount(quest).
		SetQuestIndex(questsData.Quests).
		SetQuestsAccount(questsPda).
		SetRewards(questData.Rewards).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTender(*questData.Tender).
		SetWlCandyMachines(questData.WlCandyMachines)

	if e := createQuestIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	utils.SendTx(
		"sell",
		append(make([]solana.Instruction, 0), createQuestIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}
