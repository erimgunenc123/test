package http_endpoints

import (
	"encoding/json"
	"genericAPI/internal/common/http_utils"
	"genericAPI/internal/customErrors"
	"io"
)

var url = "https://api.btcturk.com/api/v2/server/exchangeinfo"

func GetExchangeInfo() (*ExchangeInfo, error) {
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

func SymbolTicker() {

}
