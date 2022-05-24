// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package questing

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// RegisterQuestReward is the `registerQuestReward` instruction.
type RegisterQuestReward struct {
	QuestIndex *uint64
	QuestBump  *uint8
	Reward     *Reward

	// [0] = [WRITE, SIGNER] oracle
	//
	// [1] = [WRITE] quest
	//
	// [2] = [WRITE, SIGNER] rewardMint
	//
	// [3] = [] systemProgram
	//
	// [4] = [] tokenProgram
	//
	// [5] = [] rent
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewRegisterQuestRewardInstructionBuilder creates a new `RegisterQuestReward` instruction builder.
func NewRegisterQuestRewardInstructionBuilder() *RegisterQuestReward {
	nd := &RegisterQuestReward{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 6),
	}
	return nd
}

// SetQuestIndex sets the "questIndex" parameter.
func (inst *RegisterQuestReward) SetQuestIndex(questIndex uint64) *RegisterQuestReward {
	inst.QuestIndex = &questIndex
	return inst
}

// SetQuestBump sets the "questBump" parameter.
func (inst *RegisterQuestReward) SetQuestBump(questBump uint8) *RegisterQuestReward {
	inst.QuestBump = &questBump
	return inst
}

// SetReward sets the "reward" parameter.
func (inst *RegisterQuestReward) SetReward(reward Reward) *RegisterQuestReward {
	inst.Reward = &reward
	return inst
}

// SetOracleAccount sets the "oracle" account.
func (inst *RegisterQuestReward) SetOracleAccount(oracle ag_solanago.PublicKey) *RegisterQuestReward {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(oracle).WRITE().SIGNER()
	return inst
}

// GetOracleAccount gets the "oracle" account.
func (inst *RegisterQuestReward) GetOracleAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetQuestAccount sets the "quest" account.
func (inst *RegisterQuestReward) SetQuestAccount(quest ag_solanago.PublicKey) *RegisterQuestReward {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(quest).WRITE()
	return inst
}

// GetQuestAccount gets the "quest" account.
func (inst *RegisterQuestReward) GetQuestAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetRewardMintAccount sets the "rewardMint" account.
func (inst *RegisterQuestReward) SetRewardMintAccount(rewardMint ag_solanago.PublicKey) *RegisterQuestReward {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(rewardMint).WRITE().SIGNER()
	return inst
}

// GetRewardMintAccount gets the "rewardMint" account.
func (inst *RegisterQuestReward) GetRewardMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *RegisterQuestReward) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *RegisterQuestReward {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *RegisterQuestReward) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *RegisterQuestReward) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *RegisterQuestReward {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *RegisterQuestReward) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetRentAccount sets the "rent" account.
func (inst *RegisterQuestReward) SetRentAccount(rent ag_solanago.PublicKey) *RegisterQuestReward {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(rent)
	return inst
}

// GetRentAccount gets the "rent" account.
func (inst *RegisterQuestReward) GetRentAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

func (inst RegisterQuestReward) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_RegisterQuestReward,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst RegisterQuestReward) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *RegisterQuestReward) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.QuestIndex == nil {
			return errors.New("QuestIndex parameter is not set")
		}
		if inst.QuestBump == nil {
			return errors.New("QuestBump parameter is not set")
		}
		if inst.Reward == nil {
			return errors.New("Reward parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Oracle is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Quest is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.RewardMint is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.Rent is not set")
		}
	}
	return nil
}

func (inst *RegisterQuestReward) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("RegisterQuestReward")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=3]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("QuestIndex", *inst.QuestIndex))
						paramsBranch.Child(ag_format.Param(" QuestBump", *inst.QuestBump))
						paramsBranch.Child(ag_format.Param("    Reward", *inst.Reward))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=6]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("       oracle", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("        quest", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("   rewardMint", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("systemProgram", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta(" tokenProgram", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("         rent", inst.AccountMetaSlice.Get(5)))
					})
				})
		})
}

func (obj RegisterQuestReward) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `QuestIndex` param:
	err = encoder.Encode(obj.QuestIndex)
	if err != nil {
		return err
	}
	// Serialize `QuestBump` param:
	err = encoder.Encode(obj.QuestBump)
	if err != nil {
		return err
	}
	// Serialize `Reward` param:
	err = encoder.Encode(obj.Reward)
	if err != nil {
		return err
	}
	return nil
}
func (obj *RegisterQuestReward) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `QuestIndex`:
	err = decoder.Decode(&obj.QuestIndex)
	if err != nil {
		return err
	}
	// Deserialize `QuestBump`:
	err = decoder.Decode(&obj.QuestBump)
	if err != nil {
		return err
	}
	// Deserialize `Reward`:
	err = decoder.Decode(&obj.Reward)
	if err != nil {
		return err
	}
	return nil
}

// NewRegisterQuestRewardInstruction declares a new RegisterQuestReward instruction with the provided parameters and accounts.
func NewRegisterQuestRewardInstruction(
	// Parameters:
	questIndex uint64,
	questBump uint8,
	reward Reward,
	// Accounts:
	oracle ag_solanago.PublicKey,
	quest ag_solanago.PublicKey,
	rewardMint ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	rent ag_solanago.PublicKey) *RegisterQuestReward {
	return NewRegisterQuestRewardInstructionBuilder().
		SetQuestIndex(questIndex).
		SetQuestBump(questBump).
		SetReward(reward).
		SetOracleAccount(oracle).
		SetQuestAccount(quest).
		SetRewardMintAccount(rewardMint).
		SetSystemProgramAccount(systemProgram).
		SetTokenProgramAccount(tokenProgram).
		SetRentAccount(rent)
}