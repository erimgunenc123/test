package exchange_info

import (
	"fmt"
	"genericAPI/binanceconnector/http_endpoints"
)

// todo might need locks, assuming this whole service will initialize once and never write anything again

type binanceExchangeInfo struct {
	symbols    map[string]http_endpoints.SymbolInfo // symbol -> symbolInfo
	rateLimits map[string]http_endpoints.RateLimit  // rateLimitType -> rateLimit
}

var BinanceExchangeInfo *binanceExchangeInfo

func InitBinanceExchangeInfo() {
	BinanceExchangeInfo = &binanceExchangeInfo{
		symbols:    make(map[string]http_endpoints.SymbolInfo),
		rateLimits: make(map[string]http_endpoints.RateLimit),
	}

	excInfo, err := http_endpoints.GetExchangeInfo()
	if err != nil {
		panic(fmt.Sprintf("Failed fetching exchange info with error:%s", err.Error()))
	}
	for _, symbol := range excInfo.Symbols {
		BinanceExchangeInfo.symbols[symbol.Symbol] = symbol
	}
	for _, rateLimit := range excInfo.RateLimits {
		BinanceExchangeInfo.rateLimits[rateLimit.RateLimitType] = rateLimit
	}
}

func (exc *binanceExchangeInfo) GetSymbols() []string {
	var res []string
	for k, _ := range exc.symbols {
		res = append(res, k)
	}
	return res
}
