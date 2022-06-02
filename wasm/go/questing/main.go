package main

import (
	"syscall/js"

	questing_program "creaturez.nft/questing"
	"creaturez.nft/wasm/v2/integrations/questing"
	"github.com/gagliardetto/solana-go"
)

func main() {
	global := js.Global()
	done := make(chan struct{})
	questing_program.SetProgramID(solana.MustPublicKeyFromBase58("Cr4keTx8UQiQ5F9TzTGdQ5dkcMHjxhYSAaHkHhUSABCk"))

	getQuests := js.FuncOf(questing.GetQuests)
	defer getQuests.Release()
	global.Set("get_quests", getQuests)

	enrollQuestor := js.FuncOf(questing.EnrollQuestor)
	defer enrollQuestor.Release()
	global.Set("enroll_questor", enrollQuestor)

	enrollQuestees := js.FuncOf(questing.EnrollQuestees)
	defer enrollQuestees.Release()
	global.Set("enroll_questees", enrollQuestees)

	startQuests := js.FuncOf(questing.StartQuests)
	defer startQuests.Release()
	global.Set("start_quests", startQuests)

	getQuested := js.FuncOf(questing.GetQuested)
	defer getQuested.Release()
	global.Set("get_quested", getQuested)

	endQuests := js.FuncOf(questing.EndQuests)
	defer endQuests.Release()
	global.Set("end_quests", endQuests)

	<-done
}
