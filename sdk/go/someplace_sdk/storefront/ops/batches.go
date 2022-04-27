package ops

import (
	"fmt"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func EnableBatches(oracle solana.PrivateKey) {
	batches, _ := storefront.GetBatches(oracle.PublicKey())
	enableIx := someplace.NewEnableBatchUploadingInstructionBuilder().
		SetBatchesAccount(batches).
		SetOracleAccount(oracle.PublicKey()).
		SetSystemProgramAccount(solana.SystemProgramID)

	if e := enableIx.Validate(); e != nil {
		fmt.Println(e.Error())
	}

	utils.SendTx(
		"enable",
		append(make([]solana.Instruction, 0), enableIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)

}

func VerifyBatchUpload(oracle solana.PublicKey) {
	batches, _ := storefront.GetBatches(oracle)
	batchesMeta := storefront.GetBatchesData(batches)
	var i uint64 = 0
	for i < batchesMeta.Counter {
		batchReceipt, _ := storefront.GetBatchReceipt(oracle, i)
		storefront.GetBatchReceiptData(batchReceipt)

		i++
	}

}
