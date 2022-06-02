package storefront

import (
	"encoding/json"
	"os"

	"creaturez.nft/someplace"
	"creaturez.nft/someplace/storefront/ops"
	"github.com/gagliardetto/solana-go"
)

func ReportHashMap(oracle solana.PublicKey, hashMapFile string) {
	mintsHashMap := ops.PrepareMintsHashMap(oracle)
	viaMints := ops.GetViasHashMap(oracle)
	hashMap, err := os.Create(hashMapFile)
	if err != nil {
		panic(err)
	}
	defer hashMap.Close()

	for via, hashes := range viaMints {
		records := make([][]someplace.MintHash, 1)
		hashList := make([]someplace.MintHash, 0)

		for _, hash := range hashes {
			hashList = append(hashList, hash.MintHash)
		}

		records[0] = hashList
		mintsHashMap[via.String()] = records
	}

	js, _ := json.MarshalIndent(mintsHashMap, "", "  ")
	hashMap.Write(js)

}
