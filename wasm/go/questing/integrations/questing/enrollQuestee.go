package questing

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"creaturez.nft/questing/quests"
	"creaturez.nft/questing/quests/ops"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func EnrollQuestees(this js.Value, args []js.Value) interface{} {
	holder := solana.MustPublicKeyFromBase58(args[0].String())
	var questeeMints []string
	questeeMintsErr := json.Unmarshal([]byte(args[1].String()), &questeeMints)

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
			if len(questeeMints) > 8 {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("too many questee mints - tx would be too large")
				reject.Invoke(errorObject)
				return
			}
			enrollmentJson, err := enrollQuestees(holder, questeeMints)
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

func enrollQuestees(oracle solana.PublicKey, pixelBallzMints []string) ([]byte, error) {
	instructions := make([]solana.Instruction, 0)
	txJson := []byte("{}")

	for _, pixelBallzMint := range pixelBallzMints {
		if enrollQuesteeIx := enrollQuestee(oracle, solana.MustPublicKeyFromBase58(pixelBallzMint)); enrollQuesteeIx != nil {
			instructions = append(
				instructions,
				*enrollQuesteeIx,
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

func enrollQuestee(oracle, pixelBallzMint solana.PublicKey) *solana.Instruction {
	questee, _ := quests.GetQuesteeAccount(pixelBallzMint)
	questeeData := quests.GetQuesteeData(questee)
	if questeeData == nil {
		pixelBallzTokenAddress, _ := utils.GetTokenWallet(oracle, pixelBallzMint)
		enrollmentIx := ops.EnrollQuestee(oracle, pixelBallzMint, pixelBallzTokenAddress)

		return &enrollmentIx
	}
	return nil
}
