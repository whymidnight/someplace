package ops

import (
	"encoding/json"
	"fmt"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func ReportCardinalities(oracle solana.PrivateKey, batchAccount solana.PublicKey, cardinalitiesIndices [][]uint64, cardinalitiesKeys []string) {
	batchCardinalitiesReport, _ := storefront.GetBatchCardinalitiesReport(batchAccount)
	batchCardinalitiesReportData := storefront.GetBatchCardinalitiesReportData(batchCardinalitiesReport)
	if batchCardinalitiesReportData != nil {
		js, _ := json.MarshalIndent(batchCardinalitiesReportData, "", "  ")
		fmt.Println(batchCardinalitiesReport, string(js))
		return
	}

	reportCardinalitiesIx := someplace.NewReportBatchCardinalitiesInstructionBuilder().
		SetBatchAccount(batchAccount).
		SetBatchCardinalitiesReportAccount(batchCardinalitiesReport).
		SetCardinalitiesIndices(cardinalitiesIndices).
		SetCardinalitiesKeys(cardinalitiesKeys).
		SetOracleAccount(oracle.PublicKey()).
		SetSystemProgramAccount(solana.SystemProgramID)

	if e := reportCardinalitiesIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	utils.SendTx(
		"sell",
		append(make([]solana.Instruction, 0), reportCardinalitiesIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}
