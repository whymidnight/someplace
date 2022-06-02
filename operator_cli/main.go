package main

import (
	"errors"
	"fmt"
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
	HashMap       string `yaml:"HashMap"`
	Splits        []struct {
		TokenAddress string `yaml:"TokenAddress"`
		OpCode       uint8  `yaml:"OpCode"`
		Share        uint8  `yaml:"Share"`
	} `yaml:"Splits"`
}

func init() {
	someplace.SetProgramID(solana.MustPublicKeyFromBase58("GXFE4Ym1vxhbXLBx2RxqL5y1Ee3XyFUqDksD7tYjAi8z"))
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
	config := readConfig(os.Args[2])
	if config == nil {
		panic(errors.New("bad config"))
	}

	oracle, err := solana.PrivateKeyFromSolanaKeygenFile(config.OraclePath)
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "sync_listings":
		{
			storefront.SyncListingsTable(oracle, config.ListingsTable)
			storefront.ReportCatalog(oracle.PublicKey(), config.ListingsTable)
			break
		}
	case "sync_storefront_splits":
		{
			// storefront.ReportCatalog(oracle.PublicKey(), config.ListingsTable)
			fmt.Println(config.Splits)
			splits := make([]someplace.Split, 0)
			for i := range make([]int, len(config.Splits)) {
				splits = append(splits, someplace.Split{TokenAddress: solana.MustPublicKeyFromBase58(config.Splits[i].TokenAddress), OpCode: config.Splits[i].OpCode, Share: config.Splits[i].Share})
				i++
			}
			storefront.AmmendStorefrontSplits(oracle, splits)
			break
		}
	case "instance":
		{
			treasuryMint, err := solana.PublicKeyFromBase58(config.TreasuryMint)
			if err != nil {
				panic(errors.New("bad treasury mint in config"))
			}
			splits := make([]someplace.Split, 0)
			for i := range make([]int, len(config.Splits)) {
				splits = append(splits, someplace.Split{TokenAddress: solana.MustPublicKeyFromBase58(config.Splits[i].TokenAddress), OpCode: config.Splits[i].OpCode, Share: config.Splits[i].Share})
				i++
			}
			storefront.Instance(oracle, treasuryMint, splits)
			break
		}
	case "report_listings":
		{
			storefront.ReportCatalog(oracle.PublicKey(), config.ListingsTable)
			break
		}
	case "report_hashmap":
		{
			storefront.ReportHashMap(oracle.PublicKey(), config.HashMap)
			break
		}
	case "report_cardinalities":
		{
			storefront.ReportCardinalities(oracle)
			break
		}
	case "report_via_mints":
		{
			storefront.ReportViaMintingHashMap(oracle.PublicKey())
			break
		}
	}

}
