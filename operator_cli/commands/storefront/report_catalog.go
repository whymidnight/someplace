package storefront

import (
	"errors"
	"fmt"
	"os"

	"creaturez.nft/someplace/storefront"
	"creaturez.nft/utils"
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

	var i uint64 = 0
	for i < batchesData.Counter {
		var ii uint64 = 0

		fmt.Println(i)
		batchReceipt, _ := storefront.GetBatchReceipt(oracle, i)
		batchReceiptData := storefront.GetBatchReceiptData(batchReceipt)
		if batchReceiptData == nil {
			fmt.Println("asdfasdfasdf")
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

		fmt.Println("items", batchReceiptData.Items)
		for ii < batchReceiptData.Items {
			fmt.Println(batchReceiptData.BatchAccount, ii)
			catalog := Catalog{}
			catalog.Mints = 0
			catalog.Resync = false
			catalog.BatchAccount = batchReceiptData.BatchAccount.String()

			listing, _ := storefront.GetListing(oracle, batchReceiptData.BatchAccount, ii)
			listingData := storefront.GetListingData(listing)
			if listingData != nil {
				treasury, _ := storefront.GetTreasuryAuthority(oracle)
				treasuryData := storefront.GetTreasuryAuthorityData(treasury)
				catalog.IsListed = listingData.IsListed
				catalog.Price = utils.ConvertAmountToUiAmount(listingData.Price, treasuryData.TreasuryDecimals)
				catalog.LifecycleStart = int(listingData.LifecycleStart)
				catalog.Mints = int(listingData.Mints)
			}
			fmt.Println("a.....")

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

