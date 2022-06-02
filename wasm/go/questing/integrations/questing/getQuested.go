package questing

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"creaturez.nft/questing/quests/ops"
	"github.com/gagliardetto/solana-go"
)

func GetQuested(this js.Value, args []js.Value) interface{} {
	oracle := solana.MustPublicKeyFromBase58(args[0].String())
	holder := solana.MustPublicKeyFromBase58(args[1].String())

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {

			quested := getQuested(oracle, holder)
			questedJSON, err := json.Marshal(quested)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}
			fmt.Println("json", string(questedJSON))
			dst := js.Global().Get("Uint8Array").New(len(questedJSON))
			js.CopyBytesToJS(dst, questedJSON)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func getQuested(oracle, holder solana.PublicKey) ops.QuestedMetaMap {
	return ops.GetQuested(oracle, holder)
}
