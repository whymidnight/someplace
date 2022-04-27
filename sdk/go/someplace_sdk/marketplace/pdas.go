package marketplace

import (
	"encoding/binary"

	"creaturez.nft/someplace"
	"github.com/gagliardetto/solana-go"
)

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
