package models

import (
	"context"
	"encoding/json"
	"time"

	qdb "github.com/questdb/go-questdb-client/v3"
)

type OrderbookSnapshots struct {
	At           time.Time
	Symbol       string
	Bids         [][2]string
	Asks         [][2]string
	LastUpdateId int64
}

func (o *OrderbookSnapshots) Insert(ctx context.Context, sender qdb.LineSender, tableName string) error {
	// todo use Float64Array2DColumn once this merges: https://github.com/questdb/go-questdb-client/pull/55
	bidsJSON, _ := json.Marshal(o.Bids)
	asksJSON, _ := json.Marshal(o.Asks)
	return sender.Table(tableName).
		Symbol("symbol", o.Symbol).
		StringColumn("bids", string(bidsJSON)).
		StringColumn("asks", string(asksJSON)).
		Int64Column("last_update_id", o.LastUpdateId).
		At(ctx, o.At)
}
