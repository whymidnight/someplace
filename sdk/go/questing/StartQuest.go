// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package questing

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// StartQuest is the `startQuest` instruction.
type StartQuest struct {
	QuestIndex *uint64

	// [0] = [WRITE] quest
	//
	// [1] = [WRITE, SIGNER] initializer
	//
	// [2] = [WRITE] ballzTokenAccount
	//
	// [3] = [WRITE] ballzMint
	//
	// [4] = [WRITE] depositTokenAccount
	//
	// [5] = [WRITE] pixelballzMint
	//
	// [6] = [WRITE] pixelballzTokenAccount
	//
	// [7] = [WRITE] questAccount
	//
	// [8] = [] systemProgram
	//
	// [9] = [] tokenProgram
	//
	// [10] = [] rent
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewStartQuestInstructionBuilder creates a new `StartQuest` instruction builder.
func NewStartQuestInstructionBuilder() *StartQuest {
	nd := &StartQuest{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 11),
	}
	return nd
}

// SetQuestIndex sets the "questIndex" parameter.
func (inst *StartQuest) SetQuestIndex(questIndex uint64) *StartQuest {
	inst.QuestIndex = &questIndex
	return inst
}

// SetQuestAccount sets the "quest" account.
func (inst *StartQuest) SetQuestAccount(quest ag_solanago.PublicKey) *StartQuest {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(quest).WRITE()
	return inst
}

// GetQuestAccount gets the "quest" account.
func (inst *StartQuest) GetQuestAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetInitializerAccount sets the "initializer" account.
func (inst *StartQuest) SetInitializerAccount(initializer ag_solanago.PublicKey) *StartQuest {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(initializer).WRITE().SIGNER()
	return inst
}

// GetInitializerAccount gets the "initializer" account.
func (inst *StartQuest) GetInitializerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetBallzTokenAccountAccount sets the "ballzTokenAccount" account.
func (inst *StartQuest) SetBallzTokenAccountAccount(ballzTokenAccount ag_solanago.PublicKey) *StartQuest {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(ballzTokenAccount).WRITE()
	return inst
}

// GetBallzTokenAccountAccount gets the "ballzTokenAccount" account.
func (inst *StartQuest) GetBallzTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetBallzMintAccount sets the "ballzMint" account.
func (inst *StartQuest) SetBallzMintAccount(ballzMint ag_solanago.PublicKey) *StartQuest {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(ballzMint).WRITE()
	return inst
}

// GetBallzMintAccount gets the "ballzMint" account.
func (inst *StartQuest) GetBallzMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetDepositTokenAccountAccount sets the "depositTokenAccount" account.
func (inst *StartQuest) SetDepositTokenAccountAccount(depositTokenAccount ag_solanago.PublicKey) *StartQuest {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(depositTokenAccount).WRITE()
	return inst
}

// GetDepositTokenAccountAccount gets the "depositTokenAccount" account.
func (inst *StartQuest) GetDepositTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetPixelballzMintAccount sets the "pixelballzMint" account.
func (inst *StartQuest) SetPixelballzMintAccount(pixelballzMint ag_solanago.PublicKey) *StartQuest {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(pixelballzMint).WRITE()
	return inst
}

// GetPixelballzMintAccount gets the "pixelballzMint" account.
func (inst *StartQuest) GetPixelballzMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetPixelballzTokenAccountAccount sets the "pixelballzTokenAccount" account.
func (inst *StartQuest) SetPixelballzTokenAccountAccount(pixelballzTokenAccount ag_solanago.PublicKey) *StartQuest {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(pixelballzTokenAccount).WRITE()
	return inst
}

// GetPixelballzTokenAccountAccount gets the "pixelballzTokenAccount" account.
func (inst *StartQuest) GetPixelballzTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetQuestAccountAccount sets the "questAccount" account.
func (inst *StartQuest) SetQuestAccountAccount(questAccount ag_solanago.PublicKey) *StartQuest {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(questAccount).WRITE()
	return inst
}

// GetQuestAccountAccount gets the "questAccount" account.
func (inst *StartQuest) GetQuestAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *StartQuest) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *StartQuest {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *StartQuest) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *StartQuest) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *StartQuest {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *StartQuest) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

// SetRentAccount sets the "rent" account.
func (inst *StartQuest) SetRentAccount(rent ag_solanago.PublicKey) *StartQuest {
	inst.AccountMetaSlice[10] = ag_solanago.Meta(rent)
	return inst
}

// GetRentAccount gets the "rent" account.
func (inst *StartQuest) GetRentAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(10)
}

func (inst StartQuest) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_StartQuest,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst StartQuest) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *StartQuest) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.QuestIndex == nil {
			return errors.New("QuestIndex parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Quest is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Initializer is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.BallzTokenAccount is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.BallzMint is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.DepositTokenAccount is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.PixelballzMint is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.PixelballzTokenAccount is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.QuestAccount is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[10] == nil {
			return errors.New("accounts.Rent is not set")
		}
	}
	return nil
}

func (inst *StartQuest) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("StartQuest")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("QuestIndex", *inst.QuestIndex))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=11]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("          quest", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("    initializer", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("     ballzToken", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("      ballzMint", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("   depositToken", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta(" pixelballzMint", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("pixelballzToken", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("          quest", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("  systemProgram", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta("   tokenProgram", inst.AccountMetaSlice.Get(9)))
						accountsBranch.Child(ag_format.Meta("           rent", inst.AccountMetaSlice.Get(10)))
					})
				})
		})
}

func (obj StartQuest) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `QuestIndex` param:
	err = encoder.Encode(obj.QuestIndex)
	if err != nil {
		return err
	}
	return nil
}
func (obj *StartQuest) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `QuestIndex`:
	err = decoder.Decode(&obj.QuestIndex)
	if err != nil {
		return err
	}
	return nil
}

// NewStartQuestInstruction declares a new StartQuest instruction with the provided parameters and accounts.
func NewStartQuestInstruction(
	// Parameters:
	questIndex uint64,
	// Accounts:
	quest ag_solanago.PublicKey,
	initializer ag_solanago.PublicKey,
	ballzTokenAccount ag_solanago.PublicKey,
	ballzMint ag_solanago.PublicKey,
	depositTokenAccount ag_solanago.PublicKey,
	pixelballzMint ag_solanago.PublicKey,
	pixelballzTokenAccount ag_solanago.PublicKey,
	questAccount ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	rent ag_solanago.PublicKey) *StartQuest {
	return NewStartQuestInstructionBuilder().
		SetQuestIndex(questIndex).
		SetQuestAccount(quest).
		SetInitializerAccount(initializer).
		SetBallzTokenAccountAccount(ballzTokenAccount).
		SetBallzMintAccount(ballzMint).
		SetDepositTokenAccountAccount(depositTokenAccount).
		SetPixelballzMintAccount(pixelballzMint).
		SetPixelballzTokenAccountAccount(pixelballzTokenAccount).
		SetQuestAccountAccount(questAccount).
		SetSystemProgramAccount(systemProgram).
		SetTokenProgramAccount(tokenProgram).
		SetRentAccount(rent)
}