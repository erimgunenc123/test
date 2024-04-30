package http_endpoints

import (
	"encoding/json"
	binance_constants2 "genericAPI/exchange/binanceconnector/binance_constants"
	"genericAPI/internal/common/http_utils"
	"genericAPI/internal/customErrors"
	"io"
)

// GetOrderbookSnapshot symbols example: USDTBTC
func GetOrderbookSnapshot(symbols string) (*OrderbookSnapshot, error) {
	url := binance_constants2.BaseHttpUrl + binance_constants2.Depth
	// todo change limit
	resp, err := http_utils.GetRequest(url, nil, map[string]string{"symbol": symbols, "limit": "100"})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		res, _ := io.ReadAll(resp.Body)
		var orderbookSnapshot OrderbookSnapshot
		err := json.Unmarshal(res, &orderbookSnapshot)
		if err != nil {
			return nil, err
		}
		return &orderbookSnapshot, nil
	} else {
		return nil, customErrors.ErrFailedRequest
	}

}

func GetExchangeInfo() (*ExchangeInfo, error) {
	url := binance_constants2.BaseHttpUrl + binance_constants2.ExchangeInfo
	resp, err := http_utils.GetRequest(url, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		res, _ := io.ReadAll(resp.Body)
		var excInfo ExchangeInfo
		err := json.Unmarshal(res, &excInfo)
		if err != nil {
			return nil, err
		}
		return &excInfo, nil
	} else {
		return nil, customErrors.ErrFailedRequest
	}

}
