package storefront

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"github.com/gagliardetto/solana-go"

	"github.com/go-gota/gota/dataframe"
)

func reportCatalogFromCandyMachine(candyMachine solana.PublicKey) []ConfigLine {
	cmd := exec.Command("./libs/someplace_rusty", "--candy-machine", candyMachine.String(), "--endpoint", someplace.NETWORK)
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
	}

	var configLines []ConfigLine
	err = json.Unmarshal(data, &configLines)
	if err != nil {
		panic(err)
	}

	return configLines
}

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
				catalog.IsListed = listingData.IsListed
				catalog.Price = int(listingData.Price)
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
