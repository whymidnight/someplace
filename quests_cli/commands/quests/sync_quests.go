package quests

import (
	"errors"
	"fmt"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests/ops"
	"creaturez.nft/someplace"
	storefront_ops "creaturez.nft/someplace/storefront/ops"
	"creaturez.nft/utils"
	"github.com/gagliardetto/solana-go"
)

func SyncQuests(oracle solana.PrivateKey, questsPath string) {
	questsMap, questsMetas, questsMetasCreate := ReadQuestsMetas(questsPath)
	if len(questsMap) == 0 {
		panic(errors.New("no quests"))
	}

	signers := make([]solana.PrivateKey, 0)
	for _, quest := range *questsMetasCreate {
		fmt.Println(",....")
		resyncInstructions := make([]solana.Instruction, 0)
		_ = resyncInstructions
		rewards := make([]questing.Reward, 0)
		for _, reward := range quest.Rewards {
			mintKey := solana.NewWallet().PrivateKey
			signers = append(signers, mintKey)
			rewards = append(rewards, questing.Reward{
				MintAddress:  mintKey.PublicKey(),
				RngThreshold: reward.RngThreshold,
				Amount:       reward.Amount,
				Cardinality:  reward.Cardinality,
			})
		}

		tender := func() *questing.Tender {
			if quest.Tender == nil {
				return nil
			}
			tender := questing.Tender{
				MintAddress: quest.Tender.MintAddress,
				Amount:      quest.Tender.Amount,
			}
			return &tender
		}()

		tenderSplits := func() *[]questing.Split {
			if quest.TenderSplits == nil {
				return nil
			}
			tenderSplits := make([]questing.Split, 0)
			for _, tenderSplit := range *quest.TenderSplits {
				tenderSplits = append(tenderSplits, questing.Split{
					TokenAddress: tenderSplit.TokenAddress,
					OpCode:       tenderSplit.OpCode,
					Share:        tenderSplit.Share,
				})
			}
			return &tenderSplits
		}()

		questData := questing.Quest{
			Index:           quest.Index,
			Name:            quest.Name,
			Duration:        quest.Duration,
			Oracle:          quest.Oracle,
			WlCandyMachines: quest.WlCandyMachines,
			Entitlement:     &questing.Reward{},
			Rewards:         rewards,
			Tender:          tender,
			TenderSplits:    tenderSplits,
			Xp:              quest.Xp,
		}
		createQuestIx, questIndex := ops.CreateQuest(oracle.PublicKey(), questData)
		questData.Index = questIndex
		resyncInstructions = append(
			resyncInstructions,
			createQuestIx,
		)

		rewardIxs := ops.AppendQuestRewards(oracle.PublicKey(), questData)
		resyncInstructions = append(resyncInstructions, rewardIxs...)
		utils.SendTx(
			"list",
			resyncInstructions,
			append(signers, oracle),
			oracle.PublicKey(),
		)

	}
	for _, quest := range *questsMetas {
		if !quest.Resync {
			continue
		}
		questData, additionalSigners := quest.to_questing_quest()
		if len(additionalSigners) > 0 {
			fmt.Println(additionalSigners)
			signers = append(signers, additionalSigners...)
		}
		resyncInstructions := make([]solana.Instruction, 0)

		resetQuestRewardsIx := ops.ResetQuestRewardsForQuest(oracle.PublicKey(), questData.Index)
		resyncInstructions = append(resyncInstructions, resetQuestRewardsIx)

		rewardIxs := ops.AppendQuestRewards(oracle.PublicKey(), questData)
		resyncInstructions = append(resyncInstructions, rewardIxs...)

		if len(resyncInstructions) > 0 {
			viaIxs := storefront_ops.EnableViasForRarityTokens(oracle.PublicKey(), func() []someplace.ViaMint {
				viaMints := make([]someplace.ViaMint, 0)
				for _, reward := range questData.Rewards {
					fmt.Println(reward)
					viaMints = append(viaMints, someplace.ViaMint{
						MintAddress: reward.MintAddress,
						Rarity:      *reward.Cardinality,
					})
				}
				return viaMints
			}())
			resyncInstructions = append(resyncInstructions, viaIxs...)
			utils.SendTx(
				"list",
				resyncInstructions,
				append(signers, oracle),
				oracle.PublicKey(),
			)
		}
	}
}
