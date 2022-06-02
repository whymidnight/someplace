package questing

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"syscall/js"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"github.com/gagliardetto/solana-go"
)

func GetRewards(this js.Value, args []js.Value) interface{} {
	holder := solana.MustPublicKeyFromBase58(args[0].String())
	var questeeMints []string
	questeeMintsErr := json.Unmarshal([]byte(args[1].String()), &questeeMints)
	oracle := solana.MustPublicKeyFromBase58(args[2].String())
	questIndex, questIndexError := strconv.Atoi(args[3].String())

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			if questeeMintsErr != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("invalid questee mints")
				reject.Invoke(errorObject)
				return
			}
			if questIndexError != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("invalid quest index")
				reject.Invoke(errorObject)
				return
			}
			rewardTicketsJson, err := getRewardTickets(holder, oracle, questeeMints, uint64(questIndex))
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}

			dst := js.Global().Get("Uint8Array").New(len(rewardTicketsJson))
			js.CopyBytesToJS(dst, rewardTicketsJson)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func getRewardTickets(holder, oracle solana.PublicKey, pixelBallzMints []string, questIndex uint64) ([]byte, error) {
	rewardTicketsJson := []byte("{}")
	type RewardTicketResponse struct {
		RewardTicket           *someplace.RewardTicket
		QuestQuesteeEndReceipt *questing.QuestQuesteeEndReceipt
	}
	type RewardTicketsResponse map[string]*RewardTicketResponse
	rewardTickets := make(RewardTicketsResponse)

	for _, pixelBallzMint := range pixelBallzMints {
		_, rewardTicketData, questQuesteeReceiptData, err := getRewardTicket(oracle, holder, solana.MustPublicKeyFromBase58(pixelBallzMint), questIndex)
		if err != nil {
			fmt.Println(err.Error())
			rewardTickets[pixelBallzMint] = nil
			continue
		}
		rewardTickets[pixelBallzMint] = &RewardTicketResponse{
			RewardTicket:           rewardTicketData,
			QuestQuesteeEndReceipt: questQuesteeReceiptData,
		}
	}

	if len(pixelBallzMints) > 0 {
		rewardTicketsJson, _ = json.MarshalIndent(rewardTickets, "", "  ")
	}

	fmt.Println(string(rewardTicketsJson))
	return rewardTicketsJson, nil

}

func getRewardTicket(oracle, initializer, pixelBallzMint solana.PublicKey, questIndex uint64) (solana.PublicKey, *someplace.RewardTicket, *questing.QuestQuesteeEndReceipt, error) {
	questPda, _ := quests.GetQuest(oracle, questIndex)
	questor, _ := quests.GetQuestorAccount(initializer)
	questee, _ := quests.GetQuesteeAccount(pixelBallzMint)
	questQuesteeReceipt, _ := quests.GetQuestQuesteeReceiptAccount(questor, questee, questPda)
	questQuesteeReceiptData := quests.GetQuestQuesteeReceiptAccountData(questQuesteeReceipt)
	if questQuesteeReceiptData == nil {
		return solana.PublicKey{}, nil, nil, errors.New("bad quest questee end receipt")
	}
	viaMap, _ := storefront.GetViaMapping(oracle, questQuesteeReceiptData.RewardMint)
	viaMapData := storefront.GetViaMappingData(viaMap)
	if viaMapData == nil {
		return solana.PublicKey{}, nil, nil, errors.New("bad via map")
	}
	via, _ := storefront.GetVia(oracle, viaMapData.ViasIndex)
	rewardTicket, _ := storefront.GetRewardTicket(via, questPda, questee, initializer)
	rewardTicketData := storefront.GetRewardTicketData(rewardTicket)

	return rewardTicket, rewardTicketData, questQuesteeReceiptData, nil
}

