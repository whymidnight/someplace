package questing

import (
	"encoding/json"
	"fmt"
	"strconv"
	"syscall/js"

	"creaturez.nft/questing/quests"
	"creaturez.nft/someplace/storefront/ops"
	"github.com/gagliardetto/solana-go"
)

func DoRNGs(this js.Value, args []js.Value) interface{} {
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
			if len(questeeMints) > 8 {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("too many questee mints - tx would be too large")
				reject.Invoke(errorObject)
				return
			}
			enrollmentJson, err := doRNGs(holder, oracle, questeeMints, uint64(questIndex))
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}

			dst := js.Global().Get("Uint8Array").New(len(enrollmentJson))
			js.CopyBytesToJS(dst, enrollmentJson)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func doRNGs(holder, oracle solana.PublicKey, pixelBallzMints []string, questIndex uint64) ([]byte, error) {
	fmt.Println("looool", holder, oracle, pixelBallzMints, uint64(questIndex))
	instructions := make([]solana.Instruction, 0)
	txJson := []byte("{}")

	// TODO: implement batching of n instructions
	for _, pixelBallzMint := range pixelBallzMints {
		fmt.Println(pixelBallzMint)
		if rngIx := claimQuestBenefits(oracle, holder, solana.MustPublicKeyFromBase58(pixelBallzMint), questIndex); rngIx != nil {
			instructions = append(
				instructions,
				rngIx,
			)
		}
	}

	if len(instructions) > 0 {
		txBuilder := solana.NewTransactionBuilder()
		for _, ix := range instructions {
			txBuilder = txBuilder.AddInstruction(ix)
		}
		txB, _ := txBuilder.Build()
		txJson, _ = json.MarshalIndent(txB, "", "  ")
	}

	fmt.Println(string(txJson))
	return txJson, nil

}

// claimQuestBenefits - accertifies an rng operation for rewarding.
func claimQuestBenefits(oracle, initializer, pixelBallzMint solana.PublicKey, questIndex uint64) solana.Instruction {
	questPda, _ := quests.GetQuest(oracle, questIndex)
	questor, _ := quests.GetQuestorAccount(initializer)
	questee, _ := quests.GetQuesteeAccount(pixelBallzMint)
	questQuesteeReceipt, _ := quests.GetQuestQuesteeReceiptAccount(questor, questee, questPda)
	questQuesteeReceiptData := quests.GetQuestQuesteeReceiptAccountData(questQuesteeReceipt)
	fmt.Println("-.-.-.-.------.-.-.-.-.-.-.--", questPda, questQuesteeReceipt, questQuesteeReceiptData)

	return ops.DoRng(oracle, initializer, pixelBallzMint, questPda, questor, questee, questQuesteeReceiptData.RewardMint)
}
