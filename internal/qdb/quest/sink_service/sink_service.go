package sink_service

import (
	"context"
	"fmt"
	"genericAPI/internal/qdb/quest/client"
	"log/slog"
	"time"
)

// used generic interface type throughout the service because there is no need to figure out the type during compile
// time, this service isn't that performance critical and also it won't become a bottleneck before db insertion does
type QuestSinkService struct {
	client               *client.QuestClient
	sink                 chan TableData
	tableName            string
	batchSize            int // for maximum allowed message count in memory, should trigger a flush now and then
	defaultFlushInterval time.Duration
}

func NewQuestSinkService(client *client.QuestClient, tableName string, batchSize int) *QuestSinkService {
	return &QuestSinkService{
		client:               client,
		sink:                 make(chan TableData, 10000), // arbitrary size, can increase it if 10k isn't enough
		tableName:            tableName,
		batchSize:            batchSize,
		defaultFlushInterval: 5 * time.Second, // can parameterize it if needed
	}
}

func (q *QuestSinkService) Start(ctx context.Context) error {
	currentBatchCount := 0
	ticker := time.NewTicker(q.defaultFlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if currentBatchCount > 0 {
				if err := q.client.Sender.Flush(ctx); err != nil {
					return fmt.Errorf("failed flush before exiting: %w", err)
				}
			}
			return ctx.Err()

		case <-ticker.C:
			if currentBatchCount > 0 {
				slog.Info("Flushing (by ticker)")
				if err := q.client.Sender.Flush(ctx); err != nil {
					slog.ErrorContext(ctx, "failed flushing on timer",
						slog.Any("error", err),
						slog.String("table", q.tableName))
				}
				currentBatchCount = 0
			}

		case data := <-q.sink:
			slog.Info("QDB: Received bids and asks")
			if err := data.Insert(ctx, q.client.Sender, q.tableName); err != nil {
				slog.ErrorContext(ctx, "failed inserting data",
					slog.Any("error", err),
					slog.String("table", q.tableName))
			}

			currentBatchCount++
			if currentBatchCount >= q.batchSize {
				slog.Info("Flushing")
				if err := q.client.Sender.Flush(ctx); err != nil {
					slog.ErrorContext(ctx, "failed flushing batch",
						slog.Any("error", err),
						slog.String("table", q.tableName))
				}
				currentBatchCount = 0
				slog.Info("Flushed")
			}
		}
	}
}

func (q *QuestSinkService) GetDataChan() chan TableData {
	return q.sink
}
