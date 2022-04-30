// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package someplace

import (
	"bytes"
	"fmt"
	ag_spew "github.com/davecgh/go-spew/spew"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_text "github.com/gagliardetto/solana-go/text"
	ag_treeout "github.com/gagliardetto/treeout"
)

var ProgramID ag_solanago.PublicKey = ag_solanago.MustPublicKeyFromBase58("8otw5mCMUtwx91e7q7MAyhWoQVnc3Ng72qwDH58z72VW")

func SetProgramID(pubkey ag_solanago.PublicKey) {
	ProgramID = pubkey
	ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
}

const ProgramName = "Someplace"

func init() {
	if !ProgramID.IsZero() {
		ag_solanago.RegisterInstructionDecoder(ProgramID, registryDecodeInstruction)
	}
}

var (
	Instruction_CreateMarketListing = ag_binary.TypeID([8]byte{84, 255, 145, 58, 74, 134, 93, 15})

	Instruction_FulfillMarketListing = ag_binary.TypeID([8]byte{171, 101, 117, 183, 127, 117, 107, 96})

	Instruction_UnlistMarketListing = ag_binary.TypeID([8]byte{185, 12, 146, 92, 148, 242, 148, 229})

	Instruction_CreateListing = ag_binary.TypeID([8]byte{18, 168, 45, 24, 191, 31, 117, 54})

	Instruction_ModifyListing = ag_binary.TypeID([8]byte{36, 132, 230, 119, 139, 147, 164, 183})

	Instruction_EnableBatchUploading = ag_binary.TypeID([8]byte{212, 38, 162, 41, 25, 159, 102, 80})

	Instruction_InitMarket = ag_binary.TypeID([8]byte{33, 253, 15, 116, 89, 25, 127, 236})

	Instruction_InitTreasury = ag_binary.TypeID([8]byte{105, 152, 173, 51, 158, 151, 49, 14})

	Instruction_AddWhitelistedCm = ag_binary.TypeID([8]byte{164, 150, 77, 187, 189, 197, 159, 250})

	Instruction_AmmendStorefrontSplits = ag_binary.TypeID([8]byte{191, 150, 235, 205, 222, 97, 48, 220})

	Instruction_SellFor = ag_binary.TypeID([8]byte{213, 119, 234, 122, 233, 227, 107, 158})

	Instruction_AddConfigLines = ag_binary.TypeID([8]byte{223, 50, 224, 227, 151, 8, 115, 106})

	Instruction_InitializeCandyMachine = ag_binary.TypeID([8]byte{142, 137, 167, 107, 47, 39, 240, 124})

	Instruction_MintNft = ag_binary.TypeID([8]byte{211, 57, 6, 167, 15, 219, 35, 251})

	Instruction_MintNftRarity = ag_binary.TypeID([8]byte{56, 162, 100, 81, 69, 67, 107, 212})
)

// InstructionIDToName returns the name of the instruction given its ID.
func InstructionIDToName(id ag_binary.TypeID) string {
	switch id {
	case Instruction_CreateMarketListing:
		return "CreateMarketListing"
	case Instruction_FulfillMarketListing:
		return "FulfillMarketListing"
	case Instruction_UnlistMarketListing:
		return "UnlistMarketListing"
	case Instruction_CreateListing:
		return "CreateListing"
	case Instruction_ModifyListing:
		return "ModifyListing"
	case Instruction_EnableBatchUploading:
		return "EnableBatchUploading"
	case Instruction_InitMarket:
		return "InitMarket"
	case Instruction_InitTreasury:
		return "InitTreasury"
	case Instruction_AddWhitelistedCm:
		return "AddWhitelistedCm"
	case Instruction_AmmendStorefrontSplits:
		return "AmmendStorefrontSplits"
	case Instruction_SellFor:
		return "SellFor"
	case Instruction_AddConfigLines:
		return "AddConfigLines"
	case Instruction_InitializeCandyMachine:
		return "InitializeCandyMachine"
	case Instruction_MintNft:
		return "MintNft"
	case Instruction_MintNftRarity:
		return "MintNftRarity"
	default:
		return ""
	}
}

type Instruction struct {
	ag_binary.BaseVariant
}

func (inst *Instruction) EncodeToTree(parent ag_treeout.Branches) {
	if enToTree, ok := inst.Impl.(ag_text.EncodableToTree); ok {
		enToTree.EncodeToTree(parent)
	} else {
		parent.Child(ag_spew.Sdump(inst))
	}
}

var InstructionImplDef = ag_binary.NewVariantDefinition(
	ag_binary.AnchorTypeIDEncoding,
	[]ag_binary.VariantType{
		{
			"create_market_listing", (*CreateMarketListing)(nil),
		},
		{
			"fulfill_market_listing", (*FulfillMarketListing)(nil),
		},
		{
			"unlist_market_listing", (*UnlistMarketListing)(nil),
		},
		{
			"create_listing", (*CreateListing)(nil),
		},
		{
			"modify_listing", (*ModifyListing)(nil),
		},
		{
			"enable_batch_uploading", (*EnableBatchUploading)(nil),
		},
		{
			"init_market", (*InitMarket)(nil),
		},
		{
			"init_treasury", (*InitTreasury)(nil),
		},
		{
			"add_whitelisted_cm", (*AddWhitelistedCm)(nil),
		},
		{
			"ammend_storefront_splits", (*AmmendStorefrontSplits)(nil),
		},
		{
			"sell_for", (*SellFor)(nil),
		},
		{
			"add_config_lines", (*AddConfigLines)(nil),
		},
		{
			"initialize_candy_machine", (*InitializeCandyMachine)(nil),
		},
		{
			"mint_nft", (*MintNft)(nil),
		},
		{
			"mint_nft_rarity", (*MintNftRarity)(nil),
		},
	},
)

func (inst *Instruction) ProgramID() ag_solanago.PublicKey {
	return ProgramID
}

func (inst *Instruction) Accounts() (out []*ag_solanago.AccountMeta) {
	return inst.Impl.(ag_solanago.AccountsGettable).GetAccounts()
}

func (inst *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := ag_binary.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *Instruction) TextEncode(encoder *ag_text.Encoder, option *ag_text.Option) error {
	return encoder.Encode(inst.Impl, option)
}

func (inst *Instruction) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	return inst.BaseVariant.UnmarshalBinaryVariant(decoder, InstructionImplDef)
}

func (inst *Instruction) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	err := encoder.WriteBytes(inst.TypeID.Bytes(), false)
	if err != nil {
		return fmt.Errorf("unable to write variant type: %w", err)
	}
	return encoder.Encode(inst.Impl)
}

func registryDecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (interface{}, error) {
	inst, err := DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func DecodeInstruction(accounts []*ag_solanago.AccountMeta, data []byte) (*Instruction, error) {
	inst := new(Instruction)
	if err := ag_binary.NewBorshDecoder(data).Decode(inst); err != nil {
		return nil, fmt.Errorf("unable to decode instruction: %w", err)
	}
	if v, ok := inst.Impl.(ag_solanago.AccountsSettable); ok {
		err := v.SetAccounts(accounts)
		if err != nil {
			return nil, fmt.Errorf("unable to set accounts for instruction: %w", err)
		}
	}
	return inst, nil
}
