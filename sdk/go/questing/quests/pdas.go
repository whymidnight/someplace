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

func GetQuestDepositTokenAccount(
	questee solana.PublicKey,
	quest solana.PublicKey,
) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("questing"),
			questee.Bytes(),
			quest.Bytes(),
		},
		questing.ProgramID,
	)
	return addr, bump
}

func GetQuestAccount(
	questor solana.PublicKey,
	questee solana.PublicKey,
	quest solana.PublicKey,
) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("questing"),
			questor.Bytes(),
			questee.Bytes(),
			quest.Bytes(),
		},
		questing.ProgramID,
	)
	return addr, bump
}

func GetQuestorAccount(
	oracle solana.PublicKey,
) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("questing"),
			oracle.Bytes(),
		},
		questing.ProgramID,
	)
	return addr, bump
}

func GetQuesteeAccount(
	pixelBallzMint solana.PublicKey,
) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("questing"),
			pixelBallzMint.Bytes(),
		},
		questing.ProgramID,
	)
	return addr, bump
}

func GetQuestQuesteeReceiptAccount(
	questor solana.PublicKey,
	questee solana.PublicKey,
	quest solana.PublicKey,
) (solana.PublicKey, uint8) {
	addr, bump, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("quest_reward"),
			questor.Bytes(),
			questee.Bytes(),
			quest.Bytes(),
		},
		questing.ProgramID,
	)
	return addr, bump
}
