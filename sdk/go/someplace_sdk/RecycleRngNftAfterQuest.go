// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package someplace

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// RecycleRngNftAfterQuest is the `recycleRngNftAfterQuest` instruction.
type RecycleRngNftAfterQuest struct {
	ViaBump *uint8

	// [0] = [] rewardTokenAccount
	//
	// [1] = [WRITE] rewardTicket
	//
	// [2] = [] batches
	//
	// [3] = [] via
	//
	// [4] = [] viaMap
	//
	// [5] = [] quest
	//
	// [6] = [WRITE, SIGNER] initializer
	//
	// [7] = [] systemProgram
	//
	// [8] = [] slotHashes
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewRecycleRngNftAfterQuestInstructionBuilder creates a new `RecycleRngNftAfterQuest` instruction builder.
func NewRecycleRngNftAfterQuestInstructionBuilder() *RecycleRngNftAfterQuest {
	nd := &RecycleRngNftAfterQuest{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 9),
	}
	return nd
}

// SetViaBump sets the "viaBump" parameter.
func (inst *RecycleRngNftAfterQuest) SetViaBump(viaBump uint8) *RecycleRngNftAfterQuest {
	inst.ViaBump = &viaBump
	return inst
}

// SetRewardTokenAccountAccount sets the "rewardTokenAccount" account.
func (inst *RecycleRngNftAfterQuest) SetRewardTokenAccountAccount(rewardTokenAccount ag_solanago.PublicKey) *RecycleRngNftAfterQuest {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(rewardTokenAccount)
	return inst
}

// GetRewardTokenAccountAccount gets the "rewardTokenAccount" account.
func (inst *RecycleRngNftAfterQuest) GetRewardTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetRewardTicketAccount sets the "rewardTicket" account.
func (inst *RecycleRngNftAfterQuest) SetRewardTicketAccount(rewardTicket ag_solanago.PublicKey) *RecycleRngNftAfterQuest {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(rewardTicket).WRITE()
	return inst
}

// GetRewardTicketAccount gets the "rewardTicket" account.
func (inst *RecycleRngNftAfterQuest) GetRewardTicketAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetBatchesAccount sets the "batches" account.
func (inst *RecycleRngNftAfterQuest) SetBatchesAccount(batches ag_solanago.PublicKey) *RecycleRngNftAfterQuest {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(batches)
	return inst
}

// GetBatchesAccount gets the "batches" account.
func (inst *RecycleRngNftAfterQuest) GetBatchesAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetViaAccount sets the "via" account.
func (inst *RecycleRngNftAfterQuest) SetViaAccount(via ag_solanago.PublicKey) *RecycleRngNftAfterQuest {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(via)
	return inst
}

// GetViaAccount gets the "via" account.
func (inst *RecycleRngNftAfterQuest) GetViaAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetViaMapAccount sets the "viaMap" account.
func (inst *RecycleRngNftAfterQuest) SetViaMapAccount(viaMap ag_solanago.PublicKey) *RecycleRngNftAfterQuest {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(viaMap)
	return inst
}

// GetViaMapAccount gets the "viaMap" account.
func (inst *RecycleRngNftAfterQuest) GetViaMapAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetQuestAccount sets the "quest" account.
func (inst *RecycleRngNftAfterQuest) SetQuestAccount(quest ag_solanago.PublicKey) *RecycleRngNftAfterQuest {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(quest)
	return inst
}

// GetQuestAccount gets the "quest" account.
func (inst *RecycleRngNftAfterQuest) GetQuestAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetInitializerAccount sets the "initializer" account.
func (inst *RecycleRngNftAfterQuest) SetInitializerAccount(initializer ag_solanago.PublicKey) *RecycleRngNftAfterQuest {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(initializer).WRITE().SIGNER()
	return inst
}

// GetInitializerAccount gets the "initializer" account.
func (inst *RecycleRngNftAfterQuest) GetInitializerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *RecycleRngNftAfterQuest) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *RecycleRngNftAfterQuest {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *RecycleRngNftAfterQuest) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetSlotHashesAccount sets the "slotHashes" account.
func (inst *RecycleRngNftAfterQuest) SetSlotHashesAccount(slotHashes ag_solanago.PublicKey) *RecycleRngNftAfterQuest {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(slotHashes)
	return inst
}

// GetSlotHashesAccount gets the "slotHashes" account.
func (inst *RecycleRngNftAfterQuest) GetSlotHashesAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

func (inst RecycleRngNftAfterQuest) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_RecycleRngNftAfterQuest,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst RecycleRngNftAfterQuest) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *RecycleRngNftAfterQuest) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.ViaBump == nil {
			return errors.New("ViaBump parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.RewardTokenAccount is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.RewardTicket is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.Batches is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.Via is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.ViaMap is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.Quest is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.Initializer is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.SlotHashes is not set")
		}
	}
	return nil
}

func (inst *RecycleRngNftAfterQuest) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("RecycleRngNftAfterQuest")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("ViaBump", *inst.ViaBump))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=9]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("  rewardToken", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta(" rewardTicket", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("      batches", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("          via", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("       viaMap", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("        quest", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("  initializer", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("systemProgram", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("   slotHashes", inst.AccountMetaSlice.Get(8)))
					})
				})
		})
}

func (obj RecycleRngNftAfterQuest) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `ViaBump` param:
	err = encoder.Encode(obj.ViaBump)
	if err != nil {
		return err
	}
	return nil
}
func (obj *RecycleRngNftAfterQuest) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `ViaBump`:
	err = decoder.Decode(&obj.ViaBump)
	if err != nil {
		return err
	}
	return nil
}

// NewRecycleRngNftAfterQuestInstruction declares a new RecycleRngNftAfterQuest instruction with the provided parameters and accounts.
func NewRecycleRngNftAfterQuestInstruction(
	// Parameters:
	viaBump uint8,
	// Accounts:
	rewardTokenAccount ag_solanago.PublicKey,
	rewardTicket ag_solanago.PublicKey,
	batches ag_solanago.PublicKey,
	via ag_solanago.PublicKey,
	viaMap ag_solanago.PublicKey,
	quest ag_solanago.PublicKey,
	initializer ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey,
	slotHashes ag_solanago.PublicKey) *RecycleRngNftAfterQuest {
	return NewRecycleRngNftAfterQuestInstructionBuilder().
		SetViaBump(viaBump).
		SetRewardTokenAccountAccount(rewardTokenAccount).
		SetRewardTicketAccount(rewardTicket).
		SetBatchesAccount(batches).
		SetViaAccount(via).
		SetViaMapAccount(viaMap).
		SetQuestAccount(quest).
		SetInitializerAccount(initializer).
		SetSystemProgramAccount(systemProgram).
		SetSlotHashesAccount(slotHashes)
}
