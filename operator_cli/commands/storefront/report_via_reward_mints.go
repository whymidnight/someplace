package storefront

import (
	"encoding/json"
	"fmt"

	"creaturez.nft/someplace/storefront/ops"
	"github.com/gagliardetto/solana-go"
)

func ReportViaMintingHashMap(oracle solana.PublicKey) {
	viaMints := ops.GetViasHashMap(oracle)
	viaMintsJs, _ := json.MarshalIndent(viaMints, "", "  ")
	fmt.Println(string(viaMintsJs))
}
