package main

import (
	"context"
	"genericAPI/api"
	"genericAPI/api/api_config"
	"genericAPI/api/environment"
	"genericAPI/exchange/binanceconnector/connection_manager"
	"genericAPI/exchange/btcturk_connector/tickers"
	"genericAPI/internal/qdb/quest/client"
	"genericAPI/internal/qdb/quest/sink_service"
	"genericAPI/internal/services/marketdata/exchange_info"
	"genericAPI/internal/services/marketdata/orderbook"
	"log/slog"
	_ "net/http/pprof"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	environment.ParseArgs()
	api_config.InitConfig()
	ctx := context.Background()
	//database_connection.InitDB(database_logger.InitDbLogger())
	app := gin.Default()
	api.ConfigureGin(app)
	api.InitRouter(app)
	//dbops.Migrate() // disabled on prod env
	connection_manager.InitBinanceConnectionManager()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		exchange_info.InitBtcTurkExchangeInfo()
		exchange_info.InitBinanceExchangeInfo()
		questConfig := client.QuestConfig{
			Host:     "localhost",
			Port:     9000,
			Username: "admin",
			Password: "quest",
		}
		questClient := client.NewQuestClient(ctx, questConfig)
		if questClient == nil {
			slog.Error("failed to create QuestDB client")
			return
		}

		sinkService := sink_service.NewQuestSinkService(questClient, "orderbook_deltas", 5)

		go sinkService.Start(ctx)
		orderbook.InitOrderbookService(sinkService)
	}()
	go func() {
		defer wg.Done()
		tickers.InitTickerService()
	}()
	wg.Wait()

	panic(app.Run(":" + api_config.Config.App.Port))
}
