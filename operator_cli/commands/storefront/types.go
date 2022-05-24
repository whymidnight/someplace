package storefront

type Catalog struct {
	BatchAccount   string `json:"batchAccount"`
	ConfigIndex    int    `json:"configIndex"`
	IsListed       bool   `json:"isListed"`
	Price          int    `json:"price"`
	LifecycleStart int    `json:"lifecycleStart"`
	Mints          int    `json:"mints"`
	Resync         bool   `json:"resync"`
}

type ConfigLine struct {
	Name        string `json:"name"`
	Cardinality string `json:"cardinality"`
	URI         string `json:"uri"`
}
