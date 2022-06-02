package main

import (
	"syscall/js"

	"creaturez.nft/questing"
	"creaturez.nft/someplace"
	"creaturez.nft/wasm/v2/integrations"
	questing_bridge "creaturez.nft/wasm/v2/integrations/questing_bridge"
	"creaturez.nft/wasm/v2/integrations/storefront"
	"github.com/gagliardetto/solana-go"
)

func main() {
	questing.SetProgramID(solana.MustPublicKeyFromBase58("Cr4keTx8UQiQ5F9TzTGdQ5dkcMHjxhYSAaHkHhUSABCk"))
	someplace.SetProgramID(solana.MustPublicKeyFromBase58("GXFE4Ym1vxhbXLBx2RxqL5y1Ee3XyFUqDksD7tYjAi8z"))
	done := make(chan struct{})

	global := js.Global()

	fetchNftsFunc := js.FuncOf(storefront.FetchNfts)
	defer fetchNftsFunc.Release()
	global.Set("reportCatalog", fetchNftsFunc)

	FetchNftHashMapFunc := js.FuncOf(storefront.FetchNftHashMap)
	defer FetchNftHashMapFunc.Release()
	global.Set("reportHashMap", FetchNftHashMapFunc)

	fetchListingsFunc := js.FuncOf(storefront.GetListings)
	defer fetchListingsFunc.Release()
	global.Set("getListings", fetchListingsFunc)

	sellablesFunc := js.FuncOf(integrations.Sellables)
	defer sellablesFunc.Release()
	global.Set("sellables", sellablesFunc)

	sellCommitFunc := js.FuncOf(integrations.SellCommit)
	defer sellCommitFunc.Release()
	global.Set("sellCommit", sellCommitFunc)

	buyFunc := js.FuncOf(integrations.Buy)
	defer buyFunc.Release()
	global.Set("buy", buyFunc)

	marketMetaFunc := js.FuncOf(integrations.FetchMarketplaceMeta)
	defer marketMetaFunc.Release()
	global.Set("getMarketMeta", marketMetaFunc)

	marketListFunc := js.FuncOf(integrations.MarketList)
	defer marketListFunc.Release()
	global.Set("marketListNft", marketListFunc)

	//MarketListBuyables
	marketListBuyablesFunc := js.FuncOf(integrations.MarketListBuyables)
	defer marketListBuyablesFunc.Release()
	global.Set("marketListBuyables", marketListBuyablesFunc)

	//MarketBuy
	marketBuyFunc := js.FuncOf(integrations.MarketBuy)
	defer marketBuyFunc.Release()
	global.Set("marketBuy", marketBuyFunc)

	marketDelistFunc := js.FuncOf(integrations.MarketDelist)
	defer marketDelistFunc.Release()
	global.Set("marketDelist", marketDelistFunc)

	doRngs := js.FuncOf(questing_bridge.DoRNGs)
	defer doRngs.Release()
	global.Set("do_rngs", doRngs)

	mintRewards := js.FuncOf(questing_bridge.MintRewards)
	defer mintRewards.Release()
	global.Set("mint_rewards", mintRewards)

	getRewards := js.FuncOf(questing_bridge.GetRewards)
	defer getRewards.Release()
	global.Set("get_rewards", getRewards)

	<-done
}


