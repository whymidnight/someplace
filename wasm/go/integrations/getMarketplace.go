package integrations

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall/js"

	ag_binary "github.com/gagliardetto/binary"

	"creaturez.nft/wasm/v2/someplace"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func FetchMarketplaceMeta(this js.Value, args []js.Value) interface{} {
	oracle := solana.MustPublicKeyFromBase58(args[0].String())
	marketUid := solana.MustPublicKeyFromBase58(args[1].String())

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {

			marketMeta, err := GetMarketMeta(oracle, marketUid)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}
			dst := js.Global().Get("Uint8Array").New(len(marketMeta))
			js.CopyBytesToJS(dst, marketMeta)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

type TokenListMeta struct {
	Address  solana.PublicKey `json:"address"`
	Symbol   string           `json:"symbol"`
	Name     string           `json:"name"`
	Decimals uint8            `json:"decimals"`
}

type TokenList struct {
	Tokens []TokenListMeta `json:"tokens"`
}

func GetMarketMeta(oracle, marketUid solana.PublicKey) ([]byte, error) {
	marketAuthority, _ := GetMarketAuthority(oracle, marketUid)
	marketAuthorityData := GetMarketAuthorityData(marketAuthority)

	var tokenMeta TokenListMeta
	tokens := FetchTokenMeta()
	for _, token := range tokens {
		if token.Address.Equals(marketAuthorityData.MarketMint) {
			tokenMeta = token
		}
	}

	tokenMetaJson, err := json.Marshal(tokenMeta)
	if err != nil {
		return []byte{}, err
	}

	return tokenMetaJson, nil
}

func FetchTokenMeta() []TokenListMeta {
	var tokenList TokenList
	tokenListUrl := fmt.Sprint(CDN + "/solana.tokenlist.json")
	res, err := http.DefaultClient.Get(tokenListUrl)
	if err != nil {
		return tokenList.Tokens
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return tokenList.Tokens
	}

	err = json.Unmarshal(data, &tokenList)
	if err != nil {
		return tokenList.Tokens
	}

	return tokenList.Tokens

}

func GetMarketAuthority(oracle, marketUid solana.PublicKey) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("someplace"),
			[]byte("market"),
			oracle.Bytes(),
			marketUid.Bytes(),
		},
		someplace.ProgramID,
	)
	return addr, bump
}

func GetMarketAuthorityData(marketAuthority solana.PublicKey) *someplace.Market {
	rpcClient := rpc.New(NETWORK)
	batchReceiptBin, _ := rpcClient.GetAccountInfo(context.TODO(), marketAuthority)
	var batchReceiptData someplace.Market
	fmt.Println(batchReceiptBin.Value)
	if batchReceiptBin != nil {
		decoder := ag_binary.NewBorshDecoder(batchReceiptBin.Value.Data.GetBinary())
		err := batchReceiptData.UnmarshalWithDecoder(decoder)
		if err != nil {
			return nil
		}
	}

	return &batchReceiptData

}

/*
   seeds = [PREFIX.as_ref(), LISTING.as_ref(), market_authority.key().as_ref(), market_authority.listings.to_le_bytes().as_ref()],
   seeds = [PREFIX.as_ref(), LISTINGTOKEN.as_ref(), market_authority.key().as_ref(), index.to_le_bytes().as_ref()],
*/
func GetMarketListing(marketAuthority solana.PublicKey, index uint64) (solana.PublicKey, uint8) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, index)
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("someplace"),
			[]byte("publiclisting"),
			marketAuthority.Bytes(),
			buf,
		},
		someplace.ProgramID,
	)
	return addr, bump
}
func GetMarketListingData(marketListing solana.PublicKey) *someplace.MarketListing {
	var batchReceiptData someplace.MarketListing
	rpcClient := rpc.New(NETWORK)
	batchReceiptBin, _ := rpcClient.GetAccountInfo(context.TODO(), marketListing)
	fmt.Println(batchReceiptBin.Value)
	if batchReceiptBin != nil {
		decoder := ag_binary.NewBorshDecoder(batchReceiptBin.Value.Data.GetBinary())
		err := batchReceiptData.UnmarshalWithDecoder(decoder)
		if err != nil {
			fmt.Println(err)
			return nil
		}

	}

	return &batchReceiptData

}
func GetMarketListingTokenAccount(marketAuthority solana.PublicKey, index uint64) (solana.PublicKey, uint8) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, index)
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("someplace"),
			[]byte("listingtoken"),
			marketAuthority.Bytes(),
			buf,
		},
		someplace.ProgramID,
	)
	return addr, bump
}
