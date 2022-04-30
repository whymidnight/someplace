package ops

import (
	"fmt"
	"sync"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"github.com/gagliardetto/solana-go"
)

func PrepareMintsHashMap(oracle solana.PublicKey) map[string][][]someplace.MintHash {
	var globalWg sync.WaitGroup

	batches, _ := storefront.GetBatches(oracle)
	batchesMeta := storefront.GetBatchesData(batches)
	hashList := make([][][]someplace.MintHash, batchesMeta.Counter)
	mintsHashMap := make(map[string][][]someplace.MintHash)
	batchesList := make([]string, batchesMeta.Counter)

	var i uint64 = 0
	for i < batchesMeta.Counter {
		globalWg.Add(1)
		go func(globalWgPtr *sync.WaitGroup, hashListPtr [][][]someplace.MintHash, batchesListPtr []string, batchesCounterPtr uint64) {
			var mintsBatchesWg sync.WaitGroup
			var configIndex uint64 = 0

			batchReceipt, _ := storefront.GetBatchReceipt(oracle, batchesCounterPtr)
			batchReceiptData := storefront.GetBatchReceiptData(batchReceipt)
			batch := batchReceiptData.BatchAccount
			mintsBatches := make([][]someplace.MintHash, batchReceiptData.Items)

			for configIndex < batchReceiptData.Items {
				mintsBatchesWg.Add(1)
				go func(mintsBatchesWgPtr *sync.WaitGroup, mintsBatchesPtr [][]someplace.MintHash, configIndexPtr uint64) {
					var mintsBatchesMintsWg sync.WaitGroup

					listing, _ := storefront.GetListing(oracle, batch, configIndexPtr)
					listingData := storefront.GetListingData(listing)
					if listingData != nil {
						var mintsBatchesMints = make([]someplace.MintHash, listingData.Mints)
						var mintIndex uint64 = 0

						for mintIndex < listingData.Mints {
							mintsBatchesMintsWg.Add(1)

							go func(wg *sync.WaitGroup, mintsPtr []someplace.MintHash, mintIndexPtr uint64) {
								mintHash, _, _ := storefront.GetMintHash(oracle, listing, mintIndexPtr)
								mintHashData := storefront.GetMintHashData(mintHash)
								if mintHashData != nil {
									mintsPtr[mintIndexPtr] = *mintHashData
								}
								wg.Done()
							}(&mintsBatchesMintsWg, mintsBatchesMints, mintIndex)

							mintIndex++
						}
						mintsBatchesPtr[configIndexPtr] = mintsBatchesMints
					}
					mintsBatchesMintsWg.Wait()
					mintsBatchesWgPtr.Done()
				}(&mintsBatchesWg, mintsBatches, configIndex)
				configIndex++
			}
			mintsBatchesWg.Wait()
			batchesListPtr[batchesCounterPtr] = batch.String()
			hashListPtr[batchesCounterPtr] = mintsBatches
			globalWgPtr.Done()
		}(&globalWg, hashList, batchesList, i)
		i++
	}
	globalWg.Wait()

	for i, batch := range hashList {
		mintsHashMap[batchesList[i]] = batch
	}
	fmt.Println(mintsHashMap)

	return mintsHashMap
}
