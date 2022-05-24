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
	if bin == nil {
		return nil
	}
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
	if bin == nil {
		return nil
	}
	var data questing.Quest
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}

func GetQuestorData(questor solana.PublicKey) *questing.Questor {
	rpcClient := rpc.New(questing.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), questor, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	if bin == nil {
		return nil
	}
	var data questing.Questor
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}
func GetQuesteeData(questee solana.PublicKey) *questing.Questee {
	rpcClient := rpc.New(questing.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), questee, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	if bin == nil {
		return nil
	}
	var data questing.Questee
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}

func GetQuestQuesteeReceiptAccountData(questQuesteeReceipt solana.PublicKey) *questing.QuestQuesteeEndReceipt {
	rpcClient := rpc.New(questing.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), questQuesteeReceipt, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	if bin == nil {
		return nil
	}
	var data questing.QuestQuesteeEndReceipt
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}
