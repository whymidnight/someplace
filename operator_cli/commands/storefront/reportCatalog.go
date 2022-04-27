package storefront

import (
	"errors"
	"os"

	"creaturez.nft/someplace/storefront"
	"github.com/gagliardetto/solana-go"

	"github.com/go-gota/gota/dataframe"
)

// ReportCatalog will report all the listings under `marketUid` and `oracle`.
func ReportCatalog(oracle solana.PublicKey, listingsTableFile string) {
	entireCatalog := make([]Catalog, 0)

	batches, _ := storefront.GetBatches(oracle)
	batchesData := storefront.GetBatchesData(batches)
	if batchesData == nil {
		panic(errors.New("no batches"))
	}

	var i, ii uint64 = 0, 0
	for i < batchesData.Counter {
		batchReceipt, _ := storefront.GetBatchReceipt(oracle, i)
		batchReceiptData := storefront.GetBatchReceiptData(batchReceipt)
		if batchReceiptData == nil {
			i++
			continue
		}
		if batchReceiptData.Items == 0 {
			catalog := Catalog{}
			catalog.BatchAccount = batchReceiptData.BatchAccount.String()
			entireCatalog = append(entireCatalog, catalog)
			i++
			continue
		}

		for ii < batchReceiptData.Items {
			catalog := Catalog{}
			catalog.Mints = 0
			catalog.Resync = false
			catalog.BatchAccount = batchReceiptData.BatchAccount.String()

			listing, _ := storefront.GetListing(oracle, batchReceiptData.BatchAccount, ii)
			listingData := storefront.GetListingData(listing)
			if listingData != nil {
				catalog.IsListed = listingData.IsListed
				catalog.Price = int(listingData.Price)
				catalog.LifecycleStart = int(listingData.LifecycleStart)
				catalog.Mints = int(listingData.Mints)
			}

			catalog.ConfigIndex = int(ii)
			entireCatalog = append(entireCatalog, catalog)

			ii++
		}
		i++
	}

	df := dataframe.LoadStructs(entireCatalog)
	listingsTable, err := os.Create(listingsTableFile)
	if err != nil {
		panic(err)
	}
	defer listingsTable.Close()

	if err = df.WriteCSV(listingsTable); err != nil {
		panic(err)
	}
}