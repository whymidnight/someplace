// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package someplace

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// FulfillMarketListing is the `fulfillMarketListing` instruction.
type FulfillMarketListing struct {
	MarketAuthorityBump *uint8

	// [0] = [WRITE] marketAuthority
	//
	// [1] = [WRITE] marketListing
	//
	// [2] = [WRITE] marketListingTokenAccount
	//
	// [3] = [WRITE, SIGNER] buyer
	//
	// [4] = [WRITE] nftMint
	//
	// [5] = [WRITE] buyerNftTokenAccount
	//
	// [6] = [WRITE] buyerMarketTokenAccount
	//
	// [7] = [WRITE] sellerMarketTokenAccount
	//
	// [8] = [WRITE] oracle
	//
	// [9] = [] tokenProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewFulfillMarketListingInstructionBuilder creates a new `FulfillMarketListing` instruction builder.
func NewFulfillMarketListingInstructionBuilder() *FulfillMarketListing {
	nd := &FulfillMarketListing{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 10),
	}
	return nd
}

// SetMarketAuthorityBump sets the "marketAuthorityBump" parameter.
func (inst *FulfillMarketListing) SetMarketAuthorityBump(marketAuthorityBump uint8) *FulfillMarketListing {
	inst.MarketAuthorityBump = &marketAuthorityBump
	return inst
}

// SetMarketAuthorityAccount sets the "marketAuthority" account.
func (inst *FulfillMarketListing) SetMarketAuthorityAccount(marketAuthority ag_solanago.PublicKey) *FulfillMarketListing {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(marketAuthority).WRITE()
	return inst
}

// GetMarketAuthorityAccount gets the "marketAuthority" account.
func (inst *FulfillMarketListing) GetMarketAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetMarketListingAccount sets the "marketListing" account.
func (inst *FulfillMarketListing) SetMarketListingAccount(marketListing ag_solanago.PublicKey) *FulfillMarketListing {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(marketListing).WRITE()
	return inst
}

// GetMarketListingAccount gets the "marketListing" account.
func (inst *FulfillMarketListing) GetMarketListingAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetMarketListingTokenAccountAccount sets the "marketListingTokenAccount" account.
func (inst *FulfillMarketListing) SetMarketListingTokenAccountAccount(marketListingTokenAccount ag_solanago.PublicKey) *FulfillMarketListing {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(marketListingTokenAccount).WRITE()
	return inst
}

// GetMarketListingTokenAccountAccount gets the "marketListingTokenAccount" account.
func (inst *FulfillMarketListing) GetMarketListingTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetBuyerAccount sets the "buyer" account.
func (inst *FulfillMarketListing) SetBuyerAccount(buyer ag_solanago.PublicKey) *FulfillMarketListing {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(buyer).WRITE().SIGNER()
	return inst
}

// GetBuyerAccount gets the "buyer" account.
func (inst *FulfillMarketListing) GetBuyerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetNftMintAccount sets the "nftMint" account.
func (inst *FulfillMarketListing) SetNftMintAccount(nftMint ag_solanago.PublicKey) *FulfillMarketListing {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(nftMint).WRITE()
	return inst
}

// GetNftMintAccount gets the "nftMint" account.
func (inst *FulfillMarketListing) GetNftMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetBuyerNftTokenAccountAccount sets the "buyerNftTokenAccount" account.
func (inst *FulfillMarketListing) SetBuyerNftTokenAccountAccount(buyerNftTokenAccount ag_solanago.PublicKey) *FulfillMarketListing {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(buyerNftTokenAccount).WRITE()
	return inst
}

// GetBuyerNftTokenAccountAccount gets the "buyerNftTokenAccount" account.
func (inst *FulfillMarketListing) GetBuyerNftTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetBuyerMarketTokenAccountAccount sets the "buyerMarketTokenAccount" account.
func (inst *FulfillMarketListing) SetBuyerMarketTokenAccountAccount(buyerMarketTokenAccount ag_solanago.PublicKey) *FulfillMarketListing {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(buyerMarketTokenAccount).WRITE()
	return inst
}

// GetBuyerMarketTokenAccountAccount gets the "buyerMarketTokenAccount" account.
func (inst *FulfillMarketListing) GetBuyerMarketTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetSellerMarketTokenAccountAccount sets the "sellerMarketTokenAccount" account.
func (inst *FulfillMarketListing) SetSellerMarketTokenAccountAccount(sellerMarketTokenAccount ag_solanago.PublicKey) *FulfillMarketListing {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(sellerMarketTokenAccount).WRITE()
	return inst
}

// GetSellerMarketTokenAccountAccount gets the "sellerMarketTokenAccount" account.
func (inst *FulfillMarketListing) GetSellerMarketTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetOracleAccount sets the "oracle" account.
func (inst *FulfillMarketListing) SetOracleAccount(oracle ag_solanago.PublicKey) *FulfillMarketListing {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(oracle).WRITE()
	return inst
}

// GetOracleAccount gets the "oracle" account.
func (inst *FulfillMarketListing) GetOracleAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *FulfillMarketListing) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *FulfillMarketListing {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *FulfillMarketListing) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

func (inst FulfillMarketListing) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_FulfillMarketListing,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst FulfillMarketListing) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *FulfillMarketListing) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.MarketAuthorityBump == nil {
			return errors.New("MarketAuthorityBump parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.MarketAuthority is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.MarketListing is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.MarketListingTokenAccount is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.Buyer is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.NftMint is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.BuyerNftTokenAccount is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.BuyerMarketTokenAccount is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.SellerMarketTokenAccount is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.Oracle is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
	}
	return nil
}

func (inst *FulfillMarketListing) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("FulfillMarketListing")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("MarketAuthorityBump", *inst.MarketAuthorityBump))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=10]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("   marketAuthority", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("     marketListing", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("marketListingToken", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("             buyer", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("           nftMint", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("     buyerNftToken", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("  buyerMarketToken", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta(" sellerMarketToken", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("            oracle", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta("      tokenProgram", inst.AccountMetaSlice.Get(9)))
					})
				})
		})
}

func (obj FulfillMarketListing) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `MarketAuthorityBump` param:
	err = encoder.Encode(obj.MarketAuthorityBump)
	if err != nil {
		return err
	}
	return nil
}
func (obj *FulfillMarketListing) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `MarketAuthorityBump`:
	err = decoder.Decode(&obj.MarketAuthorityBump)
	if err != nil {
		return err
	}
	return nil
}

// NewFulfillMarketListingInstruction declares a new FulfillMarketListing instruction with the provided parameters and accounts.
func NewFulfillMarketListingInstruction(
	// Parameters:
	marketAuthorityBump uint8,
	// Accounts:
	marketAuthority ag_solanago.PublicKey,
	marketListing ag_solanago.PublicKey,
	marketListingTokenAccount ag_solanago.PublicKey,
	buyer ag_solanago.PublicKey,
	nftMint ag_solanago.PublicKey,
	buyerNftTokenAccount ag_solanago.PublicKey,
	buyerMarketTokenAccount ag_solanago.PublicKey,
	sellerMarketTokenAccount ag_solanago.PublicKey,
	oracle ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey) *FulfillMarketListing {
	return NewFulfillMarketListingInstructionBuilder().
		SetMarketAuthorityBump(marketAuthorityBump).
		SetMarketAuthorityAccount(marketAuthority).
		SetMarketListingAccount(marketListing).
		SetMarketListingTokenAccountAccount(marketListingTokenAccount).
		SetBuyerAccount(buyer).
		SetNftMintAccount(nftMint).
		SetBuyerNftTokenAccountAccount(buyerNftTokenAccount).
		SetBuyerMarketTokenAccountAccount(buyerMarketTokenAccount).
		SetSellerMarketTokenAccountAccount(sellerMarketTokenAccount).
		SetOracleAccount(oracle).
		SetTokenProgramAccount(tokenProgram)
}
