package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/token"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	quest_ops "creaturez.nft/questing/quests/ops"
	"creaturez.nft/someplace"
	"creaturez.nft/someplace/marketplace"
	"creaturez.nft/someplace/storefront"
	"creaturez.nft/utils"

	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
	atok "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

const DEVNET = "https://sparkling-dark-shadow.solana-devnet.quiknode.pro/0e9964e4d70fe7f856e7d03bc7e41dc6a2b84452/"
const TESTNET = "https://api.testnet.solana.com"
const NETWORK = DEVNET

var MINT = solana.MustPublicKeyFromBase58("3XeA57a7vZSD6krvu5sGUN3A5SCKJmYadeoiQ8YmACbz")

func init() {
	questing.SetProgramID(solana.MustPublicKeyFromBase58("Cv1EGc9jnop3n1fhCNBRSpoQ98uTdGCi4ESAxAwh2Ek5"))
	someplace.SetProgramID(solana.MustPublicKeyFromBase58("H4FUWoKkGC7w9B6Ef3bHScPkfHMHA4wyuXapQAA5gLrP"))
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
	GetListingsForOracle()
	/*
		candyMachineAddresses := []solana.PublicKey{solana.MustPublicKeyFromBase58("2h1BiLfvU1DH2kyYDZ7KuQh81hSDU1iEQ5ZAbjqDjH8C"), solana.MustPublicKeyFromBase58("Ek4cSAY9JfJFaiVk5MQBZ9yWZiMs5ZHTjEX8Nbm3VutN")}
		for _, candyMachine := range candyMachineAddresses {
			for range make([]int, 10) {
				mint(candyMachine)
			}
		}
	*/
	// mintRare()
	// holder_nft_metadata()
	// burn()

	// marketCreate()
	// verifyMarketCreate()
	// marketList()
	// verifyMarketList()
	// marketFulfill()

	// GetMarketMintMeta()
	// GetMarketListingsData()

	// CreateNTokenAccountsOfMint(MINT, 3)

	// enableVias()
	// enableViaForRarityToken()

	// startQuest()
	// startAndEndQuest()
	// ETZoY7cJfD8N7EVx5tShRYS1vxgv3F4Dkavjb52kGRyj
	// enableQuestsAndCreateQuest()
	// CreateAndAmmendEntitlementQuest()
	// treasureVerify()

}

func treasureVerify() {
	// BAP4H9Qki6GFtVjoDhWEFgBh9DwrUbMYrAdHqzu7a9nf
	treasuryData := storefront.GetTreasuryAuthorityData(solana.MustPublicKeyFromBase58("BAP4H9Qki6GFtVjoDhWEFgBh9DwrUbMYrAdHqzu7a9nf"))
	js, _ := json.Marshal(treasuryData)

	fmt.Println(string(js))
}

func CreateNTokenAccountsOfMint(mint solana.PublicKey, amount int) {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	tokenAccounts := make([]string, amount)
	var instructions []solana.Instruction
	for i := range tokenAccounts {
		wallet := solana.NewWallet()
		ata, _ := utils.GetTokenWallet(wallet.PublicKey(), mint)
		tokenAccounts[i] = ata.String()

		instructions = append(instructions,
			atok.NewCreateInstructionBuilder().
				SetPayer(oracle.PublicKey()).
				SetWallet(wallet.PublicKey()).
				SetMint(mint).
				Build(),
		)

	}
	utils.SendTx(
		"list",
		instructions,
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
	fmt.Println(tokenAccounts)
}

func marketList() {
	marketUid := solana.MustPublicKeyFromBase58("GAm4cGVVMi5NMBVoxg8QhuKpQk2xW4BVbhnQrESE54HA")
	marketMint := MINT
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	marketAuthority, _ := marketplace.GetMarketAuthority(oracle.PublicKey(), marketUid)
	marketAuthorityData := marketplace.GetMarketAuthorityData(marketAuthority)
	nftMint := solana.NewWallet().PrivateKey

	userTokenAccountAddress, _ := utils.GetTokenWallet(oracle.PublicKey(), nftMint.PublicKey())
	sellerMarketTokenAccountAddress, _ := utils.GetTokenWallet(oracle.PublicKey(), marketMint)

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

	marketListing, _ := marketplace.GetMarketListing(marketAuthority, marketAuthorityData.Listings)
	marketListingTokenAccount, _ := marketplace.GetMarketListingTokenAccount(marketAuthority, marketAuthorityData.Listings)
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

	utils.SendTx(
		"list",
		append(instructions, listIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle, nftMint),
		oracle.PublicKey(),
	)

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
	marketAuthority, marketAuthorityBump := marketplace.GetMarketAuthority(oracle.PublicKey(), marketUid)
	// marketAuthorityData := GetMarketAuthorityData(marketAuthority)
	marketListing, _ := marketplace.GetMarketListing(marketAuthority, 1)
	marketListingData := marketplace.GetMarketListingData(marketListing)
	marketMint := MINT
	buyerMarketTokenAccountAddress, _ := utils.GetTokenWallet(buyer.PublicKey(), marketMint)
	buyerNftTokenAccountAddress, _ := utils.GetTokenWallet(buyer.PublicKey(), marketListingData.NftMint)

	marketListingTokenAccount, _ := marketplace.GetMarketListingTokenAccount(marketAuthority, marketListingData.Index)
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

	utils.SendTx(
		"list",
		append(instructions, listIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle, buyer),
		buyer.PublicKey(),
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
				metadata, _ := utils.GetMetadata(nftMint)
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
						metadata, _ := utils.GetMetadata(nftMint)
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

	treasuryAuthority, treasuryAuthorityBump := storefront.GetTreasuryAuthority(oracle.PublicKey())
	treasuryTokenAccount, _ := storefront.GetTreasuryTokenAccount(oracle.PublicKey())
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

	utils.SendTx(
		"sell",
		append(make([]solana.Instruction, 0), sellIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)

}

func mint(candyMachineAddress solana.PublicKey) {

	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	initializerTokenAccount, _ := utils.GetTokenWallet(oracle.PublicKey(), MINT)

	treasury, _ := storefront.GetTreasuryAuthority(oracle.PublicKey())
	treasuryData := storefront.GetTreasuryAuthorityData(treasury)

	client := rpc.New(NETWORK)

	candyMachineRaw, err := client.GetAccountInfo(context.TODO(), candyMachineAddress)
	if err != nil {
		panic(err)
	}

	// signers :=

	// min, err := client.GetMinimumBalanceForRentExemption(context.TODO(), token.MINT_SIZE, rpc.CommitmentFinalized)
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

	listing, _ := storefront.GetListing(oracle.PublicKey(), candyMachineAddress, 1)
	listingData := storefront.GetListingData(listing)
	js, _ := json.MarshalIndent(treasuryData, "", "  ")
	fmt.Println(string(js))
	mint, _, _ := storefront.GetMint(oracle.PublicKey(), listing, listingData.Mints)
	/*
		    mintAta, _ := utils.GetTokenWallet(oracle.PublicKey(), mint)
			if err != nil {
				panic(err)
			}
	*/
	mintAta := solana.NewWallet()
	mintHash, _, _ := storefront.GetMintHash(oracle.PublicKey(), listing, listingData.Mints)
	metadataAddress, err := utils.GetMetadata(mint)
	if err != nil {
		panic(err)
	}
	masterEdition, err := utils.GetMasterEdition(mint)
	if err != nil {
		panic(err)
	}
	candyMachineCreator, creatorBump, err := storefront.GetCandyMachineCreator(candyMachineAddress)
	if err != nil {
		panic(err)
	}

	mintIx := someplace.NewMintNftInstructionBuilder().
		SetConfigIndex(uint64(1)).
		SetCreatorBump(creatorBump).
		SetMintHashAccount(mintHash).
		SetCandyMachineAccount(candyMachineAddress).
		SetCandyMachineCreatorAccount(candyMachineCreator).
		SetPayerAccount(oracle.PublicKey()).
		SetOracleAccount(cm.Oracle).
		SetMintAccount(mint).
		SetMintAtaAccount(mintAta.PublicKey()).
		SetMetadataAccount(metadataAddress).
		SetMasterEditionAccount(masterEdition).
		SetTreasuryAuthorityAccount(treasury).
		SetTokenMetadataProgramAccount(token_metadata.ProgramID).
		SetTokenProgramAccount(token.ProgramID).
		SetListingAccount(listing).
		SetInitializerTokenAccountAccount(initializerTokenAccount).
		SetSystemProgramAccount(system.ProgramID).
		SetRentAccount(solana.SysVarRentPubkey).
		SetClockAccount(solana.SysVarClockPubkey).
		SetInstructionSysvarAccountAccount(solana.SysVarInstructionsPubkey)

	err = mintIx.Validate()
	if err != nil {
		panic(err)
	}

	for _, split := range treasuryData.Splits {
		mintIx.Append(solana.NewAccountMeta(split.TokenAddress, true, false))
	}
	for range make([]int, 10) {
		mintIx.Append(&solana.AccountMeta{
			PublicKey:  solana.NewWallet().PublicKey(),
			IsWritable: false,
			IsSigner:   false,
		})
	}

	instructions = append(instructions,
		mintIx.Build(),
	)

	utils.SendTx(
		"mint",
		instructions,
		[]solana.PrivateKey{oracle, mintAta.PrivateKey},
		oracle.PublicKey(),
	)

}

func mintRare() {
	candyMachineAddress := solana.MustPublicKeyFromBase58("BmroeEu5zY7KvRGo2FsQ2dJM7EkvnLNdAktzPxVWPb1b")

	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}

	mint := solana.NewWallet().PrivateKey

	client := rpc.New(NETWORK)

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
	)

	/*
		        userTokenAccountAddress, err := utils.GetTokenWallet(oracle.PublicKey(), mint.PublicKey())
		        if err != nil {
		            panic(err)
		        }
				metadataAddress, err := utils.GetMetadata(mint.PublicKey())
				if err != nil {
					panic(err)
				}
				masterEdition, err := utils.GetMasterEdition(mint.PublicKey())
				if err != nil {
					panic(err)
				}
				candyMachineCreator, creatorBump, err := storefront.GetCandyMachineCreator(candyMachineAddress)
				if err != nil {
					panic(err)
				}

				treasuryTokenAccount, _ := storefront.GetTreasuryTokenAccount(oracle.PublicKey())
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
						SetTreasuryTokenAccountAccount(treasuryTokenAccount).
						SetTokenMetadataProgramAccount(token_metadata.ProgramID).
						SetTokenProgramAccount(token.ProgramID).
						SetInitializerTokenAccountAccount(solana.MustPublicKeyFromBase58("CQ5mZ1Ve4CQK1vQenH1nAnHbyF3MNKavnu1MhJ2Dr4mx")).
						SetNftTokenAccountAccount(userTokenAccountAddress).
						SetSystemProgramAccount(system.ProgramID).
						SetRentAccount(solana.SysVarRentPubkey).
						SetClockAccount(solana.SysVarClockPubkey).
						SetInstructionSysvarAccountAccount(solana.SysVarInstructionsPubkey).
						SetRecentBlockhashesAccount(solana.SysVarRecentBlockHashesPubkey)

					err = mintIx.Validate()
					if err != nil {
						panic(err)
					}
					for range make([]int, 9) {
						mintIx.Append(&solana.AccountMeta{
							PublicKey:  solana.NewWallet().PublicKey(),
							IsWritable: false,
							IsSigner:   false,
						})
					}
	*/

	/*
		sendTx(
			"mint",
			instructions,
			signers,
			oracle.PublicKey(),
		)
			// mintIx.Build(),
	*/
	utils.SendTx(
		"mint",
		append(
			instructions,
		),
		signers,
		oracle.PublicKey(),
	)

}

func batchUpload() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	batches, _ := storefront.GetBatches(oracle.PublicKey())
	ids := make([]uint64, 2)
	for i := range ids {
		index := storefront.GetBatchesData(batches).Counter
		batchReceipt, _ := storefront.GetBatchReceipt(oracle.PublicKey(), index)

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

		utils.SendTx(
			"init cm",
			append(make([]solana.Instruction, 0), initBatchAccount.Build(), initCm.Build()),
			append(make([]solana.PrivateKey, 0), oracle, batchAccount.PrivateKey),
			oracle.PublicKey(),
		)
		_, _ = initBatchAccount, initCm
		storefront.GetBatchesData(batches)

		ids[i] = index
	}

	for _, index := range ids {
		batchReceipt, _ := storefront.GetBatchReceipt(oracle.PublicKey(), index)
		storefront.GetBatchReceiptData(batchReceipt)

	}

}

func formatAsJson(data interface{}) {
	dataJson, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(dataJson))
}

func GetMarketMintMeta() {
	marketUid := solana.MustPublicKeyFromBase58("4Gm324iNEMapZV9aVyWg8EwJYLiqepYYab47sCWcPnh1")
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	marketAuthority, _ := marketplace.GetMarketAuthority(oracle.PublicKey(), marketUid)
	marketAuthorityData := marketplace.GetMarketAuthorityData(marketAuthority)

	var tokenMeta utils.TokenListMeta
	tokens := utils.FetchTokenMeta()
	for _, token := range tokens {
		if token.Address.Equals(marketAuthorityData.MarketMint) {
			tokenMeta = token
		}
	}

	fmt.Println(tokenMeta)
}

func GetMarketListingsData() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	marketUid := solana.MustPublicKeyFromBase58("4Gm324iNEMapZV9aVyWg8EwJYLiqepYYab47sCWcPnh1")
	marketAuthority, _ := marketplace.GetMarketAuthority(oracle.PublicKey(), marketUid)
	marketAuthorityData := marketplace.GetMarketAuthorityData(marketAuthority)
	var i uint64 = 0
	for i < marketAuthorityData.Listings {
		batchReceipt, _ := marketplace.GetMarketListing(marketAuthority, i)
		marketplace.GetMarketListingData(batchReceipt)

		i++
	}
}

func enableVias() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	treasuryAuthority, _ := storefront.GetTreasuryAuthority(oracle.PublicKey())
	vias, _ := storefront.GetVias(oracle.PublicKey())
	var instructions []solana.Instruction

	instructions = append(instructions,
		someplace.NewEnableViasInstructionBuilder().
			SetOracleAccount(oracle.PublicKey()).
			SetSystemProgramAccount(solana.SystemProgramID).
			SetTreasuryAuthorityAccount(treasuryAuthority).
			SetViasAccount(vias).
			Build(),
	)

	utils.SendTx(
		"list",
		instructions,
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}

func enableViaForRarityToken() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	treasuryAuthority, _ := storefront.GetTreasuryAuthority(oracle.PublicKey())
	nftMint := solana.NewWallet().PrivateKey

	userTokenAccountAddress, _ := utils.GetTokenWallet(oracle.PublicKey(), nftMint.PublicKey())

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

	vias, _ := storefront.GetVias(oracle.PublicKey())
	viasData := storefront.GetViasData(vias)
	via, _ := storefront.GetVia(oracle.PublicKey(), viasData.Vias)
	viaMapping, _ := storefront.GetViaMapping(oracle.PublicKey(), nftMint.PublicKey())
	instructions = append(
		instructions,
		someplace.NewEnableViaRarityTokenMintingInstructionBuilder().
			SetOracleAccount(oracle.PublicKey()).
			SetRarity("rare").
			SetRarityTokenMintAccount(nftMint.PublicKey()).
			SetSystemProgramAccount(solana.SystemProgramID).
			SetTreasuryAuthorityAccount(treasuryAuthority).
			SetViaAccount(via).
			SetViaMappingAccount(viaMapping).
			SetViasAccount(vias).
			Build(),
	)

	utils.SendTx(
		"list",
		instructions,
		append(make([]solana.PrivateKey, 0), oracle, nftMint),
		oracle.PublicKey(),
	)
}

func enableQuestsAndCreateQuest() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}

	quest_ops.EnableQuests(oracle)

	quest_ops.CreateQuest(oracle, questing.Quest{
		Duration:        100,
		Oracle:          oracle.PublicKey(),
		WlCandyMachines: []solana.PublicKey{oracle.PublicKey()},
		Rewards: []questing.Reward{
			{
				MintAddress: solana.NewWallet().PublicKey(),
				Amount:      1,
			},
		},
		Tender: &questing.Tender{
			MintAddress: solana.MustPublicKeyFromBase58("3XeA57a7vZSD6krvu5sGUN3A5SCKJmYadeoiQ8YmACbz"),
			Amount:      5,
		},
	})
}
func CreateAndAmmendEntitlementQuest() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}

	quest_ops.CreateQuest(oracle, questing.Quest{
		Duration:        100,
		Oracle:          oracle.PublicKey(),
		WlCandyMachines: []solana.PublicKey{oracle.PublicKey()},
		Rewards: []questing.Reward{
			{
				MintAddress: solana.NewWallet().PublicKey(),
				Amount:      1,
			},
		},
		Tender: &questing.Tender{
			MintAddress: solana.MustPublicKeyFromBase58("3XeA57a7vZSD6krvu5sGUN3A5SCKJmYadeoiQ8YmACbz"),
			Amount:      5,
		},
	})

	{
		questsPda, _ := quests.GetQuests(oracle.PublicKey())
		questsData := quests.GetQuestsData(questsPda)
		quest, _ := quests.GetQuest(oracle.PublicKey(), questsData.Quests-1)
		fmt.Println(quest, questsData.Quests)
		questData := quests.GetQuestData(quest)
		{
			questDataJson, _ := json.MarshalIndent(questData, "", "  ")
			fmt.Println(string(questDataJson))
		}

		{
			quest_ops.AmmendQuestWithEntitlement(
				oracle,
				*questData,
				questing.Reward{
					MintAddress: solana.MustPublicKeyFromBase58("ETZoY7cJfD8N7EVx5tShRYS1vxgv3F4Dkavjb52kGRyj"),
					Amount:      50,
				},
			)
		}
	}
	{
		questsPda, _ := quests.GetQuests(oracle.PublicKey())
		questsData := quests.GetQuestsData(questsPda)
		quest, _ := quests.GetQuest(oracle.PublicKey(), questsData.Quests-1)
		questData := quests.GetQuestData(quest)
		{
			questDataJson, _ := json.MarshalIndent(questData, "", "  ")
			fmt.Println(string(questDataJson))
		}
	}
}

/*
func startQuest() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	var instructions []solana.Instruction

	// marketUid := solana.NewWallet().PublicKey()
	pixelBallMint := solana.NewWallet().PrivateKey
	pixelBallTokenAddress, _ := utils.GetTokenWallet(oracle.PublicKey(), pixelBallMint.PublicKey())
	ballzMint := solana.NewWallet().PrivateKey
	ballzTokenAddress, _ := utils.GetTokenWallet(oracle.PublicKey(), ballzMint.PublicKey())
	{

		client := rpc.New(NETWORK)
		min, err := client.GetMinimumBalanceForRentExemption(context.TODO(), token.MINT_SIZE, rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}

		instructions = append(instructions,
			system.NewCreateAccountInstructionBuilder().
				SetOwner(token.ProgramID).
				SetNewAccount(ballzMint.PublicKey()).
				SetSpace(token.MINT_SIZE).
				SetFundingAccount(oracle.PublicKey()).
				SetLamports(min).
				Build(),

			token.NewInitializeMint2InstructionBuilder().
				SetMintAccount(ballzMint.PublicKey()).
				SetDecimals(0).
				SetMintAuthority(oracle.PublicKey()).
				SetFreezeAuthority(oracle.PublicKey()).
				Build(),

			atok.NewCreateInstructionBuilder().
				SetPayer(oracle.PublicKey()).
				SetWallet(oracle.PublicKey()).
				SetMint(ballzMint.PublicKey()).
				Build(),

			token.NewMintToInstructionBuilder().
				SetMintAccount(ballzMint.PublicKey()).
				SetDestinationAccount(ballzTokenAddress).
				SetAuthorityAccount(oracle.PublicKey()).
				SetAmount(3000).
				Build(),
		)
	}
	{

		client := rpc.New(NETWORK)
		min, err := client.GetMinimumBalanceForRentExemption(context.TODO(), token.MINT_SIZE, rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}

		instructions = append(instructions,
			system.NewCreateAccountInstructionBuilder().
				SetOwner(token.ProgramID).
				SetNewAccount(pixelBallMint.PublicKey()).
				SetSpace(token.MINT_SIZE).
				SetFundingAccount(oracle.PublicKey()).
				SetLamports(min).
				Build(),

			token.NewInitializeMint2InstructionBuilder().
				SetMintAccount(pixelBallMint.PublicKey()).
				SetDecimals(0).
				SetMintAuthority(oracle.PublicKey()).
				SetFreezeAuthority(oracle.PublicKey()).
				Build(),

			atok.NewCreateInstructionBuilder().
				SetPayer(oracle.PublicKey()).
				SetWallet(oracle.PublicKey()).
				SetMint(pixelBallMint.PublicKey()).
				Build(),

			token.NewMintToInstructionBuilder().
				SetMintAccount(pixelBallMint.PublicKey()).
				SetDestinationAccount(pixelBallTokenAddress).
				SetAuthorityAccount(oracle.PublicKey()).
				SetAmount(1).
				Build(),
		)
	}
	utils.SendTx(
		"list",
		instructions,
		append(make([]solana.PrivateKey, 0), oracle, ballzMint, pixelBallMint),
		oracle.PublicKey(),
	)

	quest, _ := GetQuestAccount(oracle.PublicKey(), 0)
	questDeposit, _ := GetDepositTokenAccount(oracle.PublicKey(), 0)
	startQuestIx := questing.NewStartQuestInstructionBuilder().
		SetBallzMintAccount(ballzMint.PublicKey()).
		SetBallzTokenAccountAccount(ballzTokenAddress).
		SetDepositTokenAccountAccount(questDeposit).
		SetInitializerAccount(oracle.PublicKey()).
		SetPixelballzMintAccount(pixelBallMint.PublicKey()).
		SetPixelballzTokenAccountAccount(*pixelBallTokenAddress.ToPointer()).
		SetQuestAccountAccount(quest).
		SetRentAccount(solana.SysVarRentPubkey).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTokenProgramAccount(solana.TokenProgramID)

	utils.SendTx(
		"init cm",
		append(make([]solana.Instruction, 0), startQuestIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}
*/

/*
func startAndEndQuest() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	var instructions []solana.Instruction

	// marketUid := solana.NewWallet().PublicKey()
	pixelBallMint := solana.NewWallet().PrivateKey
	pixelBallTokenAddress, _ := utils.GetTokenWallet(oracle.PublicKey(), pixelBallMint.PublicKey())
	ballzMint := solana.NewWallet().PrivateKey
	ballzTokenAddress, _ := utils.GetTokenWallet(oracle.PublicKey(), ballzMint.PublicKey())
	{

		client := rpc.New(NETWORK)
		min, err := client.GetMinimumBalanceForRentExemption(context.TODO(), token.MINT_SIZE, rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}

		instructions = append(instructions,
			system.NewCreateAccountInstructionBuilder().
				SetOwner(token.ProgramID).
				SetNewAccount(ballzMint.PublicKey()).
				SetSpace(token.MINT_SIZE).
				SetFundingAccount(oracle.PublicKey()).
				SetLamports(min).
				Build(),

			token.NewInitializeMint2InstructionBuilder().
				SetMintAccount(ballzMint.PublicKey()).
				SetDecimals(0).
				SetMintAuthority(oracle.PublicKey()).
				SetFreezeAuthority(oracle.PublicKey()).
				Build(),

			atok.NewCreateInstructionBuilder().
				SetPayer(oracle.PublicKey()).
				SetWallet(oracle.PublicKey()).
				SetMint(ballzMint.PublicKey()).
				Build(),

			token.NewMintToInstructionBuilder().
				SetMintAccount(ballzMint.PublicKey()).
				SetDestinationAccount(ballzTokenAddress).
				SetAuthorityAccount(oracle.PublicKey()).
				SetAmount(3000).
				Build(),
		)
	}
	{

		client := rpc.New(NETWORK)
		min, err := client.GetMinimumBalanceForRentExemption(context.TODO(), token.MINT_SIZE, rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}

		instructions = append(instructions,
			system.NewCreateAccountInstructionBuilder().
				SetOwner(token.ProgramID).
				SetNewAccount(pixelBallMint.PublicKey()).
				SetSpace(token.MINT_SIZE).
				SetFundingAccount(oracle.PublicKey()).
				SetLamports(min).
				Build(),

			token.NewInitializeMint2InstructionBuilder().
				SetMintAccount(pixelBallMint.PublicKey()).
				SetDecimals(0).
				SetMintAuthority(oracle.PublicKey()).
				SetFreezeAuthority(oracle.PublicKey()).
				Build(),

			atok.NewCreateInstructionBuilder().
				SetPayer(oracle.PublicKey()).
				SetWallet(oracle.PublicKey()).
				SetMint(pixelBallMint.PublicKey()).
				Build(),

			token.NewMintToInstructionBuilder().
				SetMintAccount(pixelBallMint.PublicKey()).
				SetDestinationAccount(pixelBallTokenAddress).
				SetAuthorityAccount(oracle.PublicKey()).
				SetAmount(1).
				Build(),
		)
	}
	utils.SendTx(
		"list",
		instructions,
		append(make([]solana.PrivateKey, 0), oracle, ballzMint, pixelBallMint),
		oracle.PublicKey(),
	)

	quest, _ := GetQuestAccount(oracle.PublicKey(), 1)
	questDeposit, questDepositBump := GetDepositTokenAccount(oracle.PublicKey(), 1)
	startQuestIx := questing.NewStartQuestInstructionBuilder().
		SetBallzMintAccount(ballzMint.PublicKey()).
		SetBallzTokenAccountAccount(ballzTokenAddress).
		SetDepositTokenAccountAccount(questDeposit).
		SetInitializerAccount(oracle.PublicKey()).
		SetPixelballzMintAccount(pixelBallMint.PublicKey()).
		SetPixelballzTokenAccountAccount(*pixelBallTokenAddress.ToPointer()).
		SetQuestAccountAccount(quest).
		SetQuestIndex(1).
		SetRentAccount(solana.SysVarRentPubkey).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTokenProgramAccount(solana.TokenProgramID)

	utils.SendTx(
		"init cm",
		append(make([]solana.Instruction, 0), startQuestIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)

	endQuestIx := questing.NewEndQuestInstructionBuilder().
		SetDepositTokenAccountAccount(questDeposit).
		SetDepositTokenAccountBump(questDepositBump).
		SetInitializerAccount(oracle.PublicKey()).
		SetPixelballzMintAccount(pixelBallMint.PublicKey()).
		SetPixelballzTokenAccountAccount(pixelBallTokenAddress).
		SetQuestAccountAccount(quest).
		SetQuestIndex(1).
		SetTokenProgramAccount(solana.TokenProgramID)

	utils.SendTx(
		"init cm",
		append(make([]solana.Instruction, 0), endQuestIx.Build()),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)
}
*/

func GetListingsForOracle() {
	/*
		oracle := solana.MustPublicKeyFromBase58("Hj4fk2SNGHFa33Qwi9nkpuxE5yBs2cvpxXTRgLdXdWiP")

		batch := solana.MustPublicKeyFromBase58("Hw8aHQ9YE8YjypjgMmhrHZicp3mFhLmjxCN17sJ9knET")
		for i := range make([]int, 5) {
			listingPda, _ := storefront.GetListing(oracle, batch, uint64(i))
			listingData := storefront.GetListingData(listingPda)
			fmt.Println(listingPda, listingData)
		}

	*/
	listingPda := solana.MustPublicKeyFromBase58("EpP7m2njH5dmVyMQYBXYUjiWGYGq95aXLwDakLwFxJPm")
	listingData := storefront.GetListingData(listingPda)
	fmt.Println(listingPda, listingData)
}

