package ops

import (
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"github.com/gagliardetto/solana-go"
)

func EnrollQuestee(oracle, pixelBallzMint, pixelBallzTokenAddress solana.PublicKey) solana.Instruction {
	questor, _ := quests.GetQuestorAccount(oracle)
	questee, _ := quests.GetQuesteeAccount(pixelBallzMint)

	enrollQuestorIx := questing.NewEnrollQuesteeInstructionBuilder().
		SetInitializerAccount(oracle).
		SetPixelballzMintAccount(pixelBallzMint).
		SetPixelballzTokenAccountAccount(pixelBallzTokenAddress).
		SetQuesteeAccount(questee).
		SetQuestorAccount(questor).
		SetSystemProgramAccount(solana.SystemProgramID)

	if e := enrollQuestorIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	return enrollQuestorIx.Build()
}
