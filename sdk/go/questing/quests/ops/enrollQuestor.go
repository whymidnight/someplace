package ops

import (
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"github.com/gagliardetto/solana-go"
)

func EnrollQuestor(oracle, pixelBallzMint solana.PublicKey) solana.Instruction {
	questor, _ := quests.GetQuestorAccount(oracle)

	enrollQuestorIx := questing.NewEnrollQuestorInstructionBuilder().
		SetInitializerAccount(oracle).
		SetQuestorAccount(questor).
		SetSystemProgramAccount(solana.SystemProgramID)

	if e := enrollQuestorIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	return enrollQuestorIx.Build()
}
