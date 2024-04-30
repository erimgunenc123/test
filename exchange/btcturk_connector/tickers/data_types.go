package tickers

type Tick struct {
	Data    []TickData  `json:"data"`
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
	Code    int         `json:"code"`
}

type TickData struct {
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
	Average           float64 `json:"average"`
	Daily             int     `json:"daily"`
	DailyPercent      float64 `json:"dailyPercent"`
	DenominatorSymbol string  `json:"denominatorSymbol"`
	NumeratorSymbol   string  `json:"numeratorSymbol"`
}
