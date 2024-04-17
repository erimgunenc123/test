package exchange_info

import (
	"fmt"
	btcturk "genericAPI/btcturk_connector/http_endpoints"
)

// todo might need locks, assuming this whole service will initialize once and never write anything again

type btcTurkExchangeInfo struct {
	symbols map[string]btcturk.Symbol // symbol -> symbolInfo
}

var BtcTurkExchangeInfo *btcTurkExchangeInfo

func InitBtcTurkExchangeInfo() {
	BtcTurkExchangeInfo = &btcTurkExchangeInfo{
		symbols: make(map[string]btcturk.Symbol),
	}

	excInfo, err := btcturk.GetExchangeInfo()
	if err != nil {
		panic(fmt.Sprintf("Failed fetching exchange info with error:%s", err.Error()))
	}
	for _, symbol := range excInfo.Data.Symbols {
		BtcTurkExchangeInfo.symbols[symbol.Name] = symbol
	}
}

func (exc *btcTurkExchangeInfo) GetSymbols() []string {
	var res []string
	for k, _ := range exc.symbols {
		res = append(res, k)
	}
	return res
}
