package storefront

import (
	"encoding/binary"
	"fmt"

	"creaturez.nft/someplace"
	"github.com/gagliardetto/solana-go"
)

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

func GetCandyMachineCreator(candyMachineAddress solana.PublicKey) (solana.PublicKey, uint8, error) {
	return solana.FindProgramAddress(
		[][]byte{
			[]byte("someplace"),
			candyMachineAddress.Bytes(),
		},
		someplace.ProgramID,
	)
}

func GetMintHash(oracle, listing solana.PublicKey, mints uint64) (solana.PublicKey, uint8, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, mints)
	return solana.FindProgramAddress(
		[][]byte{
			oracle.Bytes(),
			listing.Bytes(),
			buf,
		},
		someplace.ProgramID,
	)
}

func GetMint(oracle, listing solana.PublicKey, mints uint64) (solana.PublicKey, uint8, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, mints)
	return solana.FindProgramAddress(
		[][]byte{
			[]byte("mintyhash"),
			oracle.Bytes(),
			listing.Bytes(),
			buf,
		},
		someplace.ProgramID,
	)
}

func GetVias(oracle solana.PublicKey) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			oracle.Bytes(),
			[]byte("via"),
		},
		someplace.ProgramID,
	)
	return addr, bump
}

func GetVia(oracle solana.PublicKey, viaIndex uint64) (solana.PublicKey, uint8) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, viaIndex)
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			oracle.Bytes(),
			[]byte("via"),
			buf,
		},
		someplace.ProgramID,
	)
	return addr, bump
}

func GetViaMapping(oracle, rarityTokenMint solana.PublicKey) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			oracle.Bytes(),
			[]byte("via"),
			rarityTokenMint.Bytes(),
		},
		someplace.ProgramID,
	)
	return addr, bump
}
