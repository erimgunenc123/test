package client

import (
	"context"
	"log/slog"

	qdb "github.com/questdb/go-questdb-client/v3"
)

type QuestClient struct {
	Sender qdb.LineSender
}

func NewQuestClient(ctx context.Context, conf QuestConfig) *QuestClient {
	sender, err := qdb.LineSenderFromConf(ctx, conf.ToConnectionString())
	if err != nil {
		slog.Error("failed creating QuestDB sender with error: ", slog.Any("error", err))
		return nil
	}
	return &QuestClient{
		Sender: sender,
	}
}
