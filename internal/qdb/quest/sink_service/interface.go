package sink_service

import (
	"context"

	qdb "github.com/questdb/go-questdb-client/v3"
)

type TableData interface {
	Insert(context.Context, qdb.LineSender, string) error
}
