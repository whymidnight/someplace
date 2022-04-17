package integrations

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"syscall/js"

	ag_binary "github.com/gagliardetto/binary"

	"creaturez.nft/wasm/v2/someplace"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var programId = solana.MustPublicKeyFromBase58("5WwhzMCFSgWYxiuKrbsB9wtg9T49Mm1fD1v2UdhD5oYi")

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

			batches, _ := GetBatches(oracle)
			fmt.Println(oracle.String(), batches, programId.String())
			batchesData, err := GetBatchesData(batches)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}
			fmt.Println(batchesData.Counter)
			batchReceipts := make([]batchReceiptsResponse, batchesData.Counter)
			for index := range make([]uint64, batchesData.Counter) {
				batchReceipt, _ := GetBatchReceipt(oracle, uint64(index))
				batchReceiptData, err := GetBatchReceiptData(batchReceipt)
				if err != nil {
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New("unauthorized")
					reject.Invoke(errorObject)
					return
				}
				response := batchReceiptsResponse{
					BatchReceipt:     batchReceipt,
					BatchReceiptData: batchReceiptData,
				}
				fmt.Println("....", response)
				batchReceipts[index] = response
			}

			fmt.Println("----", batchReceipts)
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

func catalogBatches(batches solana.PublicKey) {

}

func GetBatchReceipt(
	oracle solana.PublicKey,
	index uint64,
) (solana.PublicKey, uint8) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, index)
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			oracle.Bytes(),
			buf,
		},
		programId,
	)
	return addr, bump
}

func GetBatches(
	oracle solana.PublicKey,
) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			oracle.Bytes(),
		},
		programId,
	)
	return addr, bump
}

func GetBatchesData(batches solana.PublicKey) (someplace.Batches, error) {
	rpcClient := rpc.New("https://sparkling-dark-shadow.solana-devnet.quiknode.pro/0e9964e4d70fe7f856e7d03bc7e41dc6a2b84452/")
	batchesBin, _ := rpcClient.GetAccountInfo(context.TODO(), batches)
	var batchesData someplace.Batches
	decoder := ag_binary.NewBorshDecoder(batchesBin.Value.Data.GetBinary())
	err := batchesData.UnmarshalWithDecoder(decoder)
	if err != nil {
		return batchesData, err
	}

	return batchesData, nil

}

func GetBatchReceiptData(batchReceipt solana.PublicKey) (someplace.BatchReceipt, error) {
	rpcClient := rpc.New("https://sparkling-dark-shadow.solana-devnet.quiknode.pro/0e9964e4d70fe7f856e7d03bc7e41dc6a2b84452/")
	batchReceiptBin, _ := rpcClient.GetAccountInfo(context.TODO(), batchReceipt)
	var batchReceiptData someplace.BatchReceipt
	decoder := ag_binary.NewBorshDecoder(batchReceiptBin.Value.Data.GetBinary())
	err := batchReceiptData.UnmarshalWithDecoder(decoder)
	if err != nil {
		return batchReceiptData, err
	}

	return batchReceiptData, nil

}
