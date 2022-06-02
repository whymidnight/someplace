package storefront

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"github.com/gagliardetto/solana-go"
)

type batchReceiptsResponse struct {
	BatchReceipt     solana.PublicKey
	BatchReceiptData someplace.BatchReceipt
}

func FetchNfts(this js.Value, args []js.Value) interface{} {
	oracle := solana.MustPublicKeyFromBase58(args[0].String())

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {

			fmt.Println("asd help me")
			batches, _ := storefront.GetBatches(oracle)
			batchesData := storefront.GetBatchesData(batches)
			batchReceipts := make([]batchReceiptsResponse, batchesData.Counter)
			for index := range make([]uint64, batchesData.Counter) {
				batchReceipt, _ := storefront.GetBatchReceipt(oracle, uint64(index))
				batchReceiptData := storefront.GetBatchReceiptData(batchReceipt)

				response := batchReceiptsResponse{
					BatchReceipt:     batchReceipt,
					BatchReceiptData: *batchReceiptData,
				}
				fmt.Println("....", response)
				batchReceipts[index] = response
			}

			batchReceiptsJson, err := json.Marshal(batchReceipts)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}
			fmt.Println("json", string(batchReceiptsJson))
			dst := js.Global().Get("Uint8Array").New(len(batchReceiptsJson))
			js.CopyBytesToJS(dst, batchReceiptsJson)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
