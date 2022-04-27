package main

import (
	"errors"
	"io/ioutil"
	"os"

	"creaturez.nft/operatorCLI/v2/commands/storefront"
	"creaturez.nft/someplace"
	"github.com/gagliardetto/solana-go"
	"gopkg.in/yaml.v2"
)

type configYaml struct {
	OraclePath    string `yaml:"OraclePath"`
	TreasuryMint  string `yaml:"TreasuryMint"`
	ListingsTable string `yaml:"ListingsTable"`
	HashList      string `yaml:"HashList"`
}

func init() {
	someplace.SetProgramID(solana.MustPublicKeyFromBase58("8otw5mCMUtwx91e7q7MAyhWoQVnc3Ng72qwDH58z72VW"))
}

func readConfig(configPath string) *configYaml {
	var config = new(configYaml)
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(configFile, config); err != nil {
		panic(err)
	}

	return config
}

func main() {
	config := readConfig(os.Args[1])
	if config == nil {
		panic(errors.New("bad config"))
	}

	oracle, err := solana.PrivateKeyFromSolanaKeygenFile(config.OraclePath)
	if err != nil {
		panic(err)
	}

	switch os.Args[2] {
	case "sync_listings":
		{
			storefront.SyncListingsTable(oracle, config.ListingsTable)
			storefront.ReportCatalog(oracle.PublicKey(), config.ListingsTable)
			break
		}
	case "instance":
		{
			treasuryMint, err := solana.PublicKeyFromBase58(config.TreasuryMint)
			if err != nil {
				panic(errors.New("bad treasury mint in config"))
			}
			storefront.Instance(oracle, treasuryMint)
			break
		}
	case "report":
		{
			storefront.ReportCatalog(oracle.PublicKey(), config.ListingsTable)
			break
		}
	}

}
