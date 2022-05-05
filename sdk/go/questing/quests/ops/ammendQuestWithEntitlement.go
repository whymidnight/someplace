package ops

import (
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func AmmendQuestWithEntitlement(oracle solana.PrivateKey, questData questing.Quest, entitlement questing.Reward) {
	quest, questBump := quests.GetQuest(oracle.PublicKey(), questData.Index)
	entitlementTokenAccount, _ := quests.GetQuestEntitlementTokenAccount(oracle.PublicKey(), questData.Index)
	ammendQuestWithEntitlementIx := questing.NewAmmendQuestWithEntitlementInstructionBuilder().
		SetBallzMintAccount(entitlement.MintAddress).
		SetEntitlement(entitlement).
		SetEntitlementTokenAccountAccount(entitlementTokenAccount).
		SetOracleAccount(oracle.PublicKey()).
		SetQuestAccount(quest).
		SetQuestBump(questBump).
		SetQuestIndex(questData.Index).
		SetRentAccount(solana.SysVarRentPubkey).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTokenProgramAccount(solana.TokenProgramID)

	if e := ammendQuestWithEntitlementIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	utils.SendTx(
		"sell",
		append(make([]solana.Instruction, 0), ammendQuestWithEntitlementIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}
