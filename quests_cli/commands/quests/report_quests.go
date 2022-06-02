package quests

import (
	"errors"
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"github.com/gagliardetto/solana-go"
)

func ReportQuests(oracle solana.PublicKey, questsPath string) {
	questsData := make([]questing.Quest, 0)

	questsPda, _ := quests.GetQuests(oracle)
	questsPdaData := quests.GetQuestsData(questsPda)
  fmt.Println(questsPda, questsPdaData)
	for i := range make([]int, questsPdaData.Quests) {
		quest, _ := quests.GetQuest(oracle, uint64(i))
		questData := quests.GetQuestData(quest)
		if questData == nil {
			panic(errors.New("bad quest"))
		}

		questsData = append(questsData, *questData)

	}

	WriteQuestsAsMetas(questsData, questsPath)

}
