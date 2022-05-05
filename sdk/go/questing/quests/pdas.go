package quests

import (
	"encoding/binary"

	"creaturez.nft/questing"
	"github.com/gagliardetto/solana-go"
)

func GetQuests(
	oracle solana.PublicKey,
) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("oracle"),
			oracle.Bytes(),
		},
		questing.ProgramID,
	)
	return addr, bump
}

func GetQuest(
	oracle solana.PublicKey,
	index uint64,
) (solana.PublicKey, uint8) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, index)
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("oracle"),
			oracle.Bytes(),
			buf,
		},
		questing.ProgramID,
	)
	return addr, bump
}

func GetQuestEntitlementTokenAccount(
	oracle solana.PublicKey,
	index uint64,
) (solana.PublicKey, uint8) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, index)
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("oracle"),
			[]byte("entitlement"),
			oracle.Bytes(),
			buf,
		},
		questing.ProgramID,
	)
	return addr, bump
}
