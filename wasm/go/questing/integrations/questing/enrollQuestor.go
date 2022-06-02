package questing

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"creaturez.nft/questing/quests"
	"creaturez.nft/questing/quests/ops"
	"github.com/gagliardetto/solana-go"
)

func EnrollQuestor(this js.Value, args []js.Value) interface{} {
	holder := solana.MustPublicKeyFromBase58(args[0].String())

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			enrollmentJson, err := enrollQuestor(holder)
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

func enrollQuestor(oracle solana.PublicKey) ([]byte, error) {
	instructions := make([]solana.Instruction, 0)
	txJson := []byte("{}")

	questor, _ := quests.GetQuestorAccount(oracle)
	questorData := quests.GetQuestorData(questor)
	if questorData == nil {
		enrollmentIx := ops.EnrollQuestor(oracle)
		instructions = append(instructions, enrollmentIx)

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
