package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go/programs/token"

	"creaturez.nft/questing"
	"creaturez.nft/questing/quests"
	"creaturez.nft/questing/quests/ops"
	quest_ops "creaturez.nft/questing/quests/ops"
	"creaturez.nft/someplace"
	"creaturez.nft/someplace/marketplace"
	"creaturez.nft/someplace/storefront"
	storefront_ops "creaturez.nft/someplace/storefront/ops"
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

var MINT = solana.MustPublicKeyFromBase58("6c5EBgbPnpdZgKhXW4uTtcYojXqVNnVQbS2cdCHo8Zmu")

func init() {
	questing.SetProgramID(solana.MustPublicKeyFromBase58("Cr4keTx8UQiQ5F9TzTGdQ5dkcMHjxhYSAaHkHhUSABCk"))
	someplace.SetProgramID(solana.MustPublicKeyFromBase58("GXFE4Ym1vxhbXLBx2RxqL5y1Ee3XyFUqDksD7tYjAi8z"))
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

	// enableVias()
	// enableViaForRarityToken()

	// CreateNTokenAccountsOfMint(MINT, 2)
	// enableQuestsAndCreateQuest()
	// CreateAndAmmendEntitlementQuest()
	// startQuest()
	startAndEndQuest()
	// ETZoY7cJfD8N7EVx5tShRYS1vxgv3F4Dkavjb52kGRyj
	// treasureVerify()

  hash := sha256.Sum256([]byte("account:QuestAccount"))
  encoded := base58.Encode(hash[:8])
  fmt.Println(string(encoded))

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

	signers := make([]solana.PrivateKey, 0)
	rewardsMints := []solana.PrivateKey{
		solana.NewWallet().PrivateKey,
		solana.NewWallet().PrivateKey,
		solana.NewWallet().PrivateKey,
	}
	signers = append(signers, rewardsMints...)
	RARE := "rare"
	COMMON := "common"
	BALLZTASTIC := "ballztastic"
	rewards := []questing.Reward{
		{
			MintAddress:  rewardsMints[0].PublicKey(),
			RngThreshold: 40,
			Amount:       2,
			Cardinality:  &RARE,
		},
		{
			MintAddress:  rewardsMints[1].PublicKey(),
			RngThreshold: 40,
			Amount:       4,
			Cardinality:  &COMMON,
		},
		{
			MintAddress:  rewardsMints[2].PublicKey(),
			RngThreshold: 20,
			Amount:       1,
			Cardinality:  &BALLZTASTIC,
		},
	}

	quest_ops.EnableQuests(oracle)

	ixs := make([]solana.Instruction, 0)
	questData := questing.Quest{
		Index:           0,
		Name:            "aaa",
		Duration:        100,
		Oracle:          oracle.PublicKey(),
		WlCandyMachines: []solana.PublicKey{oracle.PublicKey()},
		Entitlement:     nil,
		Rewards:         rewards,
		Tender:          nil,
		/*
			Tender: &questing.Tender{
				MintAddress: MINT,
				Amount:      5,
			},
		*/
	}
	ix, questIndex := quest_ops.CreateQuest(oracle.PublicKey(), questData)
	questData.Index = questIndex
	ixs = append(ixs, ix)

	rewardIxs := ops.AppendQuestRewards(oracle.PublicKey(), questData)
	ixs = append(ixs, rewardIxs...)

	utils.SendTx(
		"list",
		ixs,
		append(signers, oracle),
		oracle.PublicKey(),
	)

	{
		// enable vias
		storefront_ops.EnableVias(oracle)

		// add reward mints to vias
		viaIxs := storefront_ops.EnableViasForRarityTokens(oracle.PublicKey(), func() []someplace.ViaMint {
			viaMints := make([]someplace.ViaMint, 0)
			for _, reward := range rewards {
				viaMints = append(viaMints, someplace.ViaMint{
					MintAddress: reward.MintAddress,
					Rarity:      *reward.Cardinality,
				})
			}
			return viaMints
		}())

		utils.SendTx(
			"list",
			viaIxs,
			append(signers, oracle),
			oracle.PublicKey(),
		)
	}
	{
		questsPda, _ := quests.GetQuests(oracle.PublicKey())
		questsData := quests.GetQuestsData(questsPda)
		quest, _ := quests.GetQuest(oracle.PublicKey(), questsData.Quests-1)
		fmt.Println(quest, questsData.Quests, questsData.Quests-1)
		questData := quests.GetQuestData(quest)
		{
			questDataJson, _ := json.MarshalIndent(questData, "", "  ")
			fmt.Println(string(questDataJson))
		}

	}
}
func CreateAndAmmendEntitlementQuest() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}

	ix, _ := quest_ops.CreateQuest(oracle.PublicKey(), questing.Quest{
		Index:           0,
		Name:            "aaa",
		Duration:        100,
		Oracle:          oracle.PublicKey(),
		WlCandyMachines: []solana.PublicKey{oracle.PublicKey()},
		Entitlement:     nil,
		Rewards: []questing.Reward{
			{
				MintAddress: MINT,
				Amount:      1,
			},
		},
		Tender: &questing.Tender{
			MintAddress: MINT,
			Amount:      5,
		},
	})
	utils.SendTx(
		"list",
		append(make([]solana.Instruction, 0), ix),
		append(make([]solana.PrivateKey, 0), oracle),
		oracle.PublicKey(),
	)

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

func enrollQuestor() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	pixelBallzMint := solana.MustPublicKeyFromBase58("7zaCj11reNw4FMxY5UqR8mjNdatgB4vgdN17eKAwMGie")

	_, _ = oracle, pixelBallzMint
}

func startQuest() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
	}
	for range make([]int, 2) {
		var instructions []solana.Instruction

		pixelBallzMint := solana.NewWallet()
		pixelBallzTokenAddress, _ := utils.GetTokenWallet(oracle.PublicKey(), pixelBallzMint.PublicKey())
		// ballzMint := solana.MustPublicKeyFromBase58("6NGcNWFVksoeXf1xgAvKubQgS6rW5EZ2oVwqAa1eHySz")
		// ballzTokenAddress := solana.MustPublicKeyFromBase58("57hobyD843HjijKTLAbiKcfPCdBY3bdDPgvKR4ggoGaz")
		// pixelBallzMint := solana.MustPublicKeyFromBase58("7zaCj11reNw4FMxY5UqR8mjNdatgB4vgdN17eKAwMGie")
		// pixelBallTokenAddress := solana.MustPublicKeyFromBase58("DpXfu5sQpfGM2wSRPq1nUs4iKqVkjSwCehDrErhysZLP")
		{

			client := rpc.New(NETWORK)
			min, err := client.GetMinimumBalanceForRentExemption(context.TODO(), token.MINT_SIZE, rpc.CommitmentFinalized)
			if err != nil {
				panic(err)
			}

			instructions = append(instructions,
				system.NewCreateAccountInstructionBuilder().
					SetOwner(token.ProgramID).
					SetNewAccount(pixelBallzMint.PublicKey()).
					SetSpace(token.MINT_SIZE).
					SetFundingAccount(oracle.PublicKey()).
					SetLamports(min).
					Build(),

				token.NewInitializeMint2InstructionBuilder().
					SetMintAccount(pixelBallzMint.PublicKey()).
					SetDecimals(0).
					SetMintAuthority(oracle.PublicKey()).
					SetFreezeAuthority(oracle.PublicKey()).
					Build(),

				atok.NewCreateInstructionBuilder().
					SetPayer(oracle.PublicKey()).
					SetWallet(oracle.PublicKey()).
					SetMint(pixelBallzMint.PublicKey()).
					Build(),

				token.NewMintToInstructionBuilder().
					SetMintAccount(pixelBallzMint.PublicKey()).
					SetDestinationAccount(pixelBallzTokenAddress).
					SetAuthorityAccount(oracle.PublicKey()).
					SetAmount(1).
					Build(),
			)
		}
		utils.SendTx(
			"list",
			instructions,
			append(make([]solana.PrivateKey, 0), oracle, pixelBallzMint.PrivateKey),
			oracle.PublicKey(),
		)

		/*
			fmt.Println("sleeping")
			time.Sleep(15 * time.Second)
		*/

		questInstructions := make([]solana.Instruction, 0)

		questor, _ := quests.GetQuestorAccount(oracle.PublicKey())
		questorData := quests.GetQuestorData(questor)
		if questorData == nil {
			questInstructions = append(
				questInstructions,
				ops.EnrollQuestor(oracle.PublicKey()),
			)
		}

		questee, _ := quests.GetQuesteeAccount(pixelBallzMint.PublicKey())
		questeeData := quests.GetQuesteeData(questee)
		if questeeData == nil {
			questInstructions = append(
				questInstructions,
				ops.EnrollQuestee(oracle.PublicKey(), pixelBallzMint.PublicKey(), pixelBallzTokenAddress),
			)
		}

		questPda, _ := quests.GetQuest(oracle.PublicKey(), 3)
		questAccount, _ := quests.GetQuestAccount(questor, questee, questPda)
		questDeposit, _ := quests.GetQuestDepositTokenAccount(questee, questPda)

		questData := quests.GetQuestData(questPda)

		startQuestIx := questing.NewStartQuestInstructionBuilder().
			SetDepositTokenAccountAccount(questDeposit).
			SetInitializerAccount(oracle.PublicKey()).
			SetPixelballzMintAccount(pixelBallzMint.PublicKey()).
			SetPixelballzTokenAccountAccount(pixelBallzTokenAddress).
			SetQuestAccAccount(questAccount).
			SetQuestAccount(questPda).
			SetQuesteeAccount(questee).
			SetQuestorAccount(questor).
			SetRentAccount(solana.SysVarRentPubkey).
			SetSystemProgramAccount(solana.SystemProgramID).
			SetTokenProgramAccount(solana.TokenProgramID)

		if questData.Tender != nil && questData.TenderSplits != nil {
			fmt.Println("asdfasdfasdfa")
			tenderTokenAccount, _ := utils.GetTokenWallet(oracle.PublicKey(), questData.Tender.MintAddress)
			startQuestIx.Append(&solana.AccountMeta{tenderTokenAccount, true, false})
			for _, tenderSplit := range *questData.TenderSplits {
				startQuestIx.Append(&solana.AccountMeta{tenderSplit.TokenAddress, true, false})
			}
		}

		if err = startQuestIx.Validate(); err != nil {
			panic(err)
		} else {
			questInstructions = append(
				questInstructions,
				startQuestIx.Build(),
			)
		}

		utils.SendTx(
			"init cm",
			questInstructions,
			append(make([]solana.PrivateKey, 0), oracle),
			oracle.PublicKey(),
		)
	}
}

func startAndEndQuest() {
	oracle, err := solana.PrivateKeyFromSolanaKeygenFile("./oracle.key")
	if err != nil {
		panic(err)
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

	pixelBallzMint := solana.NewWallet()
	pixelBallzTokenAddress, _ := utils.GetTokenWallet(oracle.PublicKey(), pixelBallzMint.PublicKey())
	// ballzMint := solana.MustPublicKeyFromBase58("6NGcNWFVksoeXf1xgAvKubQgS6rW5EZ2oVwqAa1eHySz")
	// ballzTokenAddress := solana.MustPublicKeyFromBase58("57hobyD843HjijKTLAbiKcfPCdBY3bdDPgvKR4ggoGaz")
	// pixelBallzMint := solana.MustPublicKeyFromBase58("7zaCj11reNw4FMxY5UqR8mjNdatgB4vgdN17eKAwMGie")
	// pixelBallTokenAddress := solana.MustPublicKeyFromBase58("DpXfu5sQpfGM2wSRPq1nUs4iKqVkjSwCehDrErhysZLP")
	{
		var instructions []solana.Instruction
		{

			client := rpc.New(NETWORK)
			min, err := client.GetMinimumBalanceForRentExemption(context.TODO(), token.MINT_SIZE, rpc.CommitmentFinalized)
			if err != nil {
				panic(err)
			}

			instructions = append(instructions,
				system.NewCreateAccountInstructionBuilder().
					SetOwner(token.ProgramID).
					SetNewAccount(pixelBallzMint.PublicKey()).
					SetSpace(token.MINT_SIZE).
					SetFundingAccount(oracle.PublicKey()).
					SetLamports(min).
					Build(),

				token.NewInitializeMint2InstructionBuilder().
					SetMintAccount(pixelBallzMint.PublicKey()).
					SetDecimals(0).
					SetMintAuthority(oracle.PublicKey()).
					SetFreezeAuthority(oracle.PublicKey()).
					Build(),

				atok.NewCreateInstructionBuilder().
					SetPayer(oracle.PublicKey()).
					SetWallet(oracle.PublicKey()).
					SetMint(pixelBallzMint.PublicKey()).
					Build(),

				token.NewMintToInstructionBuilder().
					SetMintAccount(pixelBallzMint.PublicKey()).
					SetDestinationAccount(pixelBallzTokenAddress).
					SetAuthorityAccount(oracle.PublicKey()).
					SetAmount(1).
					Build(),
			)
		}
		utils.SendTx(
			"list",
			instructions,
			append(make([]solana.PrivateKey, 0), oracle, pixelBallzMint.PrivateKey),
			oracle.PublicKey(),
		)

		/*
		   fmt.Println("sleeping")
		   time.Sleep(15 * time.Second)
		*/

		questInstructions := make([]solana.Instruction, 0)

		questor, _ := quests.GetQuestorAccount(oracle.PublicKey())
		questorData := quests.GetQuestorData(questor)
		if questorData == nil {
			questInstructions = append(
				questInstructions,
				ops.EnrollQuestor(oracle.PublicKey()),
			)
		}

		questee, _ := quests.GetQuesteeAccount(pixelBallzMint.PublicKey())
		questeeData := quests.GetQuesteeData(questee)
		if questeeData == nil {
			questInstructions = append(
				questInstructions,
				ops.EnrollQuestee(oracle.PublicKey(), pixelBallzMint.PublicKey(), pixelBallzTokenAddress),
			)
		}

		questPda, _ := quests.GetQuest(oracle.PublicKey(), 0)
		questAccount, _ := quests.GetQuestAccount(questor, questee, questPda)
		questDeposit, _ := quests.GetQuestDepositTokenAccount(questee, questPda)

		questData := quests.GetQuestData(questPda)

		startQuestIx := questing.NewStartQuestInstructionBuilder().
			SetDepositTokenAccountAccount(questDeposit).
			SetInitializerAccount(oracle.PublicKey()).
			SetPixelballzMintAccount(pixelBallzMint.PublicKey()).
			SetPixelballzTokenAccountAccount(pixelBallzTokenAddress).
			SetQuestAccAccount(questAccount).
			SetQuestAccount(questPda).
			SetQuesteeAccount(questee).
			SetQuestorAccount(questor).
			SetRentAccount(solana.SysVarRentPubkey).
			SetSystemProgramAccount(solana.SystemProgramID).
			SetTokenProgramAccount(solana.TokenProgramID)

		if questData.Tender != nil && questData.TenderSplits != nil {
			tenderTokenAccount, _ := utils.GetTokenWallet(oracle.PublicKey(), questData.Tender.MintAddress)
			startQuestIx.Append(&solana.AccountMeta{tenderTokenAccount, true, false})
			for _, tenderSplit := range *questData.TenderSplits {
				startQuestIx.Append(&solana.AccountMeta{tenderSplit.TokenAddress, true, false})
			}
		}

		if err = startQuestIx.Validate(); err != nil {
			panic(err)
		} else {
			questInstructions = append(
				questInstructions,
				startQuestIx.Build(),
			)
		}

		utils.SendTx(
			"init cm",
			questInstructions,
			append(make([]solana.PrivateKey, 0), oracle),
			oracle.PublicKey(),
		)
	}
	fmt.Println("Sleeping...")
	time.Sleep(5 * time.Second)
	fmt.Println("Slept")
	{
		questInstructions := make([]solana.Instruction, 0)

		questor, _ := quests.GetQuestorAccount(oracle.PublicKey())

		questee, _ := quests.GetQuesteeAccount(pixelBallzMint.PublicKey())

		questPda, questPdaBump := quests.GetQuest(oracle.PublicKey(), 0)
		questAccount, _ := quests.GetQuestAccount(questor, questee, questPda)
		questDeposit, questDepositBump := quests.GetQuestDepositTokenAccount(questee, questPda)
		questQuesteeReceipt, _ := quests.GetQuestQuesteeReceiptAccount(questor, questee, questPda)

		endQuestIx := questing.NewEndQuestInstructionBuilder().
			SetAssociatedTokenProgramAccount(solana.SPLAssociatedTokenAccountProgramID).
			SetDepositTokenAccountAccount(questDeposit).
			SetDepositTokenAccountBump(questDepositBump).
			SetInitializerAccount(oracle.PublicKey()).
			SetOracleAccount(oracle.PublicKey()).
			SetPixelballzMintAccount(pixelBallzMint.PublicKey()).
			SetPixelballzTokenAccountAccount(pixelBallzTokenAddress).
			SetQuestAccAccount(questAccount).
			SetQuestAccount(questPda).
			SetQuesteeAccount(questee).
			SetQuestorAccount(questor).
			SetQuestQuesteeReceiptAccount(questQuesteeReceipt).
			SetRentAccount(solana.SysVarRentPubkey).
			SetQuestBump(questPdaBump).
			SetSlotHashesAccount(solana.MustPublicKeyFromBase58("SysvarS1otHashes111111111111111111111111111")).
			SetSystemProgramAccount(solana.SystemProgramID).
			SetTokenProgramAccount(solana.TokenProgramID)

		questData := quests.GetQuestData(questPda)
		rewardMints := make([]solana.AccountMeta, 0)
		rewardAtas := make([]solana.AccountMeta, 0)
		fmt.Println("---------", len(endQuestIx.AccountMetaSlice))
		for _, reward := range questData.Rewards {
			endQuestIx.Append(&solana.AccountMeta{PublicKey: reward.MintAddress, IsWritable: true, IsSigner: false})
		}
		for _, reward := range questData.Rewards {
			rewardAta, _ := utils.GetTokenWallet(oracle.PublicKey(), reward.MintAddress)
			endQuestIx.Append(&solana.AccountMeta{PublicKey: rewardAta, IsWritable: true, IsSigner: false})
		}

		if err = endQuestIx.Validate(); err != nil {
			panic(err)
		} else {
			fmt.Println("---------", len(endQuestIx.AccountMetaSlice), len(rewardMints), len(rewardAtas))
			questInstructions = append(
				questInstructions,
				endQuestIx.Build(),
			)
		}

		utils.SendTx(
			"init cm",
			questInstructions,
			append(make([]solana.PrivateKey, 0), oracle),
			oracle.PublicKey(),
		)
	}

	{
		// rng after quest end
		batches, _ := storefront.GetBatches(oracle.PublicKey())
		batchesData := storefront.GetBatchesData(batches)
		questPda, _ := quests.GetQuest(oracle.PublicKey(), 0)
		questor, _ := quests.GetQuestorAccount(oracle.PublicKey())
		questee, _ := quests.GetQuesteeAccount(pixelBallzMint.PublicKey())

		// get quest questee reward account data for the reward chosen of this questee
		questQuesteeReceipt, _ := quests.GetQuestQuesteeReceiptAccount(questor, questee, questPda)
		questQuesteeReceiptData := quests.GetQuestQuesteeReceiptAccountData(questQuesteeReceipt)
		questQuesteeReceiptDataJs, _ := json.MarshalIndent(questQuesteeReceiptData, "", "  ")

		viaMap, _ := storefront.GetViaMapping(oracle.PublicKey(), questQuesteeReceiptData.RewardMint)
		viaMapData := storefront.GetViaMappingData(viaMap)
		fmt.Println(viaMap, viaMapData)
		via, viaBump := storefront.GetVia(batchesData.Oracle, viaMapData.ViasIndex)
		viaData := storefront.GetViaData(via)
		viaDataJs, _ := json.MarshalIndent(viaData, "", "  ")
		fmt.Println(string(questQuesteeReceiptDataJs))
		fmt.Println(string(viaDataJs))

		rewardTicket, rewardTicketBump := storefront.GetRewardTicket(via, questPda, questee, oracle.PublicKey())

		rewardTokenAccount, _ := utils.GetTokenWallet(oracle.PublicKey(), viaMapData.TokenMint)

		rngRewardIndiceIx := someplace.NewRngNftAfterQuestInstructionBuilder().
			SetBatchesAccount(batches).
			SetInitializerAccount(oracle.PublicKey()).
			SetQuestAccount(questPda).
			SetQuesteeAccount(questee).
			SetRewardTicketAccount(rewardTicket).
			SetRewardTokenAccountAccount(rewardTokenAccount).
			SetSlotHashesAccount(solana.MustPublicKeyFromBase58("SysvarS1otHashes111111111111111111111111111")).
			SetSystemProgramAccount(solana.SystemProgramID).
			SetViaAccount(via).
			SetViaBump(viaBump).
			SetViaMapAccount(viaMap)

		if err = rngRewardIndiceIx.Validate(); err != nil {
			panic(err)
		}

		for i := range make([]int, batchesData.Counter) {
			batchReceipt, _ := storefront.GetBatchReceipt(oracle.PublicKey(), uint64(i))
			batchReceiptData := storefront.GetBatchReceiptData(batchReceipt)
			batchAccount := batchReceiptData.BatchAccount
			batchCardinalitiesReport, _ := storefront.GetBatchCardinalitiesReport(batchAccount)

			rngRewardIndiceIx.Append(&solana.AccountMeta{
				PublicKey:  batchCardinalitiesReport,
				IsWritable: false,
				IsSigner:   false,
			})
		}

		utils.SendTx(
			"init cm",
			append(make([]solana.Instruction, 0), rngRewardIndiceIx.Build()),
			append(make([]solana.PrivateKey, 0), oracle),
			oracle.PublicKey(),
		)

		{
			rewardTicketData := storefront.GetRewardTicketData(rewardTicket)
			if rewardTicketData == nil {
				panic("null reward ticket data")
			}
			candyMachineAddress := rewardTicketData.BatchAccount
			batchCardinalitiesReport, _ := storefront.GetBatchCardinalitiesReport(candyMachineAddress)

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

			      mint, _, _ := storefront.GetMint(cm.Oracle, via, viaData.Mints)
			mintAta := solana.NewWallet()
			mintViaHash, _, _ := storefront.GetMintHashVia(cm.Oracle, viaData.TokenMint, viaData.Mints)
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

            mintIx := someplace.NewMintNftViaInstructionBuilder().
                SetClockAccount(solana.SysVarClockPubkey).
                SetCreatorBump(creatorBump).
                SetMintAccount(mint).
                SetViaAccount(via).
                SetMintAtaAccount(mintAta.PublicKey()).
                SetCandyMachineAccount(rewardTicketData.BatchAccount).
                SetMasterEditionAccount(masterEdition).
                SetCandyMachineCreatorAccount(candyMachineCreator).
                SetBatchCardinalitiesReportAccount(batchCardinalitiesReport).
                SetMetadataAccount(metadataAddress).
                SetRewardTicketAccount(rewardTicket).
                SetSlotHashesAccount(solana.MustPublicKeyFromBase58("SysvarS1otHashes111111111111111111111111111")).
                SetInstructionSysvarAccountAccount(solana.SysVarInstructionsPubkey).
                SetMintHashAccount(mintViaHash).
                SetOracleAccount(cm.Oracle).
                SetPayerAccount(oracle.PublicKey()).
                SetRentAccount(solana.SysVarRentPubkey).
                SetRewardTicketBump(rewardTicketBump).
                SetSystemProgramAccount(solana.SystemProgramID).
                SetRewardTokenAccountAccount(rewardTokenAccount).
                SetRewardTokenMintAccountAccount(viaMapData.TokenMint).
                SetTokenProgramAccount(solana.TokenProgramID).
                SetTokenMetadataProgramAccount(solana.TokenMetadataProgramID)


			err = mintIx.Validate()
			if err != nil {
				panic(err)
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
	}
}
