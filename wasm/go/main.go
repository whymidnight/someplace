package main

import (
	"syscall/js"

	"creaturez.nft/wasm/v2/integrations"
)

func main() {
	done := make(chan struct{})

	global := js.Global()

	fetchNftsFunc := js.FuncOf(integrations.FetchNfts)
	defer fetchNftsFunc.Release()
	global.Set("reportCatalog", fetchNftsFunc)

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

	<-done
}
