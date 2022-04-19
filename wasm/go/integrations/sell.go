package integrations

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	ag_binary "github.com/gagliardetto/binary"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	atok "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"strconv"
	"strings"
	"syscall/js"

	"creaturez.nft/wasm/v2/someplace"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func Sellables(this js.Value, args []js.Value) interface{} {
	holder := solana.MustPublicKeyFromBase58(args[0].String())
	oracle := solana.MustPublicKeyFromBase58(args[1].String())

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {

			nftsJson, err := nfts(holder, oracle)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}

			dst := js.Global().Get("Uint8Array").New(len(nftsJson))
			js.CopyBytesToJS(dst, nftsJson)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func SellCommit(this js.Value, args []js.Value) interface{} {
	holder := solana.MustPublicKeyFromBase58(args[0].String())
	oracle := solana.MustPublicKeyFromBase58(args[1].String())
	mint := solana.MustPublicKeyFromBase58(args[2].String())

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			nftsJson, err := burn(holder, mint, oracle)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}

			dst := js.Global().Get("Uint8Array").New(len(nftsJson))
			js.CopyBytesToJS(dst, nftsJson)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func MarketList(this js.Value, args []js.Value) interface{} {
	holder := solana.MustPublicKeyFromBase58(args[0].String())
	oracle := solana.MustPublicKeyFromBase58(args[1].String())
	marketUid := solana.MustPublicKeyFromBase58(args[2].String())
	nftMint := solana.MustPublicKeyFromBase58(args[3].String())
	price, priceError := strconv.Atoi(args[4].String())

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			if priceError != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}
			listTx, err := marketList(holder, oracle, marketUid, nftMint, uint64(price))
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New("unauthorized")
				reject.Invoke(errorObject)
				return
			}

			dst := js.Global().Get("Uint8Array").New(len(listTx))
			js.CopyBytesToJS(dst, listTx)

			resolve.Invoke(dst)
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func nfts(holder, oracle solana.PublicKey) ([]byte, error) {
	treasuryAuthority, _ := GetTreasuryAuthority(oracle)
	client := rpc.New(NETWORK)
	tokenAccounts, err := client.GetTokenAccountsByOwner(context.TODO(), holder, &rpc.GetTokenAccountsConfig{ProgramId: &solana.TokenProgramID}, &rpc.GetTokenAccountsOpts{Encoding: "jsonParsed"})
	if err != nil {
		return []byte{}, err
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
		Ata  string `json:"ata"`
		Mint string `json:"mint"`
		Uri  string `json:"uri"`
		Name string `json:"name"`
	}

	fmt.Println("asdfasdfasdfasdfasdfhgahjsdfghjagsdkfhjahjsdfglasdhjkfakljsdhfjkaghsdhjkfgahjksdfghjklasgdfhjk")
	tokens := func() []token {
		tokens := make([]token, 0)
		for _, tokenAccount := range tokenAccounts.Value {
			var tokenData tokenDataParsed
			err = json.Unmarshal(tokenAccount.Account.Data.GetRawJSON(), &tokenData)
			if err != nil {
				continue
			}
			if tokenData.Parsed.Info.TokenAmount.UiAmount == 1.0 {
				nftMint := solana.MustPublicKeyFromBase58(tokenData.Parsed.Info.Mint)
				metadata, _ := getMetadata(nftMint)
				metadataAccount, err := client.GetAccountInfo(context.TODO(), metadata)
				if err != nil {
					continue
				}
				var metadataData token_metadata.Metadata
				metadataDecoder := ag_binary.NewBorshDecoder(metadataAccount.Value.Data.GetBinary())
				err = metadataData.UnmarshalWithDecoder(metadataDecoder)
				if err != nil {
					continue
				}
				// derive treasuryWhitelist account from candy machine id in creator vector
				creators := *metadataData.Data.Creators
				candyMachine := creators[0].Address
				treasuryWhitelist, _ := GetTreasuryWhitelist(oracle, treasuryAuthority, candyMachine)
				fmt.Println("riptard?????", candyMachine, treasuryWhitelist)
				treasuryWhitelistAccount, err := client.GetAccountInfo(context.TODO(), treasuryWhitelist)
				if err != nil {
					continue
				}
				var treasuryWhitelistData someplace.TreasuryWhitelist
				treasuryWhitelistDecoder := ag_binary.NewBorshDecoder(treasuryWhitelistAccount.Value.Data.GetBinary())
				err = treasuryWhitelistData.UnmarshalWithDecoder(treasuryWhitelistDecoder)
				if err != nil {
					continue
				}
				fmt.Println("im a riptard", metadataData.Data.Name, candyMachine, treasuryWhitelistData.CandyMachineCreator)
				if !candyMachine.Equals(treasuryWhitelistData.CandyMachineCreator) {
					continue
				}
				fmt.Println("total riptard")
				var t token
				t.Uri = strings.Trim(metadataData.Data.Uri, "\u0000")
				t.Name = strings.Trim(metadataData.Data.Name, "\u0000")
				t.Ata = tokenAccount.Pubkey.String()
				t.Mint = tokenData.Parsed.Info.Mint
				tokens = append(tokens, t)
			}
			continue
		}

		return tokens
	}()

	nftsJson, err := json.Marshal(tokens)
	if err != nil {
		return []byte{}, err
	}
	return nftsJson, nil
}

func burn(holder, mint, oracle solana.PublicKey) ([]byte, error) {
	treasuryAuthority, treasuryAuthorityBump := GetTreasuryAuthority(oracle)
	treasuryMint := GetTreasuryAuthorityData(treasuryAuthority).TreasuryMint
	treasuryTokenAccount, _ := GetTreasuryTokenAccount(oracle)
	fmt.Println(mint)
	client := rpc.New(NETWORK)
	tokenAccounts, err := client.GetTokenAccountsByOwner(context.TODO(), holder, &rpc.GetTokenAccountsConfig{ProgramId: &solana.TokenProgramID}, &rpc.GetTokenAccountsOpts{Encoding: "jsonParsed"})
	if err != nil {
		fmt.Println("bad0")
		return []byte{}, errors.New("bad")
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
		ata               solana.PublicKey
		mint              solana.PublicKey
		metadata          solana.PublicKey
		treasuryWhitelist solana.PublicKey
	}

	tokens := func() []token {
		tokens := make([]token, 2)
		for _, tokenAccount := range tokenAccounts.Value {
			var tokenData tokenDataParsed
			err = json.Unmarshal(tokenAccount.Account.Data.GetRawJSON(), &tokenData)
			if err != nil {
				return tokens
			}
			fmt.Println("???", tokenData.Parsed.Info.Mint, mint)
			switch tokenData.Parsed.Info.Mint {
			case treasuryMint.String():
				{
					tokens[0] = token{
						ata:               tokenAccount.Pubkey,
						mint:              solana.MustPublicKeyFromBase58(tokenData.Parsed.Info.Mint),
						metadata:          solana.SystemProgramID,
						treasuryWhitelist: solana.SystemProgramID,
					}
					continue
				}
			case mint.String():
				{
					if tokenData.Parsed.Info.TokenAmount.UiAmount > 0 {
						nftMint := solana.MustPublicKeyFromBase58(tokenData.Parsed.Info.Mint)
						metadata, _ := getMetadata(nftMint)
						metadataAccount, err := client.GetAccountInfo(context.TODO(), metadata)
						if err != nil {
							continue
						}
						var metadataData token_metadata.Metadata
						metadataDecoder := ag_binary.NewBorshDecoder(metadataAccount.Value.Data.GetBinary())
						err = metadataData.UnmarshalWithDecoder(metadataDecoder)
						if err != nil {
							continue
						}
						creators := *metadataData.Data.Creators
						candyMachine := creators[0].Address
						fmt.Println("candy machine", candyMachine)
						treasuryWhitelist, _ := GetTreasuryWhitelist(oracle, treasuryAuthority, candyMachine)
						tokens[1] = token{
							ata:               tokenAccount.Pubkey,
							mint:              nftMint,
							metadata:          metadata,
							treasuryWhitelist: treasuryWhitelist,
						}
					}
					continue
				}
			}
		}

		return tokens
	}()
	fmt.Println(tokens)
	if !tokens[0].mint.Equals(treasuryMint) {
		fmt.Println("bad1")
		return []byte{}, errors.New("bad")
	}

	sellIxBuilder := someplace.NewSellForInstructionBuilder().
		SetDepoMintAccount(tokens[1].mint).
		SetDepoTokenAccountAccount(tokens[1].ata).
		SetInitializerAccount(holder).
		SetInitializerTokenAccountAccount(tokens[0].ata).
		SetMetadataAccount(tokens[1].metadata).
		SetTokenProgramAccount(solana.TokenProgramID).
		SetTreasuryAuthorityAccount(treasuryAuthority).
		SetTreasuryBump(treasuryAuthorityBump).
		SetTreasuryTokenAccountAccount(treasuryTokenAccount).
		SetTreasuryTokenMintAccount(treasuryMint).
		SetTreasuryWhitelistAccount(tokens[1].treasuryWhitelist).
		SetOracleAccount(oracle)

	sellIx := sellIxBuilder.Build()

	if e := sellIxBuilder.Validate(); e != nil {
		fmt.Println(e.Error())
		return []byte{}, errors.New("bad")
	}

	tx := solana.NewTransactionBuilder().AddInstruction(sellIx)
	txB, _ := tx.Build()
	txJson, _ := json.MarshalIndent(txB, "", "  ")
	fmt.Println(string(txJson))

	return txJson, nil
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

func marketList(holder, oracle, marketUid, nftMint solana.PublicKey, price uint64) ([]byte, error) {
	marketAuthority, _ := GetMarketAuthority(oracle, marketUid)
	marketAuthorityData := GetMarketAuthorityData(marketAuthority)
	marketMint := marketAuthorityData.MarketMint

	sellerTokenAccountAddress, _ := getTokenWallet(holder, nftMint)
	sellerMarketTokenAccountAddress, _ := getTokenWallet(holder, marketMint)

	marketListing, _ := GetMarketListing(marketAuthority, marketAuthorityData.Listings)
	marketListingTokenAccount, _ := GetMarketListingTokenAccount(marketAuthority, marketAuthorityData.Listings)
	listIx := someplace.NewCreateMarketListingInstructionBuilder().
		SetIndex(marketAuthorityData.Listings).
		SetMarketAuthorityAccount(marketAuthority).
		SetMarketListingAccount(marketListing).
		SetMarketListingTokenAccountAccount(marketListingTokenAccount).
		SetNftMintAccount(nftMint).
		SetPrice(price).
		SetRentAccount(solana.SysVarRentPubkey).
		SetSellerAccount(holder).
		SetSellerMarketTokenAccountAccount(sellerMarketTokenAccountAddress).
		SetSellerNftTokenAccountAccount(sellerTokenAccountAddress).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTokenProgramAccount(solana.TokenProgramID)

	err := listIx.Validate()
	if err != nil {
		return []byte{}, errors.New("bad")
	}
	var instructions []solana.Instruction
	client := rpc.New(NETWORK)
	tokenAccounts, err := client.GetTokenAccountsByOwner(context.TODO(), holder, &rpc.GetTokenAccountsConfig{Mint: &marketAuthorityData.MarketMint}, &rpc.GetTokenAccountsOpts{Encoding: "jsonParsed"})
	if err != nil {
		return []byte{}, errors.New("bad")
	}
	if len(tokenAccounts.Value) == 0 {
		instructions = append(
			instructions,
			atok.NewCreateInstructionBuilder().
				SetPayer(holder).
				SetWallet(holder).
				SetMint(marketMint).
				Build(),
		)
	}
	instructions = append(
		instructions,
		listIx.
			Build(),
	)

	txBuilder := solana.NewTransactionBuilder()
	for _, ix := range instructions {
		txBuilder = txBuilder.AddInstruction(ix)
	}
	txB, _ := txBuilder.Build()
	txJson, _ := json.MarshalIndent(txB, "", "  ")

	fmt.Println(string(txJson))
	return txJson, nil

}
