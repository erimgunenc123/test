package http_endpoints

type Symbol struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	NameNormalized   string `json:"nameNormalized"`
	Status           string `json:"status"`
	Numerator        string `json:"numerator"`
	Denominator      string `json:"denominator"`
	NumeratorScale   int    `json:"numeratorScale"`
	DenominatorScale int    `json:"denominatorScale"`
	HasFraction      bool   `json:"hasFraction"`
	Filters          []struct {
		FilterType       string      `json:"filterType"`
		MinPrice         string      `json:"minPrice"`
		MaxPrice         string      `json:"maxPrice"`
		TickSize         string      `json:"tickSize"`
		MinExchangeValue string      `json:"minExchangeValue"`
		MinAmount        interface{} `json:"minAmount"`
		MaxAmount        interface{} `json:"maxAmount"`
	} `json:"filters"`
	OrderMethods                          []string `json:"orderMethods"`
	DisplayFormat                         string   `json:"displayFormat"`
	CommissionFromNumerator               bool     `json:"commissionFromNumerator"`
	Order                                 int      `json:"order"`
	PriceRounding                         bool     `json:"priceRounding"`
	IsNew                                 bool     `json:"isNew"`
	MarketPriceWarningThresholdPercentage float64  `json:"marketPriceWarningThresholdPercentage"`
	MaximumOrderAmount                    *float64 `json:"maximumOrderAmount"`
}

type ExchangeInfo struct {
	Data struct {
		TimeZone   string   `json:"timeZone"`
		ServerTime int64    `json:"serverTime"`
		Symbols    []Symbol `json:"symbols"`
		Currencies []struct {
			Id            int     `json:"id"`
			Symbol        string  `json:"symbol"`
			MinWithdrawal float64 `json:"minWithdrawal"`
			MinDeposit    float64 `json:"minDeposit"`
			Precision     int     `json:"precision"`
			Address       struct {
				MinLen *int `json:"minLen"`
				MaxLen *int `json:"maxLen"`
			} `json:"address"`
			CurrencyType string `json:"currencyType"`
			Tag          struct {
				Enable bool    `json:"enable"`
				Name   *string `json:"name"`
				MinLen *int    `json:"minLen"`
				MaxLen *int    `json:"maxLen"`
			} `json:"tag"`
			Color                      string `json:"color"`
			Name                       string `json:"name"`
			IsAddressRenewable         bool   `json:"isAddressRenewable"`
			GetAutoAddressDisabled     bool   `json:"getAutoAddressDisabled"`
			IsPartialWithdrawalEnabled bool   `json:"isPartialWithdrawalEnabled"`
			IsNew                      bool   `json:"isNew"`
		} `json:"currencies"`
		CurrencyOperationBlocks []struct {
			CurrencySymbol     string `json:"currencySymbol"`
			WithdrawalDisabled bool   `json:"withdrawalDisabled"`
			DepositDisabled    bool   `json:"depositDisabled"`
		} `json:"currencyOperationBlocks"`
	} `json:"data"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type SymbolTick struct {
	Data []struct {
		Pair              string  `json:"pair"`
		PairNormalized    string  `json:"pairNormalized"`
		Timestamp         int64   `json:"timestamp"`
		Last              int     `json:"last"`
		High              int     `json:"high"`
		Low               int     `json:"low"`
		Bid               int     `json:"bid"`
		Ask               int     `json:"ask"`
		Open              int     `json:"open"`
		Volume            float64 `json:"volume"`
		Average           int     `json:"average"`
		Daily             int     `json:"daily"`
		DailyPercent      float64 `json:"dailyPercent"`
		DenominatorSymbol string  `json:"denominatorSymbol"`
		NumeratorSymbol   string  `json:"numeratorSymbol"`
		Order             int     `json:"order"`
	} `json:"data"`
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
	Code    int         `json:"code"`
}
