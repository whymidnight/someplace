package ops

import (
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"github.com/gagliardetto/solana-go"
)

func AppendQuestRewards(oracle solana.PublicKey, questData questing.Quest) []solana.Instruction {
	questsPda, _ := quests.GetQuests(oracle)
	questsData := quests.GetQuestsData(questsPda)
    fmt.Println(questsData.Quests)
	quest, questBump := quests.GetQuest(oracle, questData.Index)

	appendRewardIxs := make([]solana.Instruction, 0)

	rewards := questData.Rewards
	for _, reward := range rewards {
		appendRewardIx := questing.NewRegisterQuestRewardInstructionBuilder().
			SetOracleAccount(oracle).
			SetQuestAccount(quest).
			SetQuestBump(questBump).
			SetQuestIndex(questData.Index).
			SetRentAccount(solana.SysVarRentPubkey).
			SetReward(reward).
			SetRewardMintAccount(reward.MintAddress).
			SetSystemProgramAccount(solana.SystemProgramID).
			SetTokenProgramAccount(solana.TokenProgramID)

		if e := appendRewardIx.Validate(); e != nil {
			fmt.Println(e.Error())
			panic("...")
		}

		appendRewardIxs = append(appendRewardIxs, appendRewardIx.Build())

	}

	return appendRewardIxs
}
