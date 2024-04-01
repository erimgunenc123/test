package arbitrage

import "math/rand"

type ArbitrageService struct {
}

func (as *ArbitrageService) GetCandlestickSnapshot(symbols []string) []CandlestickData {
	var data []CandlestickData
	for i := 0; i < 100; i++ {
		open := rand.Float64() * 100
		close_ := rand.Float64() * 100
		high := max(open, close_) + rand.Float64()*10
		low := min(open, close_) - rand.Float64()*10
		if low < 0 {
			low = min(open, close_)
		}

		data = append(data, CandlestickData{
			Timestamp: 1710631777000 + int64(i*25),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close_,
		})
	}
	return data
}
