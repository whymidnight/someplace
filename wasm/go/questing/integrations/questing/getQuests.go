package questing

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"github.com/gagliardetto/solana-go"
)

func GetQuests(this js.Value, args []js.Value) interface{} {
	oracle := solana.MustPublicKeyFromBase58(args[0].String())

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {

			quests := getQuests(oracle)

			questsJSON, err := json.Marshal(quests)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}
			fmt.Println("json", string(questsJSON))
			dst := js.Global().Get("Uint8Array").New(len(questsJSON))
			js.CopyBytesToJS(dst, questsJSON)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func getQuests(oracle solana.PublicKey) map[solana.PublicKey]questing.Quest {

	questsData := make(map[solana.PublicKey]questing.Quest)

	questsPda, _ := quests.GetQuests(oracle)
	questsPdaData := quests.GetQuestsData(questsPda)
  if questsPdaData == nil {
    return questsData
  }
	for i := range make([]int, questsPdaData.Quests) {
		quest, _ := quests.GetQuest(oracle, uint64(i))
		questData := quests.GetQuestData(quest)
		if questData == nil {
			panic(errors.New("bad quest"))
		}

		questsData[quest] = *questData

	}

	return questsData
}
