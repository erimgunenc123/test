package models

import (
	"context"
	"time"

	qdb "github.com/questdb/go-questdb-client/v3"
)

type Candle struct {
	Symbol              string    `json:"Symbol,omitempty"`
	Open                float64   `json:"Open,omitempty"`
	High                float64   `json:"High,omitempty"`
	Low                 float64   `json:"Low,omitempty"`
	Close               float64   `json:"Close,omitempty"`
	Volume              float64   `json:"Volume,omitempty"`
	OpenTime            time.Time `json:"OpenTime,omitempty"`
	CloseTime           time.Time `json:"CloseTime,omitempty"`
	TotalTradesInCandle uint32    `json:"TotalTradesInCandle,omitempty"`
	QuoteAssetVolume    float64   `json:"QuoteAssetVolume,omitempty"`
	BuyBaseAssetVolume  float64   `json:"BuyBaseAssetVolume,omitempty"`
	BuyQuoteAssetVolume float64   `json:"BuyQuoteAssetVolume,omitempty"`
}

func (c *Candle) Insert(ctx context.Context, sender qdb.LineSender, tableName string) error {
	return sender.Table(tableName).
		Symbol("symbol", c.Symbol).
		Float64Column("o", c.Open).
		Float64Column("h", c.High).
		Float64Column("l", c.Low).
		Float64Column("c", c.Close).
		Int64Column("volume", int64(c.Volume)). // realistically there shouldn't be a loss in these conversions
		Int64Column("quote_asset_vol", int64(c.QuoteAssetVolume)).
		Int64Column("buy_base_asset_vol", int64(c.BuyBaseAssetVolume)).
		Int64Column("buy_quote_asset_vol", int64(c.BuyQuoteAssetVolume)).
		Int64Column("num_trades", int64(c.TotalTradesInCandle)).
		At(ctx, c.OpenTime)
}
