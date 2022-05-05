package quests

import (
	"context"

	ag_binary "github.com/gagliardetto/binary"

	"github.com/gagliardetto/solana-go/rpc"

	"creaturez.nft/questing"
	"github.com/gagliardetto/solana-go"
)

func GetQuestsData(quests solana.PublicKey) *questing.Quests {
	rpcClient := rpc.New(questing.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), quests, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var data questing.Quests
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}

func GetQuestData(quest solana.PublicKey) *questing.Quest {
	rpcClient := rpc.New(questing.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), quest, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var data questing.Quest
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}
