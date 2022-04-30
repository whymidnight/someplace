package storefront

import (
	"encoding/json"
	"os"

	"creaturez.nft/someplace/storefront/ops"
	"github.com/gagliardetto/solana-go"
)

func ReportHashMap(oracle solana.PublicKey, hashMapFile string) {
	mintsHashMap := ops.PrepareMintsHashMap(oracle)
	hashMap, err := os.Create(hashMapFile)
	if err != nil {
		panic(err)
	}
	defer hashMap.Close()
	js, _ := json.MarshalIndent(mintsHashMap, "", "  ")
	hashMap.Write(js)

}
