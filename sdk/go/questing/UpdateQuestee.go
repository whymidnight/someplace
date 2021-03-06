// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package questing

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// UpdateQuestee is the `updateQuestee` instruction.
type UpdateQuestee struct {

	// [0] = [WRITE, SIGNER] newOwner
	//
	// [1] = [WRITE] questee
	//
	// [2] = [] pixelballzMint
	//
	// [3] = [] pixelballzTokenAccount
	//
	// [4] = [] owner
	//
	// [5] = [WRITE] questor
	//
	// [6] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewUpdateQuesteeInstructionBuilder creates a new `UpdateQuestee` instruction builder.
func NewUpdateQuesteeInstructionBuilder() *UpdateQuestee {
	nd := &UpdateQuestee{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 7),
	}
	return nd
}

// SetNewOwnerAccount sets the "newOwner" account.
func (inst *UpdateQuestee) SetNewOwnerAccount(newOwner ag_solanago.PublicKey) *UpdateQuestee {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(newOwner).WRITE().SIGNER()
	return inst
}

// GetNewOwnerAccount gets the "newOwner" account.
func (inst *UpdateQuestee) GetNewOwnerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetQuesteeAccount sets the "questee" account.
func (inst *UpdateQuestee) SetQuesteeAccount(questee ag_solanago.PublicKey) *UpdateQuestee {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(questee).WRITE()
	return inst
}

// GetQuesteeAccount gets the "questee" account.
func (inst *UpdateQuestee) GetQuesteeAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetPixelballzMintAccount sets the "pixelballzMint" account.
func (inst *UpdateQuestee) SetPixelballzMintAccount(pixelballzMint ag_solanago.PublicKey) *UpdateQuestee {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(pixelballzMint)
	return inst
}

// GetPixelballzMintAccount gets the "pixelballzMint" account.
func (inst *UpdateQuestee) GetPixelballzMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetPixelballzTokenAccountAccount sets the "pixelballzTokenAccount" account.
func (inst *UpdateQuestee) SetPixelballzTokenAccountAccount(pixelballzTokenAccount ag_solanago.PublicKey) *UpdateQuestee {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(pixelballzTokenAccount)
	return inst
}

// GetPixelballzTokenAccountAccount gets the "pixelballzTokenAccount" account.
func (inst *UpdateQuestee) GetPixelballzTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetOwnerAccount sets the "owner" account.
func (inst *UpdateQuestee) SetOwnerAccount(owner ag_solanago.PublicKey) *UpdateQuestee {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(owner)
	return inst
}

// GetOwnerAccount gets the "owner" account.
func (inst *UpdateQuestee) GetOwnerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetQuestorAccount sets the "questor" account.
func (inst *UpdateQuestee) SetQuestorAccount(questor ag_solanago.PublicKey) *UpdateQuestee {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(questor).WRITE()
	return inst
}

// GetQuestorAccount gets the "questor" account.
func (inst *UpdateQuestee) GetQuestorAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *UpdateQuestee) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *UpdateQuestee {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *UpdateQuestee) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

func (inst UpdateQuestee) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_UpdateQuestee,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst UpdateQuestee) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *UpdateQuestee) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.NewOwner is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Questee is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.PixelballzMint is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.PixelballzTokenAccount is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.Questor is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *UpdateQuestee) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("UpdateQuestee")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=7]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("       newOwner", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("        questee", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta(" pixelballzMint", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("pixelballzToken", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("          owner", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("        questor", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("  systemProgram", inst.AccountMetaSlice.Get(6)))
					})
				})
		})
}

func (obj UpdateQuestee) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *UpdateQuestee) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewUpdateQuesteeInstruction declares a new UpdateQuestee instruction with the provided parameters and accounts.
func NewUpdateQuesteeInstruction(
	// Accounts:
	newOwner ag_solanago.PublicKey,
	questee ag_solanago.PublicKey,
	pixelballzMint ag_solanago.PublicKey,
	pixelballzTokenAccount ag_solanago.PublicKey,
	owner ag_solanago.PublicKey,
	questor ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *UpdateQuestee {
	return NewUpdateQuesteeInstructionBuilder().
		SetNewOwnerAccount(newOwner).
		SetQuesteeAccount(questee).
		SetPixelballzMintAccount(pixelballzMint).
		SetPixelballzTokenAccountAccount(pixelballzTokenAccount).
		SetOwnerAccount(owner).
		SetQuestorAccount(questor).
		SetSystemProgramAccount(systemProgram)
}
