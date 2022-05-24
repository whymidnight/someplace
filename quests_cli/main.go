package main

import (
	"errors"
	"os"
	"strconv"

	"creaturez.nft/questing"
	"creaturez.nft/questingCLI/v2/commands/quests"
	"creaturez.nft/someplace"
	"github.com/gagliardetto/solana-go"
)

func init() {
	someplace.SetProgramID(solana.MustPublicKeyFromBase58("GXFE4Ym1vxhbXLBx2RxqL5y1Ee3XyFUqDksD7tYjAi8z"))
	questing.SetProgramID(solana.MustPublicKeyFromBase58("Cr4keTx8UQiQ5F9TzTGdQ5dkcMHjxhYSAaHkHhUSABCk"))
}

func main() {
	config := quests.ReadConfig(os.Args[2])
	if config == nil {
		panic(errors.New("bad config"))
	}

	oracle, err := solana.PrivateKeyFromSolanaKeygenFile(config.OraclePath)
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "instance":
		{
			quests.Instance(oracle)
			break
		}
	case "report_quests":
		{
			quests.ReportQuests(oracle.PublicKey(), config.QuestsPath)
			break
		}
	case "sync_quests":
		{
			quests.SyncQuests(oracle, config.QuestsPath)
			quests.ReportQuests(oracle.PublicKey(), config.QuestsPath)
			break
		}
	// resets all the reward entries in a quest.
	// use as way to rewrite the (new) reward
	// entries.
	case "reset_quest_rewards":
		{
			questIndex, err := strconv.Atoi(os.Args[2])
			if err != nil {
				panic(err)
			}
			_ = questIndex
		}
	}
}
