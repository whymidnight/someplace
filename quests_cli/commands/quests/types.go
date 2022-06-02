package quests

import (
	"fmt"
	"io/ioutil"

	"creaturez.nft/questing"
	"github.com/gagliardetto/solana-go"
	ag_solanago "github.com/gagliardetto/solana-go"
	"gopkg.in/yaml.v2"
)

type QuestMetaReward struct {
	MintAddress  *ag_solanago.PublicKey `yaml:"MintAddress"`
	RngThreshold *uint8                  `yaml:"RngThreshold"`
	Amount       uint64                 `yaml:"Amount"`
	Cardinality  *string                `yaml:"Cardinality"`
}
type QuestMetaTender struct {
	MintAddress ag_solanago.PublicKey `yaml:"MintAddress"`
	Amount      uint64                `yaml:"Amount"`
}
type QuestMetaSplit struct {
	TokenAddress ag_solanago.PublicKey `yaml:"TokenAddress"`
	OpCode       uint8                 `yaml:"OpCode"`
	Share        uint8                 `yaml:"Share"`
}

type QuestMeta struct {
	Enabled         bool                    `yaml:"Enabled"`
	Index           uint64                  `yaml:"Index"`
	Name            string                  `yaml:"Name"`
	Duration        int64                   `yaml:"Duration"`
	Oracle          ag_solanago.PublicKey   `yaml:"Oracle"`
	RequiredLevel   uint64                  `yaml:"RequiredLevel"`
	WlCandyMachines []ag_solanago.PublicKey `yaml:"WlCandyMachines"`
	Entitlement     *QuestMetaReward        `yaml:"Entitlement"`
	Rewards         []QuestMetaReward       `yaml:"Rewards"`
	Tender          *QuestMetaTender        `yaml:"Tender"`
	TenderSplits    *[]QuestMetaSplit       `yaml:"TenderSplits"`
	Xp              uint64                  `yaml:"Xp"`
	Resync          bool                    `yaml:"Resync"`
}

type ConfigYaml struct {
	OraclePath string `yaml:"OraclePath"`
	QuestsPath string `yaml:"QuestsPath"`
}

func ReadConfig(configPath string) *ConfigYaml {
	var config = new(ConfigYaml)
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(configFile, config); err != nil {
		panic(err)
	}

	return config
}

func ReadQuestsMeta(questsPath string) QuestMeta {
	var questMeta QuestMeta
	questsFile, err := ioutil.ReadFile(questsPath)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(questsFile, &questMeta); err != nil {
		panic(err)
	}

	return questMeta
}

func ReadQuestsMetas(questsPath string) ([]map[string]string, *[]QuestMeta, *[]QuestMeta) {
	var questsMetas = make([]QuestMeta, 0)
	var questsMetasCreate = make([]QuestMeta, 0)
	questsFile, err := ioutil.ReadFile(questsPath)
	if err != nil {
		panic(err)
	}
	questsMap := make(map[string][]map[string]string)
	if err = yaml.Unmarshal(questsFile, &questsMap); err != nil {
		panic(err)
	}
	for _, quest := range questsMap["quests"] {
		for k, v := range quest {
			if k == "upload" {
				// process creation

				questCreationFiles := func() []string {
					files := make([]string, 0)
					creationFiles, err := ioutil.ReadDir(v)
					if err != nil {
						panic(err)
					}
					for _, creationFile := range creationFiles {
						name := creationFile.Name()
						fmt.Println(name[len(name)-3:])
						if name[len(name)-3:] == "yml" {
							files = append(files, fmt.Sprint(v, creationFile.Name()))
						}
					}
					return files
				}()

				for _, questCreationFile := range questCreationFiles {
					meta := ReadQuestsMeta(questCreationFile)
					questsMetasCreate = append(questsMetasCreate, meta)
				}
				fmt.Println(questCreationFiles)
			} else {
				// read from cache
				meta := ReadQuestsMeta(v)
				questsMetas = append(questsMetas, meta)
			}
		}
	}

	return questsMap["quests"], &questsMetas, &questsMetasCreate
}

func (quest QuestMeta) to_questing_quest() (questing.Quest, []solana.PrivateKey) {
	additionalSigners := make([]solana.PrivateKey, 0)
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

	entitlement := func() *questing.Reward {
		if quest.Entitlement == nil {
			return nil
		}
		entitlement := questing.Reward{
			MintAddress: *quest.Entitlement.MintAddress,
			Amount:      quest.Entitlement.Amount,
		}
		return &entitlement
	}()

	rewards := func() []questing.Reward {
		if quest.Rewards == nil {
			return nil
		}
		rewards := make([]questing.Reward, 0)
		for _, reward := range quest.Rewards {
			mint := solana.NewWallet().PrivateKey
			var mintAddress *solana.PublicKey = reward.MintAddress
			if mintAddress == nil {
				mintAddress = mint.PublicKey().ToPointer()
				additionalSigners = append(additionalSigners, mint)
			}
			rewards = append(rewards, questing.Reward{
				MintAddress:  *mintAddress,
				Amount:       reward.Amount,
				Cardinality:  reward.Cardinality,
				RngThreshold: *reward.RngThreshold,
			})
		}
		return rewards
	}()
	return questing.Quest{
		Index:           quest.Index,
		Name:            quest.Name,
		Duration:        quest.Duration,
		Oracle:          quest.Oracle,
		WlCandyMachines: quest.WlCandyMachines,
		Entitlement:     entitlement,
		Rewards:         rewards,
		Tender:          tender,
		TenderSplits:    tenderSplits,
		Xp:              quest.Xp,
	}, additionalSigners
}

func WriteQuestsAsMetas(questsData []questing.Quest, questsPath string) {
	questsMap := make(map[string][]map[string]string)
	questsMap["quests"] = make([]map[string]string, 0)

	for _, questData := range questsData {

		questFilePath := fmt.Sprint("./quests/", questData.Index, "_", questData.Name, ".yml")
		questMap := make(map[string]string)
		questMap[fmt.Sprint(questData.Index)] = questFilePath
		questsMap["quests"] = append(questsMap["quests"], questMap)

		questYamlData := QuestMeta{
			Index:           questData.Index,
			Name:            questData.Name,
			Duration:        questData.Duration,
			Oracle:          questData.Oracle,
			RequiredLevel:   questData.RequiredLevel,
			WlCandyMachines: questData.WlCandyMachines,
			Entitlement: func() *QuestMetaReward {
				if questData.Entitlement == nil {
					return nil
				}
				return &QuestMetaReward{
					MintAddress:  &questData.Entitlement.MintAddress,
					Amount:       questData.Entitlement.Amount,
				}
			}(),
			Rewards: func() []QuestMetaReward {
				questMetaRewards := make([]QuestMetaReward, 0)
				for _, reward := range questData.Rewards {
					questMetaRewards = append(questMetaRewards, QuestMetaReward{
						MintAddress:  reward.MintAddress.ToPointer(),
						RngThreshold: &reward.RngThreshold,
						Amount:       reward.Amount,
						Cardinality:  reward.Cardinality,
					})
				}
				return questMetaRewards
			}(),
			Tender: func() *QuestMetaTender {
				if questData.Tender == nil {
					return nil
				}
				return &QuestMetaTender{
					MintAddress: questData.Tender.MintAddress,
					Amount:      questData.Tender.Amount,
				}
			}(),
			TenderSplits: func() *[]QuestMetaSplit {
				if questData.TenderSplits == nil {
					return nil
				}
				questMetaTenderSplits := make([]QuestMetaSplit, 0)
				for _, tender := range *questData.TenderSplits {
					questMetaTenderSplits = append(questMetaTenderSplits, QuestMetaSplit{
						TokenAddress: tender.TokenAddress,
						OpCode:       tender.OpCode,
						Share:        tender.Share,
					})
				}
				return &questMetaTenderSplits
			}(),
			Xp:     questData.Xp,
			Resync: false,
		}

		questYaml, err := yaml.Marshal(questYamlData)
		if err != nil {
			panic(err)
		}
		if err = ioutil.WriteFile(questFilePath, questYaml, 0644); err != nil {
			panic(err)
		}
	}

	questMapYaml, err := yaml.Marshal(questsMap)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(questsPath, questMapYaml, 0644); err != nil {
		panic(err)
	}
}
