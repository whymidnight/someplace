package storefront

import (
	"context"
	"fmt"

	ag_binary "github.com/gagliardetto/binary"

	"github.com/gagliardetto/solana-go/rpc"

	"creaturez.nft/someplace"
	"github.com/gagliardetto/solana-go"
)

func GetTreasuryWhitelistData(treasuryAuthority solana.PublicKey) *someplace.TreasuryWhitelist {
	rpcClient := rpc.New(someplace.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), treasuryAuthority, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var data someplace.TreasuryWhitelist
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}
func GetTreasuryAuthorityData(treasuryAuthority solana.PublicKey) *someplace.TreasuryAuthority {
	fmt.Println(someplace.ProgramID)
	rpcClient := rpc.New(someplace.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), treasuryAuthority, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	if bin == nil {
		return nil
	}
	var data someplace.TreasuryAuthority
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}
func GetBatchesData(batches solana.PublicKey) *someplace.Batches {
	rpcClient := rpc.New(someplace.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), batches, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	if bin == nil {
		return nil
	}
	var data someplace.Batches
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}

func GetBatchReceiptData(batchReceipt solana.PublicKey) *someplace.BatchReceipt {
	rpcClient := rpc.New(someplace.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), batchReceipt, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	if bin == nil {
		return nil
	}
	var data someplace.BatchReceipt
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}
func GetListingData(listing solana.PublicKey) *someplace.Listing {
	rpcClient := rpc.New(someplace.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), listing, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	if bin == nil {
		return nil
	}
	var data someplace.Listing
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}
func GetMintHashData(mintHash solana.PublicKey) *someplace.MintHash {
	rpcClient := rpc.New(someplace.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), mintHash, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	if bin == nil {
		return nil
	}
	var data someplace.MintHash
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}
