# Questing Contract

## Architecture
`Quests` (one) -> `Quest` (many) -> `Enrollment` (many)

A `Quest` requires `Quests` to account the `Quest` index as `+1` in order to be able to derive all the quests created by the given `oracle`.

## Notes

Since quests _mint_ out a reward token, quests lack the treasury regime when rewarding participation by default. You can execute `ammend_quest_with_entitlement` instruction to associate a treasury token account for the reward mint. However since the questing program does not have minting authority over such token, the supply of such reward token needs to be managed and resupplied.

In order to achieve n% chance/probability, the process of achieving a random number is via the lastest blockhash. We slice 8 bytes off the blockhash and convert from little endian bytes to an unsigned integer which gets modulated over 100. The modulated number should be 0-99 which can represent the percentile of items from our "drop table".
