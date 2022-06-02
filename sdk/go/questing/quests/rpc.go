package quests

import (
	"context"
	"crypto/sha256"

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

func GetQuestAccountData(questAccount solana.PublicKey) *questing.QuestAccount {
	rpcClient := rpc.New(questing.NETWORK)
	bin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), questAccount, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	if bin == nil {
		return nil
	}
	var data questing.QuestAccount
	decoder := ag_binary.NewBorshDecoder(bin.Value.Data.GetBinary())
	err := data.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	return &data

}

func GetQuestAccountsDataForInitializer(initializer solana.PublicKey) []questing.QuestAccount {
	hash := sha256.Sum256([]byte("account:QuestAccount"))
	questAccounts := make([]questing.QuestAccount, 0)
	rpcClient := rpc.New(questing.NETWORK)
	bin, _ := rpcClient.GetProgramAccountsWithOpts(context.TODO(), questing.ProgramID, &rpc.GetProgramAccountsOpts{
		Encoding: "base64",
		Filters: append(
			make([]rpc.RPCFilter, 0),
			rpc.RPCFilter{
				Memcmp: &rpc.RPCFilterMemcmp{
					Offset: 0,
					Bytes:  hash[:8],
				},
			},
			rpc.RPCFilter{
				Memcmp: &rpc.RPCFilterMemcmp{
					Offset: 64,
					Bytes:  initializer.Bytes(),
				},
			},
		),
	})
	if bin == nil {
		return nil
	}
	for _, accountData := range bin {
		var data questing.QuestAccount
		decoder := ag_binary.NewBorshDecoder(accountData.Account.Data.GetBinary())
		err := data.UnmarshalWithDecoder(decoder)
		if err != nil {
			panic(err)
		}
		questAccounts = append(questAccounts, data)
	}

	return questAccounts

}
