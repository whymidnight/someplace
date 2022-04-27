package marketplace

import (
	"context"

	ag_binary "github.com/gagliardetto/binary"

	"creaturez.nft/someplace"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetMarketListingData(marketListing solana.PublicKey) *someplace.MarketListing {
	rpcClient := rpc.New(someplace.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), marketListing, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var data someplace.MarketListing
	if bin != nil {
		decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
		err := data.UnmarshalWithDecoder(decoder)
		if err != nil {
			panic(err)
		}
	}

	return &data

}

func GetMarketAuthorityData(marketAuthority solana.PublicKey) *someplace.Market {
	rpcClient := rpc.New(someplace.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), marketAuthority, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var data someplace.Market
	if bin != nil {
		decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
		err := data.UnmarshalWithDecoder(decoder)
		if err != nil {
			panic(err)
		}
	}

	return &data

}
