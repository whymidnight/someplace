package storefront

import (
	"encoding/json"
	"syscall/js"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront/ops"
	"github.com/gagliardetto/solana-go"
)

func FetchNftHashMap(this js.Value, args []js.Value) interface{} {
	oracle := solana.MustPublicKeyFromBase58(args[0].String())

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			hashMap := ops.PrepareMintsHashMap(oracle)
			viaMints := ops.GetViasHashMap(oracle)
			for via, hashes := range viaMints {
				records := make([][]someplace.MintHash, 1)
				hashList := make([]someplace.MintHash, 0)

				for _, hash := range hashes {
					hashList = append(hashList, hash.MintHash)
				}

				records[0] = hashList
				hashMap[via.String()] = records
			}

			hashMapJson, err := json.Marshal(hashMap)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}
			dst := js.Global().Get("Uint8Array").New(len(hashMapJson))
			js.CopyBytesToJS(dst, hashMapJson)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
