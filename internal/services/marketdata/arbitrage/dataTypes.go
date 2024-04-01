package arbitrage

type CandlestickData struct {
	Timestamp int64   `json:"x"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
}
