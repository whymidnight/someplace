package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/token"

	"creaturez.nft/someplace/v2/someplace"
	sendAndConfirmTransaction "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/text"

	"github.com/gagliardetto/solana-go/rpc/ws"

	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
	atok "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"io/ioutil"
	"net/http"
)

const DEVNET = "https://sparkling-dark-shadow.solana-devnet.quiknode.pro/0e9964e4d70fe7f856e7d03bc7e41dc6a2b84452/"
const TESTNET = "https://api.testnet.solana.com"
const NETWORK = TESTNET
const CDN = "https://triptychlabs.io:4445"

var MINT = solana.MustPublicKeyFromBase58("9K9h3f5dEPyqEvaJ2kjNSbjwBq7j9ri1Bn8soF41J2w1")

func init() {
	someplace.SetProgramID(solana.MustPublicKeyFromBase58("8otw5mCMUtwx91e7q7MAyhWoQVnc3Ng72qwDH58z72VW"))
}

func main() {
	// enable()
	// verifyBatchUpload()
	// catalogBatches()
	// treasure()
	// list()
	// verifyList()
	// treasureCMs()
	// treasureVerify()
	// treasureVerifyCM()
	// mint()
	mintRare()
	// holder_nft_metadata()
	// burn()

	// marketCreate()
	// verifyMarketCreate()
	// marketList()
	// verifyMarketList()
	// marketFulfill()

	// GetMarketMintMeta()
	// GetMarketListingsData()

	public := solana.MustPublicKeyFromBase58("DdGbS9PZyxLwJLisa8a9tuTE7V3xJ9Hs6yRPmzzX75gF")
	fmt.Println(getCandyMachineCreator(public))
	fmt.Println([]byte("candy_machine"))

	key := []byte{68, 100, 71, 98, 83, 57, 80, 90, 121, 120, 76, 119, 74, 76, 105, 115, 97, 56, 97, 57, 116, 117, 84, 69, 55, 86, 51, 120, 74, 57, 72, 115, 54, 121, 82, 80, 109, 122, 122, 88, 55, 53, 103, 70}
	fmt.Println(len(key))
	fmt.Println(solana.PublicKeyFromBytes(key))
	fmt.Println(len(public.Bytes()))
	fmt.Println(len("DdGbS9PZyxLwJLisa8a9tuTE7V3xJ9Hs6yRPmzzX75gF"))

}

func verifyMarketCreate() {
	marketUid := solana.MustPublicKeyFromBase58("4Gm324iNEMapZV9aVyWg8EwJYLiqepYYab47sCWcPnh1")
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}

	marketAuthority, _ := GetMarketAuthority(oracle.PublicKey(), marketUid)
	GetMarketAuthorityData(marketAuthority)
}

func marketCreate() {
	marketUid := solana.MustPublicKeyFromBase58("4Gm324iNEMapZV9aVyWg8EwJYLiqepYYab47sCWcPnh1")
	mint := MINT
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	marketAuthority, _ := GetMarketAuthority(oracle.PublicKey(), marketUid)

	listIx := someplace.NewInitMarketInstructionBuilder().
		SetMarketAuthorityAccount(marketAuthority).
		SetMarketMintAccount(mint).
		SetMarketUid(marketUid).
		SetName("market test").
		SetOracleAccount(oracle.PublicKey()).
		SetSystemProgramAccount(solana.SystemProgramID)

	err = listIx.Validate()
	if err != nil {
		panic(err)
	}

	sendTx(
		"list",
		append(make([]solana.Instruction, 0), listIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)

}

func marketList() {
	marketUid := solana.MustPublicKeyFromBase58("GAm4cGVVMi5NMBVoxg8QhuKpQk2xW4BVbhnQrESE54HA")
	marketMint := MINT
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	marketAuthority, _ := GetMarketAuthority(oracle.PublicKey(), marketUid)
	marketAuthorityData := GetMarketAuthorityData(marketAuthority)
	nftMint := solana.NewWallet().PrivateKey

	userTokenAccountAddress, _ := getTokenWallet(oracle.PublicKey(), nftMint.PublicKey())
	sellerMarketTokenAccountAddress, _ := getTokenWallet(oracle.PublicKey(), marketMint)

	client := rpc.New(NETWORK)
	min, err := client.GetMinimumBalanceForRentExemption(context.TODO(), token.MINT_SIZE, rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	var instructions []solana.Instruction
	instructions = append(instructions,
		system.NewCreateAccountInstructionBuilder().
			SetOwner(token.ProgramID).
			SetNewAccount(nftMint.PublicKey()).
			SetSpace(token.MINT_SIZE).
			SetFundingAccount(oracle.PublicKey()).
			SetLamports(min).
			Build(),

		token.NewInitializeMint2InstructionBuilder().
			SetMintAccount(nftMint.PublicKey()).
			SetDecimals(0).
			SetMintAuthority(oracle.PublicKey()).
			SetFreezeAuthority(oracle.PublicKey()).
			Build(),

		atok.NewCreateInstructionBuilder().
			SetPayer(oracle.PublicKey()).
			SetWallet(oracle.PublicKey()).
			SetMint(nftMint.PublicKey()).
			Build(),

		token.NewMintToInstructionBuilder().
			SetMintAccount(nftMint.PublicKey()).
			SetDestinationAccount(userTokenAccountAddress).
			SetAuthorityAccount(oracle.PublicKey()).
			SetAmount(1).
			Build(),
	)

	marketListing, _ := GetMarketListing(marketAuthority, marketAuthorityData.Listings)
	marketListingTokenAccount, _ := GetMarketListingTokenAccount(marketAuthority, marketAuthorityData.Listings)
	listIx := someplace.NewCreateMarketListingInstructionBuilder().
		SetIndex(marketAuthorityData.Listings).
		SetMarketAuthorityAccount(marketAuthority).
		SetMarketListingAccount(marketListing).
		SetMarketListingTokenAccountAccount(marketListingTokenAccount).
		SetNftMintAccount(nftMint.PublicKey()).
		SetPrice(1).
		SetRentAccount(solana.SysVarRentPubkey).
		SetSellerAccount(oracle.PublicKey()).
		SetSellerMarketTokenAccountAccount(sellerMarketTokenAccountAddress).
		SetSellerNftTokenAccountAccount(userTokenAccountAddress).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTokenProgramAccount(solana.TokenProgramID)

	err = listIx.Validate()
	if err != nil {
		panic(err)
	}

	sendTx(
		"list",
		append(instructions, listIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle, nftMint),
		oracle.PublicKey(),
	)

}

func verifyMarketList() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	marketUid := solana.MustPublicKeyFromBase58("GAm4cGVVMi5NMBVoxg8QhuKpQk2xW4BVbhnQrESE54HA")
	marketAuthority, _ := GetMarketAuthority(oracle.PublicKey(), marketUid)
	marketAuthorityData := GetMarketAuthorityData(marketAuthority)
	i := uint64(0)
	for i < marketAuthorityData.Listings {
		marketListing, _ := GetMarketListing(marketAuthority, i)
		GetMarketListingData(marketListing)
		i++
	}

}

func marketFulfill() {
	marketUid := solana.MustPublicKeyFromBase58("GAm4cGVVMi5NMBVoxg8QhuKpQk2xW4BVbhnQrESE54HA")
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	buyer, err := solana.PrivateKeyFromSolanaKeygenFile("./burner.key")
	if err != nil {
		panic(err)
	}
	marketAuthority, marketAuthorityBump := GetMarketAuthority(oracle.PublicKey(), marketUid)
	// marketAuthorityData := GetMarketAuthorityData(marketAuthority)
	marketListing, _ := GetMarketListing(marketAuthority, 1)
	marketListingData := GetMarketListingData(marketListing)
	marketMint := MINT
	buyerMarketTokenAccountAddress, _ := getTokenWallet(buyer.PublicKey(), marketMint)
	buyerNftTokenAccountAddress, _ := getTokenWallet(buyer.PublicKey(), marketListingData.NftMint)

	marketListingTokenAccount, _ := GetMarketListingTokenAccount(marketAuthority, marketListingData.Index)
	var instructions []solana.Instruction
	instructions = append(instructions,
		atok.NewCreateInstructionBuilder().
			SetPayer(buyer.PublicKey()).
			SetWallet(buyer.PublicKey()).
			SetMint(marketListingData.NftMint).
			Build(),
	)
	listIx := someplace.NewFulfillMarketListingInstructionBuilder().
		SetBuyerAccount(buyer.PublicKey()).
		SetBuyerMarketTokenAccountAccount(buyerMarketTokenAccountAddress).
		SetBuyerNftTokenAccountAccount(buyerNftTokenAccountAddress).
		SetMarketAuthorityAccount(marketAuthority).
		SetMarketAuthorityBump(marketAuthorityBump).
		SetMarketListingAccount(marketListing).
		SetMarketListingTokenAccountAccount(marketListingTokenAccount).
		SetNftMintAccount(marketListingData.NftMint).
		SetOracleAccount(oracle.PublicKey()).
		SetSellerMarketTokenAccountAccount(marketListingData.SellerMarketTokenAccount).
		SetTokenProgramAccount(solana.TokenProgramID)

	err = listIx.Validate()
	if err != nil {
		panic(err)
	}

	sendTx(
		"list",
		append(instructions, listIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle, buyer),
		buyer.PublicKey(),
	)

}

func verifyList() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	batch := solana.MustPublicKeyFromBase58("Dfp73qrotgTtecSTkakgTX5rfAqW3JtXgwqFMCL26zaz")
	listing, _ := GetListing(oracle.PublicKey(), batch, 0)

	GetListingData(listing)
}

func list() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	batch := solana.MustPublicKeyFromBase58("GTTXtoSYPhaH27tREZ6PZDvKsQMhaCrqyp7rPcTAzjtg")
	treasuryAuthority, _ := GetTreasuryAuthority(oracle.PublicKey())

	listing, _ := GetListing(oracle.PublicKey(), batch, 0)
	listIx := someplace.NewCreateListingInstructionBuilder().
		SetBatchAccount(batch).
		SetConfigIndex(0).
		SetLifecycleStart(0).
		SetListingAccount(listing).
		SetOracleAccount(oracle.PublicKey()).
		SetPrice(15).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTreasuryAuthorityAccount(treasuryAuthority)

	err = listIx.Validate()
	if err != nil {
		panic(err)
	}

	sendTx(
		"list",
		append(make([]solana.Instruction, 0), listIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)

}

func holder_nft_metadata() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./dev.key")
	if err != nil {
		panic(err)
	}
	client := rpc.New(NETWORK)
	tokenAccounts, err := client.GetTokenAccountsByOwner(context.TODO(), oracle.PublicKey(), &rpc.GetTokenAccountsConfig{ProgramId: &solana.TokenProgramID}, &rpc.GetTokenAccountsOpts{Encoding: "jsonParsed"})
	if err != nil {
		panic(err)
	}
	type tokenDataParsed struct {
		Parsed struct {
			Info struct {
				Mint        string `json:"mint"`
				TokenAmount struct {
					UiAmount float64 `json:"uiAmount"`
				} `json:"tokenAmount"`
			} `json:"info"`
		} `json:"parsed"`
	}
	type token struct {
		ata      solana.PublicKey
		mint     solana.PublicKey
		metadata solana.PublicKey
	}

	_ = func() []token {
		tokens := make([]token, 2)
		for _, tokenAccount := range tokenAccounts.Value {
			var tokenData tokenDataParsed
			err = json.Unmarshal(tokenAccount.Account.Data.GetRawJSON(), &tokenData)
			if err != nil {
				panic(err)
			}
			if tokenData.Parsed.Info.TokenAmount.UiAmount == 1.0 {
				nftMint := solana.MustPublicKeyFromBase58(tokenData.Parsed.Info.Mint)
				metadata, _ := getMetadata(nftMint)
				fmt.Println("........", nftMint, metadata)
				metadataAccount, err := client.GetAccountInfo(context.TODO(), metadata)
				if err != nil {
					panic(err)
				}
				var metadataData token_metadata.Metadata
				metadataDecoder := ag_binary.NewBorshDecoder(metadataAccount.Value.Data.GetBinary())
				err = metadataData.UnmarshalWithDecoder(metadataDecoder)
				if err != nil {
					panic(err)
				}
				_, err = json.MarshalIndent(metadataData, "", "  ")
				if err != nil {
					panic(err)
				}
				fmt.Println(strings.Trim(metadataData.Data.Uri, "\u0000"))
			}
			continue
		}

		return tokens
	}()
}

func treasureCMs() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	treasuryAuthority, _ := GetTreasuryAuthority(oracle.PublicKey())

	candyMachine := solana.MustPublicKeyFromBase58("92uRUkRgGoisXqT81kP3kWUJ1UroGdnXaV9DdURZLNRr")
	candyMachineCreator, _, _ := getCandyMachineCreator(candyMachine)
	treasuryWhitelist, _ := GetTreasuryWhitelist(oracle.PublicKey(), treasuryAuthority, candyMachineCreator)
	treasuryIx := someplace.NewAddWhitelistedCmInstructionBuilder().
		SetCandyMachine(candyMachine).
		SetCandyMachineCreator(candyMachineCreator).
		SetOracleAccount(oracle.PublicKey()).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTreasuryAuthorityAccount(treasuryAuthority).
		SetTreasuryWhitelistAccount(treasuryWhitelist)

	sendTx(
		"sell",
		append(make([]solana.Instruction, 0), treasuryIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}
func treasure() {
	mint := MINT
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	treasuryAuthority, _ := GetTreasuryAuthority(oracle.PublicKey())
	treasuryTokenAccount, _ := GetTreasuryTokenAccount(oracle.PublicKey())
	oracleTokenAccount, _ := getTokenWallet(oracle.PublicKey(), MINT)

	treasuryIx := someplace.NewInitTreasuryInstructionBuilder().
		SetAdornment("fedcoin").
		SetOracleAccount(oracle.PublicKey()).
		SetOracleTokenAccountAccount(oracleTokenAccount).
		SetRentAccount(solana.SysVarRentPubkey).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTokenProgramAccount(solana.TokenProgramID).
		SetTreasuryAuthorityAccount(treasuryAuthority).
		SetTreasuryTokenAccountAccount(treasuryTokenAccount).
		SetTreasuryTokenMintAccount(mint)

	if e := treasuryIx.Validate(); e != nil {
		fmt.Println(e.Error())
		panic("...")
	}

	sendTx(
		"sell",
		append(make([]solana.Instruction, 0), treasuryIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}
func treasureVerifyCM() {
	candyMachine := solana.MustPublicKeyFromBase58("exiVcLT1yPJi2zwkP1pXSS5jSHKTV9UUq5tTtBW6AZW")
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	// candyMachineCreator, _, _ := getCandyMachineCreator(candyMachine)
	treasuryAuthority, _ := GetTreasuryAuthority(oracle.PublicKey())
	treasuryWhitelist, _ := GetTreasuryWhitelist(oracle.PublicKey(), treasuryAuthority, candyMachine)

	GetTreasuryWhitelistData(treasuryWhitelist)
}
func treasureVerify() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	treasuryAuthority, _ := GetTreasuryAuthority(oracle.PublicKey())

	GetTreasuryAuthorityData(treasuryAuthority)
}

func burn() {

	mint := MINT
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./dev.key")
	if err != nil {
		panic(err)
	}

	client := rpc.New(NETWORK)
	tokenAccounts, err := client.GetTokenAccountsByOwner(context.TODO(), oracle.PublicKey(), &rpc.GetTokenAccountsConfig{ProgramId: &solana.TokenProgramID}, &rpc.GetTokenAccountsOpts{Encoding: "jsonParsed"})
	if err != nil {
		panic(err)
	}
	type tokenDataParsed struct {
		Parsed struct {
			Info struct {
				Mint        string `json:"mint"`
				TokenAmount struct {
					UiAmount float64 `json:"uiAmount"`
				} `json:"tokenAmount"`
			} `json:"info"`
		} `json:"parsed"`
	}
	type token struct {
		ata      solana.PublicKey
		mint     solana.PublicKey
		metadata solana.PublicKey
	}

	tokens := func() []token {
		tokens := make([]token, 2)
		for _, tokenAccount := range tokenAccounts.Value {
			_json, err := json.Marshal(tokenAccount.Account.Data.GetRawJSON())
			if err != nil {
				panic(err)
			}
			fmt.Println(string(_json))
			var tokenData tokenDataParsed
			err = json.Unmarshal(tokenAccount.Account.Data.GetRawJSON(), &tokenData)
			if err != nil {
				panic(err)
			}
			switch tokenData.Parsed.Info.Mint {
			case mint.String():
				{
					tokens[0] = token{
						ata:      tokenAccount.Pubkey,
						mint:     solana.MustPublicKeyFromBase58(tokenData.Parsed.Info.Mint),
						metadata: solana.SystemProgramID,
					}
					continue
				}
			default:
				{
					fmt.Println(tokenData.Parsed.Info.TokenAmount.UiAmount, tokenData.Parsed.Info.Mint)
					if tokenData.Parsed.Info.TokenAmount.UiAmount > 0 {
						nftMint := solana.MustPublicKeyFromBase58(tokenData.Parsed.Info.Mint)
						metadata, _ := getMetadata(nftMint)
						tokens[1] = token{
							ata:      tokenAccount.Pubkey,
							mint:     nftMint,
							metadata: metadata,
						}
					}
					continue
				}
			}
		}

		return tokens
	}()
	if !tokens[0].mint.Equals(mint) {
		panic("no mint")
	}

	treasuryAuthority, treasuryAuthorityBump := GetTreasuryAuthority(oracle.PublicKey())
	treasuryTokenAccount, _ := GetTreasuryTokenAccount(oracle.PublicKey())
	sellIx := someplace.NewSellForInstructionBuilder().
		SetDepoMintAccount(tokens[1].mint).
		SetDepoTokenAccountAccount(tokens[1].ata).
		SetInitializerAccount(oracle.PublicKey()).
		SetInitializerTokenAccountAccount(tokens[0].ata).
		SetMetadataAccount(tokens[1].metadata).
		SetTokenProgramAccount(solana.TokenProgramID).
		SetTreasuryAuthorityAccount(treasuryAuthority).
		SetTreasuryBump(treasuryAuthorityBump).
		SetTreasuryTokenAccountAccount(treasuryTokenAccount).
		SetTreasuryTokenMintAccount(mint)

	err = sellIx.Validate()
	if err != nil {
		panic(err)
	}
	sellJsonKeys, _ := json.MarshalIndent(sellIx.Build().Accounts(), "", "  ")
	sellJson, _ := json.MarshalIndent(sellIx.Build(), "", "  ")
	fmt.Println(string(sellJsonKeys))
	fmt.Println(string(sellJson))

	tx := solana.NewTransactionBuilder().AddInstruction(sellIx.Build())
	txB, _ := tx.Build()
	txJson, _ := json.MarshalIndent(txB, "", "  ")
	fmt.Println(string(txJson))

	sendTx(
		"sell",
		append(make([]solana.Instruction, 0), sellIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)

}

func mint() {
	candyMachineAddress := solana.MustPublicKeyFromBase58("Dfp73qrotgTtecSTkakgTX5rfAqW3JtXgwqFMCL26zaz")

	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}

	mint := solana.NewWallet().PrivateKey

	client := rpc.New(NETWORK)
	userTokenAccountAddress, err := getTokenWallet(oracle.PublicKey(), mint.PublicKey())
	if err != nil {
		panic(err)
	}

	candyMachineRaw, err := client.GetAccountInfo(context.TODO(), candyMachineAddress)
	if err != nil {
		panic(err)
	}

	signers := []solana.PrivateKey{mint, oracle}

	min, err := client.GetMinimumBalanceForRentExemption(context.TODO(), token.MINT_SIZE, rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	dec := ag_binary.NewBorshDecoder(candyMachineRaw.Value.Data.GetBinary())
	var cm someplace.Batch
	err = dec.Decode(&cm)
	if err != nil {
		panic(err)
	}

	var instructions []solana.Instruction
	instructions = append(instructions,
		system.NewCreateAccountInstructionBuilder().
			SetOwner(token.ProgramID).
			SetNewAccount(mint.PublicKey()).
			SetSpace(token.MINT_SIZE).
			SetFundingAccount(oracle.PublicKey()).
			SetLamports(min).
			Build(),

		token.NewInitializeMint2InstructionBuilder().
			SetMintAccount(mint.PublicKey()).
			SetDecimals(0).
			SetMintAuthority(oracle.PublicKey()).
			SetFreezeAuthority(oracle.PublicKey()).
			Build(),

		atok.NewCreateInstructionBuilder().
			SetPayer(oracle.PublicKey()).
			SetWallet(oracle.PublicKey()).
			SetMint(mint.PublicKey()).
			Build(),

		token.NewMintToInstructionBuilder().
			SetMintAccount(mint.PublicKey()).
			SetDestinationAccount(userTokenAccountAddress).
			SetAuthorityAccount(oracle.PublicKey()).
			SetAmount(1).
			Build(),
	)

	metadataAddress, err := getMetadata(mint.PublicKey())
	if err != nil {
		panic(err)
	}
	masterEdition, err := getMasterEdition(mint.PublicKey())
	if err != nil {
		panic(err)
	}
	candyMachineCreator, creatorBump, err := getCandyMachineCreator(candyMachineAddress)
	if err != nil {
		panic(err)
	}

	listing, _ := GetListing(oracle.PublicKey(), candyMachineAddress, 0)
	treasuryTokenAccount, _ := GetTreasuryTokenAccount(oracle.PublicKey())
	mintIx := someplace.NewMintNftInstructionBuilder().
		SetConfigIndex(uint64(0)).
		SetCreatorBump(creatorBump).
		SetCandyMachineAccount(candyMachineAddress).
		SetCandyMachineCreatorAccount(candyMachineCreator).
		SetPayerAccount(oracle.PublicKey()).
		SetOracleAccount(cm.Oracle).
		SetMintAccount(mint.PublicKey()).
		SetMetadataAccount(metadataAddress).
		SetMasterEditionAccount(masterEdition).
		SetMintAuthorityAccount(oracle.PublicKey()).
		SetUpdateAuthorityAccount(oracle.PublicKey()).
		SetTreasuryTokenAccountAccount(treasuryTokenAccount).
		SetTokenMetadataProgramAccount(token_metadata.ProgramID).
		SetTokenProgramAccount(token.ProgramID).
		SetListingAccount(listing).
		SetInitializerTokenAccountAccount(solana.MustPublicKeyFromBase58("AfaHA73mAsdFh4ie79smiLFsA5Zm1DCxfWQbd7pBBt7y")).
		SetSystemProgramAccount(system.ProgramID).
		SetRentAccount(solana.SysVarRentPubkey).
		SetClockAccount(solana.SysVarClockPubkey).
		SetInstructionSysvarAccountAccount(solana.SysVarInstructionsPubkey)

	err = mintIx.Validate()
	if err != nil {
		panic(err)
	}
	instructions = append(instructions,
		mintIx.Build(),
	)

	sendTx(
		"mint",
		instructions,
		signers,
		oracle.PublicKey(),
	)

}

func mintRare() {
	candyMachineAddress := solana.MustPublicKeyFromBase58("GTTXtoSYPhaH27tREZ6PZDvKsQMhaCrqyp7rPcTAzjtg")

	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}

	mint := solana.NewWallet().PrivateKey

	client := rpc.New(NETWORK)
	userTokenAccountAddress, err := getTokenWallet(oracle.PublicKey(), mint.PublicKey())
	if err != nil {
		panic(err)
	}

	candyMachineRaw, err := client.GetAccountInfo(context.TODO(), candyMachineAddress)
	if err != nil {
		panic(err)
	}

	signers := []solana.PrivateKey{mint, oracle}

	min, err := client.GetMinimumBalanceForRentExemption(context.TODO(), token.MINT_SIZE, rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	dec := ag_binary.NewBorshDecoder(candyMachineRaw.Value.Data.GetBinary())
	var cm someplace.Batch
	err = dec.Decode(&cm)
	if err != nil {
		panic(err)
	}

	var instructions []solana.Instruction
	instructions = append(instructions,
		system.NewCreateAccountInstructionBuilder().
			SetOwner(token.ProgramID).
			SetNewAccount(mint.PublicKey()).
			SetSpace(token.MINT_SIZE).
			SetFundingAccount(oracle.PublicKey()).
			SetLamports(min).
			Build(),

		token.NewInitializeMint2InstructionBuilder().
			SetMintAccount(mint.PublicKey()).
			SetDecimals(0).
			SetMintAuthority(oracle.PublicKey()).
			SetFreezeAuthority(oracle.PublicKey()).
			Build(),

		atok.NewCreateInstructionBuilder().
			SetPayer(oracle.PublicKey()).
			SetWallet(oracle.PublicKey()).
			SetMint(mint.PublicKey()).
			Build(),

		token.NewMintToInstructionBuilder().
			SetMintAccount(mint.PublicKey()).
			SetDestinationAccount(userTokenAccountAddress).
			SetAuthorityAccount(oracle.PublicKey()).
			SetAmount(1).
			Build(),
	)

	metadataAddress, err := getMetadata(mint.PublicKey())
	if err != nil {
		panic(err)
	}
	masterEdition, err := getMasterEdition(mint.PublicKey())
	if err != nil {
		panic(err)
	}
	candyMachineCreator, creatorBump, err := getCandyMachineCreator(candyMachineAddress)
	if err != nil {
		panic(err)
	}

	listing, _ := GetListing(oracle.PublicKey(), candyMachineAddress, 0)
	treasuryTokenAccount, _ := GetTreasuryTokenAccount(oracle.PublicKey())
	mintIx := someplace.NewMintNftRarityInstructionBuilder().
		SetConfigIndex(uint64(0)).
		SetCreatorBump(creatorBump).
		SetCandyMachineAccount(candyMachineAddress).
		SetCandyMachineCreatorAccount(candyMachineCreator).
		SetPayerAccount(oracle.PublicKey()).
		SetOracleAccount(cm.Oracle).
		SetMintAccount(mint.PublicKey()).
		SetMetadataAccount(metadataAddress).
		SetMasterEditionAccount(masterEdition).
		SetMintAuthorityAccount(oracle.PublicKey()).
		SetUpdateAuthorityAccount(oracle.PublicKey()).
		SetTreasuryTokenAccountAccount(treasuryTokenAccount).
		SetTokenMetadataProgramAccount(token_metadata.ProgramID).
		SetTokenProgramAccount(token.ProgramID).
		SetListingAccount(listing).
		SetInitializerTokenAccountAccount(solana.MustPublicKeyFromBase58("CQ5mZ1Ve4CQK1vQenH1nAnHbyF3MNKavnu1MhJ2Dr4mx")).
		SetSystemProgramAccount(system.ProgramID).
		SetRentAccount(solana.SysVarRentPubkey).
		SetClockAccount(solana.SysVarClockPubkey).
		SetInstructionSysvarAccountAccount(solana.SysVarInstructionsPubkey).
		SetRecentBlockhashesAccount(solana.SysVarRecentBlockHashesPubkey)

	err = mintIx.Validate()
	if err != nil {
		panic(err)
	}
	for _ = range make([]int, 16) {
		mintIx.AccountMetaSlice.Append(&solana.AccountMeta{
			PublicKey:  solana.NewWallet().PublicKey(),
			IsWritable: false,
			IsSigner:   false,
		})
	}
	/*
		sendTx(
			"mint",
			instructions,
			signers,
			oracle.PublicKey(),
		)
	*/

	sendTx(
		"mint",
		append(make([]solana.Instruction, 0),
			mintIx.Build(),
		),
		signers,
		oracle.PublicKey(),
	)

}

func catalogBatches() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}

	batches, _ := GetBatches(oracle.PublicKey())
	batchesData := GetBatchesData(batches)
	for index := range make([]uint64, batchesData.Counter) {
		batchReceipt, _ := GetBatchReceipt(oracle.PublicKey(), uint64(index))
		GetBatchReceiptData(batchReceipt)

	}

}

func enable() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}

	batches, _ := GetBatches(oracle.PublicKey())
	enableIx := someplace.NewEnableBatchUploadingInstructionBuilder().
		SetBatchesAccount(batches).
		SetOracleAccount(oracle.PublicKey()).
		SetSystemProgramAccount(solana.SystemProgramID)

	if e := enableIx.Validate(); e != nil {
		fmt.Println(e.Error())
	}

	sendTx(
		"enable",
		append(make([]solana.Instruction, 0), enableIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)

}

func batchUpload() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	batches, _ := GetBatches(oracle.PublicKey())
	ids := make([]uint64, 2)
	for i := range ids {
		index := GetBatchesData(batches).Counter
		batchReceipt, _ := GetBatchReceipt(oracle.PublicKey(), index)

		batchAccount := solana.NewWallet()
		initBatchAccount := system.NewCreateAccountInstructionBuilder().
			SetFundingAccount(oracle.PublicKey()).
			SetLamports(1 * solana.LAMPORTS_PER_SOL).
			SetNewAccount(batchAccount.PublicKey()).
			SetOwner(someplace.ProgramID).
			SetSpace(1024)

		cmData := someplace.CandyMachineData{}
		cmData.Uuid = "asdf12"

		initCm := someplace.NewInitializeCandyMachineInstructionBuilder().
			SetBatchAccountAccount(batchAccount.PublicKey()).
			SetBatchReceiptAccount(batchReceipt).
			SetBatchesAccount(batches).
			SetName("Integration Test").
			SetData(cmData).
			SetOracleAccount(oracle.PublicKey()).
			SetSystemProgramAccount(solana.SystemProgramID)

		sendTx(
			"init cm",
			append(make([]solana.Instruction, 0), initBatchAccount.Build(), initCm.Build()),
			append(make([]solana.PrivateKey, 0), oracle, batchAccount.PrivateKey),
			oracle.PublicKey(),
		)
		_, _ = initBatchAccount, initCm
		GetBatchesData(batches)

		ids[i] = index
	}

	for _, index := range ids {
		batchReceipt, _ := GetBatchReceipt(oracle.PublicKey(), index)
		GetBatchReceiptData(batchReceipt)

	}

}

func verifyBatchUpload() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	batches, _ := GetBatches(oracle.PublicKey())
	batchesMeta := GetBatchesData(batches)
	var i uint64 = 0
	for i < batchesMeta.Counter {
		batchReceipt, _ := GetBatchReceipt(oracle.PublicKey(), i)
		GetBatchReceiptData(batchReceipt)

		i++
	}

}

func GetTreasuryWhitelistData(treasuryAuthority solana.PublicKey) *someplace.TreasuryWhitelist {
	fmt.Println(someplace.ProgramID)
	rpcClient := rpc.New(NETWORK)
	batchesBin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), treasuryAuthority, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	fmt.Println("....", batchesBin)
	var batchesData someplace.TreasuryWhitelist
	decoder := ag_binary.NewBorshDecoder(batchesBin.Value.Data.GetBinary())
	err := batchesData.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	formatAsJson(batchesData)

	return &batchesData

}
func GetTreasuryAuthorityData(treasuryAuthority solana.PublicKey) *someplace.TreasuryAuthority {
	fmt.Println(someplace.ProgramID)
	rpcClient := rpc.New(NETWORK)
	batchesBin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), treasuryAuthority, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var batchesData someplace.TreasuryAuthority
	decoder := ag_binary.NewBorshDecoder(batchesBin.Value.Data.GetBinary())
	err := batchesData.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	formatAsJson(batchesData)

	return &batchesData

}
func GetBatchesData(batches solana.PublicKey) *someplace.Batches {
	rpcClient := rpc.New(NETWORK)
	batchesBin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), batches, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var batchesData someplace.Batches
	decoder := ag_binary.NewBorshDecoder(batchesBin.Value.Data.GetBinary())
	err := batchesData.UnmarshalWithDecoder(decoder)
	if err != nil {
		panic(err)
	}

	formatAsJson(batchesData)

	return &batchesData

}

func GetBatchReceiptData(batchReceipt solana.PublicKey) *someplace.BatchReceipt {
	rpcClient := rpc.New(NETWORK)
	batchReceiptBin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), batchReceipt, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var batchReceiptData someplace.BatchReceipt
	fmt.Println(batchReceiptBin.Value)
	if batchReceiptBin != nil {
		decoder := ag_binary.NewBorshDecoder(batchReceiptBin.Value.Data.GetBinary())
		err := batchReceiptData.UnmarshalWithDecoder(decoder)
		if err != nil {
			panic(err)
		}

		formatAsJson(batchReceiptData)
	}

	return &batchReceiptData

}
func GetListingData(listing solana.PublicKey) *someplace.Listing {
	rpcClient := rpc.New(NETWORK)
	batchReceiptBin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), listing, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var batchReceiptData someplace.Listing
	fmt.Println(batchReceiptBin.Value)
	if batchReceiptBin != nil {
		decoder := ag_binary.NewBorshDecoder(batchReceiptBin.Value.Data.GetBinary())
		err := batchReceiptData.UnmarshalWithDecoder(decoder)
		if err != nil {
			panic(err)
		}

		formatAsJson(batchReceiptData)
	}

	return &batchReceiptData

}

func formatAsJson(data interface{}) {
	dataJson, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(dataJson))
}

func sendTx(
	doc string,
	instructions []solana.Instruction,
	signers []solana.PrivateKey,
	feePayer solana.PublicKey,
) {
	rpcClient := rpc.New(NETWORK)
	wsClient, err := ws.Connect(context.TODO(), "wss://api.testnet.solana.com")
	if err != nil {
		log.Println("PANIC!!!", fmt.Errorf("unable to open WebSocket Client - %w", err))
	}

	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		log.Println("PANIC!!!", fmt.Errorf("unable to fetch recent blockhash - %w", err))
		return
	}

	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(feePayer),
	)
	if err != nil {
		log.Println("PANIC!!!", fmt.Errorf("unable to create transaction"))
		return
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		for _, candidate := range signers {
			if candidate.PublicKey().Equals(key) {
				return &candidate
			}
		}
		return nil
	})
	if err != nil {
		log.Println("PANIC!!!", fmt.Errorf("unable to sign transaction: %w", err))
		return
	}

	tx.EncodeTree(text.NewTreeEncoder(os.Stdout, doc))

	bin, _ := tx.MarshalBinary()
	sig, err := sendAndConfirmTransaction.SendAndConfirmTransaction(
		context.TODO(),
		rpcClient,
		wsClient,
		tx,
	)
	fmt.Println(len(bin))
	if err != nil {
		log.Println("PANIC!!!", fmt.Errorf("unable to send transaction - %w", err))
		return
	}
	wsClient.Close()
	log.Println(sig)
}

func GetBatchReceipt(
	oracle solana.PublicKey,
	index uint64,
) (solana.PublicKey, uint8) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, index)
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			oracle.Bytes(),
			buf,
		},
		someplace.ProgramID,
	)
	fmt.Println("Someplace PDA - batch receipt", oracle, index)
	return addr, bump
}

func GetBatches(
	oracle solana.PublicKey,
) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			oracle.Bytes(),
		},
		someplace.ProgramID,
	)
	fmt.Println("Someplace PDA - batches", addr, oracle)
	return addr, bump
}

func GetListing(oracle, batch solana.PublicKey, configIndex uint64) (solana.PublicKey, uint8) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, configIndex)
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			oracle.Bytes(),
			batch.Bytes(),
			buf,
		},
		someplace.ProgramID,
	)
	return addr, bump
}
func GetTreasuryWhitelist(oracle, treasuryAuthority, candyMachine solana.PublicKey) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			oracle.Bytes(),
			[]byte("treasury_whitelist"),
			treasuryAuthority.Bytes(),
			candyMachine.Bytes(),
		},
		someplace.ProgramID,
	)
	return addr, bump
}
func GetTreasuryAuthority(oracle solana.PublicKey) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("someplace"),
			[]byte("ballz"),
			oracle.Bytes(),
		},
		someplace.ProgramID,
	)
	return addr, bump
}
func GetTreasuryTokenAccount(oracle solana.PublicKey) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("someplace"),
			[]byte("ballz"),
			[]byte("treasury_mint"),
			oracle.Bytes(),
		},
		someplace.ProgramID,
	)
	return addr, bump
}

func getTokenWallet(wallet solana.PublicKey, mint solana.PublicKey) (solana.PublicKey, error) {
	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			wallet.Bytes(),
			solana.TokenProgramID.Bytes(),
			mint.Bytes(),
		},
		solana.SPLAssociatedTokenAccountProgramID,
	)
	return addr, err
}

func getCandyMachineCreator(candyMachineAddress solana.PublicKey) (solana.PublicKey, uint8, error) {
	return solana.FindProgramAddress(
		[][]byte{
			[]byte("someplace"),
			candyMachineAddress.Bytes(),
		},
		someplace.ProgramID,
	)
}

func getMetadata(mint solana.PublicKey) (solana.PublicKey, error) {
	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("metadata"),
			token_metadata.ProgramID.Bytes(),
			mint.Bytes(),
		},
		token_metadata.ProgramID,
	)
	return addr, err
}

func getMasterEdition(mint solana.PublicKey) (solana.PublicKey, error) {
	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("metadata"),
			token_metadata.ProgramID.Bytes(),
			mint.Bytes(),
			[]byte("edition"),
		},
		token_metadata.ProgramID,
	)
	return addr, err
}

func GetMarketAuthority(oracle, marketUid solana.PublicKey) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("someplace"),
			[]byte("market"),
			oracle.Bytes(),
			marketUid.Bytes(),
		},
		someplace.ProgramID,
	)
	return addr, bump
}

func GetMarketAuthorityData(marketAuthority solana.PublicKey) *someplace.Market {
	rpcClient := rpc.New(NETWORK)
	batchReceiptBin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), marketAuthority, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var batchReceiptData someplace.Market
	fmt.Println(batchReceiptBin.Value)
	if batchReceiptBin != nil {
		decoder := ag_binary.NewBorshDecoder(batchReceiptBin.Value.Data.GetBinary())
		err := batchReceiptData.UnmarshalWithDecoder(decoder)
		if err != nil {
			panic(err)
		}

		formatAsJson(batchReceiptData)
	}

	return &batchReceiptData

}

/*
   seeds = [PREFIX.as_ref(), LISTING.as_ref(), market_authority.key().as_ref(), market_authority.listings.to_le_bytes().as_ref()],
   seeds = [PREFIX.as_ref(), LISTINGTOKEN.as_ref(), market_authority.key().as_ref(), index.to_le_bytes().as_ref()],
*/
func GetMarketListing(marketAuthority solana.PublicKey, index uint64) (solana.PublicKey, uint8) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, index)
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("someplace"),
			[]byte("publiclisting"),
			marketAuthority.Bytes(),
			buf,
		},
		someplace.ProgramID,
	)
	return addr, bump
}
func GetMarketListingData(marketListing solana.PublicKey) *someplace.MarketListing {
	rpcClient := rpc.New(NETWORK)
	batchReceiptBin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), marketListing, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var batchReceiptData someplace.MarketListing
	fmt.Println(batchReceiptBin.Value)
	if batchReceiptBin != nil {
		decoder := ag_binary.NewBorshDecoder(batchReceiptBin.Value.Data.GetBinary())
		err := batchReceiptData.UnmarshalWithDecoder(decoder)
		if err != nil {
			panic(err)
		}

		formatAsJson(batchReceiptData)
	}

	return &batchReceiptData

}
func GetMarketListingTokenAccount(marketAuthority solana.PublicKey, index uint64) (solana.PublicKey, uint8) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, index)
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("someplace"),
			[]byte("listingtoken"),
			marketAuthority.Bytes(),
			buf,
		},
		someplace.ProgramID,
	)
	return addr, bump
}

type TokenListMeta struct {
	Address solana.PublicKey `json:"address"`
	Symbol  string           `json:"symbol"`
	Name    string           `json:"name"`
}

type TokenList struct {
	Tokens []TokenListMeta `json:"tokens"`
}

func GetMarketMintMeta() {
	marketUid := solana.MustPublicKeyFromBase58("4Gm324iNEMapZV9aVyWg8EwJYLiqepYYab47sCWcPnh1")
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	marketAuthority, _ := GetMarketAuthority(oracle.PublicKey(), marketUid)
	marketAuthorityData := GetMarketAuthorityData(marketAuthority)

	var tokenMeta TokenListMeta
	tokens := FetchTokenMeta()
	for _, token := range tokens {
		if token.Address.Equals(marketAuthorityData.MarketMint) {
			tokenMeta = token
		}
	}

	fmt.Println(tokenMeta)
}

func FetchTokenMeta() []TokenListMeta {
	var tokenList TokenList
	tokenListUrl := fmt.Sprint(CDN + "/solana.tokenlist.json")
	res, err := http.DefaultClient.Get(tokenListUrl)
	if err != nil {
		return tokenList.Tokens
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return tokenList.Tokens
	}

	err = json.Unmarshal(data, &tokenList)
	if err != nil {
		return tokenList.Tokens
	}

	return tokenList.Tokens

}

func GetMarketListingsData() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	marketUid := solana.MustPublicKeyFromBase58("4Gm324iNEMapZV9aVyWg8EwJYLiqepYYab47sCWcPnh1")
	marketAuthority, _ := GetMarketAuthority(oracle.PublicKey(), marketUid)
	marketAuthorityData := GetMarketAuthorityData(marketAuthority)
	var i uint64 = 0
	for i < marketAuthorityData.Listings {
		batchReceipt, _ := GetMarketListing(marketAuthority, i)
		GetMarketListingData(batchReceipt)

		i++
	}
}
