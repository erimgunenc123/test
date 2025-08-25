package models

import (
	"context"
	"encoding/json"
	"time"

	qdb "github.com/questdb/go-questdb-client/v3"
)

type OrderbookDelta struct {
	At            time.Time
	Symbol        string
	Bids          [][]string
	Asks          [][]string
	EventType     string
	FirstUpdateId int64
	FinalUpdateId int64
}

func (o *OrderbookDelta) Insert(ctx context.Context, sender qdb.LineSender, tableName string) error {
	// todo use Float64Array2DColumn once this merges: https://github.com/questdb/go-questdb-client/pull/55
	bidsJSON, _ := json.Marshal(o.Bids)
	asksJSON, _ := json.Marshal(o.Asks)
	return sender.Table(tableName).
		Symbol("symbol", o.Symbol).
		Symbol("event_type", o.EventType).
		StringColumn("bids", string(bidsJSON)).
		StringColumn("asks", string(asksJSON)).
		Int64Column("first_update_id", o.FirstUpdateId).
		Int64Column("final_update_id", o.FinalUpdateId).
		At(ctx, o.At)
}
