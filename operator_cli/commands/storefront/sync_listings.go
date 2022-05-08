package storefront

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"creaturez.nft/someplace/storefront"
	"creaturez.nft/someplace/storefront/ops"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/go-gota/gota/dataframe"
)

func SyncListingsTable(oracle solana.PrivateKey, listingsTableFile string) {
	listingsTable, err := ioutil.ReadFile(listingsTableFile)
	if err != nil {
		panic(errors.New("bad listings table path in config"))
	}

	df := dataframe.ReadCSV(bytes.NewReader(listingsTable))
	records := df.Records()

	resyncInstructions := make([]solana.Instruction, 0)
	for i, record := range records {
		if i == 0 {
			continue
		}
		_, catalogUpdateInstruction := SyncListingRecord(oracle.PublicKey(), record)
		if catalogUpdateInstruction != nil {
			resyncInstructions = append(resyncInstructions, *catalogUpdateInstruction)
		}
	}
	if len(resyncInstructions) > 0 {
		utils.SendTx(
			"list",
			resyncInstructions,
			append(make([]solana.PrivateKey, 0), oracle),
			oracle.PublicKey(),
		)
	}
	listingsTableLock, err := os.Create(fmt.Sprint(listingsTableFile, ".lock"))
	if err != nil {
		panic(err)
	}
	defer listingsTableLock.Close()

	if err = df.WriteCSV(listingsTableLock); err != nil {
		panic(err)
	}
}

func SyncListingRecord(oracle solana.PublicKey, record []string) (Catalog, *solana.Instruction) {
	catalog := Catalog{}
	catalog.BatchAccount = record[0]

	if a, err := strconv.Atoi(record[1]); err == nil {
		catalog.ConfigIndex = a
	} else {
		panic(err)
	}
	if b, err := strconv.ParseBool(record[2]); err == nil {
		catalog.IsListed = b
	} else {
		panic(err)
	}
	if a, err := strconv.Atoi(record[3]); err == nil {
		// catalog.Price = float64(utils.ConvertUiAmountToAmount(a, treasuryData.TreasuryDecimals))
		catalog.Price = a
	} else {
		panic(err)
	}
	if a, err := strconv.Atoi(record[4]); err == nil {
		catalog.LifecycleStart = a
	} else {
		panic(err)
	}
	if b, err := strconv.ParseBool(record[6]); err == nil {
		catalog.Resync = b
	} else {
		panic(err)
	}

	var resyncInstruction *solana.Instruction = nil
	if catalog.Resync {
		batchAccount := solana.MustPublicKeyFromBase58(catalog.BatchAccount)
		listing, _ := storefront.GetListing(oracle, batchAccount, uint64(catalog.ConfigIndex))
		listingData := storefront.GetListingData(listing)
		if listingData == nil {
			listInstruction := ops.List(oracle, batchAccount, uint64(catalog.ConfigIndex), uint64(catalog.LifecycleStart), uint64(catalog.Price))
			resyncInstruction = &listInstruction
		} else {
			modifyInstruction := ops.Modify(oracle, batchAccount, uint64(catalog.ConfigIndex), uint64(catalog.LifecycleStart), uint64(catalog.Price), catalog.IsListed)
			resyncInstruction = &modifyInstruction
		}

	}

	return catalog, resyncInstruction
}
