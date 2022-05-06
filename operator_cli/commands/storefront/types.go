package storefront

type Catalog struct {
	BatchAccount   string  `json:"batchAccount"`
	ConfigIndex    int     `json:"configIndex"`
	IsListed       bool    `json:"isListed"`
	Price          float64 `json:"price"`
	LifecycleStart int     `json:"lifecycleStart"`
	Mints          int     `json:"mints"`
	Resync         bool    `json:"resync"`
}

