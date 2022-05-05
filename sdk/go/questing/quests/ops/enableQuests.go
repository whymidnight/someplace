package ops

import (
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func EnableQuests(oracle solana.PrivateKey) {
	quests, _ := quests.GetQuests(oracle.PublicKey())
	enableQuestIx := questing.NewEnableQuestsInstructionBuilder().
		SetOracleAccount(oracle.PublicKey()).
		SetQuestsAccount(quests).
		SetSystemProgramAccount(solana.SystemProgramID)

	if e := enableQuestIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	utils.SendTx(
		"sell",
		append(make([]solana.Instruction, 0), enableQuestIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}
