package storefront

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"github.com/gagliardetto/solana-go"
)

func GetListings(this js.Value, args []js.Value) interface{} {
	oracle := solana.MustPublicKeyFromBase58(args[0].String())
	batchBatchReceipts := args[1].String()

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			var batchBatchReceiptsSlice []batchReceiptsResponse
			if err := json.Unmarshal([]byte(batchBatchReceipts), &batchBatchReceiptsSlice); err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}

			listingsMap, err := getListings(oracle, batchBatchReceiptsSlice)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}

			listingsJson, err := json.Marshal(listingsMap)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}
			fmt.Println("json", string(listingsJson))
			dst := js.Global().Get("Uint8Array").New(len(listingsJson))
			js.CopyBytesToJS(dst, listingsJson)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func getListings(oracle solana.PublicKey, batchBatchReceipts []batchReceiptsResponse) (map[string][]*someplace.Listing, error) {
	var listings = make(map[string][]*someplace.Listing)
	for _, batch := range batchBatchReceipts {
		listings[batch.BatchReceiptData.BatchAccount.String()] = make([]*someplace.Listing, batch.BatchReceiptData.Items)
		for i := range listings[batch.BatchReceiptData.BatchAccount.String()] {
			listing, _ := storefront.GetListing(oracle, batch.BatchReceiptData.BatchAccount, uint64(i))
			listingData := storefront.GetListingData(listing)
			listings[batch.BatchReceiptData.BatchAccount.String()][i] = listingData
		}
	}

	return listings, nil

}
