package storefront

import (
	"fmt"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"creaturez.nft/someplace/storefront/ops"
	"github.com/gagliardetto/solana-go"
)

func ReportCardinalities(oracle solana.PrivateKey) {
	/*
	   Required ops include,

	   enable batching for max(u64) candy machines
	   initialize a treasury for storefront under `oracle`
	*/
	_ = someplace.NETWORK
	fmt.Println(someplace.NETWORK)
	batches, _ := storefront.GetBatches(oracle.PublicKey())
	batchesData := storefront.GetBatchesData(batches)
	for i := range make([]int, batchesData.Counter) {
		batchReceipt, _ := storefront.GetBatchReceipt(oracle.PublicKey(), uint64(i))
		batchReceiptData := storefront.GetBatchReceiptData(batchReceipt)
		batchAccount := batchReceiptData.BatchAccount
		candyMachine := batchAccount
		catalog := reportCatalogFromCandyMachine(candyMachine)

		cardinalityMap := make(map[string][]uint64)
		keys := make([]string, 0)
		for _, configLine := range catalog {
			cardinalityMap[configLine.Cardinality] = make([]uint64, 0)
		}

		for k := range cardinalityMap {
			keys = append(keys, k)
		}

		for i, configLine := range catalog {
			cardinalityMap[configLine.Cardinality] = append(
				cardinalityMap[configLine.Cardinality],
				uint64(i),
			)
		}

		fmt.Println(cardinalityMap)
		fmt.Println(keys)
		cardinalityIndices := make([][]uint64, len(keys))
		cardinalityKeys := make([]string, len(keys))
		for i, key := range keys {
			cardinalityKeys[i] = key
			cardinalityIndices[i] = cardinalityMap[key]
		}

		ops.ReportCardinalities(oracle, batchAccount, cardinalityIndices, cardinalityKeys)
	}
}
