package integrations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	ag_binary "github.com/gagliardetto/binary"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"strconv"
	"syscall/js"

	"creaturez.nft/wasm/v2/someplace"
	"github.com/gagliardetto/solana-go"
	atok "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

func Buy(this js.Value, args []js.Value) interface{} {
	holder := solana.MustPublicKeyFromBase58(args[0].String())
	batchAccount := solana.MustPublicKeyFromBase58(args[1].String())
	fmt.Println(args[2].String(), args[2].String(), args[2].String())
	configIndex, configIndexError := strconv.Atoi(args[2].String())

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			if configIndexError != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}
			mintTx, err := mint(holder, batchAccount, configIndex)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}

			dst := js.Global().Get("Uint8Array").New(len(mintTx))
			js.CopyBytesToJS(dst, mintTx)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func mint(holder solana.PublicKey, candyMachineAddress solana.PublicKey, configIndex int) ([]byte, error) {
	mint := solana.NewWallet()
	fmt.Println("wasm-mintKp", mint.PublicKey().String())

	client := rpc.New("https://sparkling-dark-shadow.solana-devnet.quiknode.pro/0e9964e4d70fe7f856e7d03bc7e41dc6a2b84452/")
	userTokenAccountAddress, err := getTokenWallet(holder, mint.PublicKey())
	if err != nil {
		return []byte{}, errors.New("bad")
	}

	candyMachineRaw, err := client.GetAccountInfo(context.TODO(), candyMachineAddress)
	if err != nil {
		return []byte{}, errors.New("bad")
	}

	min, err := client.GetMinimumBalanceForRentExemption(context.TODO(), token.MINT_SIZE, rpc.CommitmentFinalized)
	if err != nil {
		return []byte{}, errors.New("bad")
	}

	dec := ag_binary.NewBorshDecoder(candyMachineRaw.Value.Data.GetBinary())
	var cm someplace.Batch
	err = dec.Decode(&cm)
	if err != nil {
		return []byte{}, errors.New("bad")
	}

	var instructions []solana.Instruction
	instructions = append(instructions,
		system.NewCreateAccountInstructionBuilder().
			SetOwner(token.ProgramID).
			SetNewAccount(mint.PublicKey()).
			SetSpace(token.MINT_SIZE).
			SetFundingAccount(holder).
			SetLamports(min).
			Build(),

		token.NewInitializeMint2InstructionBuilder().
			SetMintAccount(mint.PublicKey()).
			SetDecimals(0).
			SetMintAuthority(holder).
			SetFreezeAuthority(holder).
			Build(),

		atok.NewCreateInstructionBuilder().
			SetPayer(holder).
			SetWallet(holder).
			SetMint(mint.PublicKey()).
			Build(),

		token.NewMintToInstructionBuilder().
			SetMintAccount(mint.PublicKey()).
			SetDestinationAccount(userTokenAccountAddress).
			SetAuthorityAccount(holder).
			SetAmount(1).
			Build(),
	)

	metadataAddress, err := getMetadata(mint.PublicKey())
	if err != nil {
		return []byte{}, errors.New("bad")
	}
	masterEdition, err := getMasterEdition(mint.PublicKey())
	if err != nil {
		return []byte{}, errors.New("bad")
	}
	candyMachineCreator, creatorBump, err := getCandyMachineCreator(candyMachineAddress)
	if err != nil {
		return []byte{}, errors.New("bad")
	}

	treasuryAuthority, _ := GetTreasuryAuthority(cm.Oracle)
	fmt.Println(cm.Oracle, cm.Name, cm.Data)
	treasuryAuthorityData := GetTreasuryAuthorityData(treasuryAuthority)
	if treasuryAuthorityData == nil {
		fmt.Println("???????")
		return []byte{}, errors.New("bad")

	}
	userTreasuryTokenAccountAddress, err := getTokenWallet(holder, treasuryAuthorityData.TreasuryMint)
	listing, _ := GetListing(cm.Oracle, candyMachineAddress, uint64(configIndex))
	treasuryTokenAccount, _ := GetTreasuryTokenAccount(cm.Oracle)

	mintIx := someplace.NewMintNftInstructionBuilder().
		SetConfigIndex(uint64(configIndex)).
		SetCreatorBump(creatorBump).
		SetCandyMachineAccount(candyMachineAddress).
		SetCandyMachineCreatorAccount(candyMachineCreator).
		SetPayerAccount(holder).
		SetOracleAccount(cm.Oracle).
		SetMintAccount(mint.PublicKey()).
		SetMetadataAccount(metadataAddress).
		SetMasterEditionAccount(masterEdition).
		SetMintAuthorityAccount(holder).
		SetUpdateAuthorityAccount(holder).
		SetTokenMetadataProgramAccount(token_metadata.ProgramID).
		SetTokenProgramAccount(token.ProgramID).
		SetSystemProgramAccount(system.ProgramID).
		SetRentAccount(solana.SysVarRentPubkey).
		SetClockAccount(solana.SysVarClockPubkey).
		SetInstructionSysvarAccountAccount(solana.SysVarInstructionsPubkey).
		SetListingAccount(listing).
		SetInitializerTokenAccountAccount(userTreasuryTokenAccountAddress).
		SetTreasuryTokenAccountAccount(treasuryTokenAccount)
	err = mintIx.Validate()
	if err != nil {
		return []byte{}, errors.New("bad")
	}
	instructions = append(instructions,
		mintIx.
			Build(),
	)

	txBuilder := solana.NewTransactionBuilder()
	for _, ix := range instructions {
		txBuilder = txBuilder.AddInstruction(ix)
	}
	txB, _ := txBuilder.Build()
	txJson, _ := json.MarshalIndent(BuyResponse{
		Tx:      *txB,
		MintKey: mint.PrivateKey.String(),
	}, "", "  ")

	fmt.Println(string(txJson))
	fmt.Println("configline", configIndex)

	return txJson, nil

}
