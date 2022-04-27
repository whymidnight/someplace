package ops

import (
	"creaturez.nft/someplace/marketplace"
	"github.com/gagliardetto/solana-go"
)

func VerifyMarketList(oracle, marketUid solana.PublicKey) {
	marketAuthority, _ := marketplace.GetMarketAuthority(oracle, marketUid)
	marketAuthorityData := marketplace.GetMarketAuthorityData(marketAuthority)
	i := uint64(0)
	for i < marketAuthorityData.Listings {
		marketListing, _ := marketplace.GetMarketListing(marketAuthority, i)
		marketplace.GetMarketListingData(marketListing)
		i++
	}

}

