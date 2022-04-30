package integrations

import (
	"context"
	"creaturez.nft/someplace"
	ag_binary "github.com/gagliardetto/binary"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type TransactionInstructionKey struct {
	Pubkey     string `json:"pubkey"`
	IsSigner   bool   `json:"isSigner"`
	IsWritable bool   `json:"isWritable"`
}
type TransactionInstruction struct {
	Keys      []TransactionInstructionKey `json:"keys"`
	ProgramId string                      `json:"programId"`
	Data      []int                       `json:"data"`
}

type BuyResponse struct {
	Tx      solana.Transaction `json:"transaction"`
	MintKey string             `json:"mintKey"`
}

func sendTx(
	doc string,
	instructions []solana.Instruction,
	signers []solana.PrivateKey,
	feePayer solana.PublicKey,
) (string, error) {
	rpcClient := rpc.New("https://psytrbhymqlkfrhudd.dev.genesysgo.net:8899/")
	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return "", err
	}

	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(feePayer),
	)
	if err != nil {
		return "", err
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
		return "", err
	}

	encoded, err := tx.ToBase64()
	if err != nil {
		return "", err
	}
	signature, err := rpcClient.SendEncodedTransaction(context.TODO(), encoded)
	if err != nil {
		return "", err
	}
	return signature.String(), nil
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
func GetTreasuryAuthorityData(treasuryAuthority solana.PublicKey) *someplace.TreasuryAuthority {
	rpcClient := rpc.New(NETWORK)
	batchesBin, _ := rpcClient.GetAccountInfoWithOpts(context.TODO(), treasuryAuthority, &rpc.GetAccountInfoOpts{Commitment: "confirmed"})
	var batchesData someplace.TreasuryAuthority
	decoder := ag_binary.NewBorshDecoder(batchesBin.Value.Data.GetBinary())
	err := batchesData.UnmarshalWithDecoder(decoder)
	if err != nil {
		return nil
	}

	return &batchesData

}
