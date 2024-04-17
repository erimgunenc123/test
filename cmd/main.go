package main

import (
	"genericAPI/api"
	"genericAPI/api/api_config"
	"genericAPI/api/database_connection"
	"genericAPI/api/database_logger"
	"genericAPI/api/environment"
	"genericAPI/binanceconnector/connection_manager"
	"genericAPI/btcturk_connector/tickers"
	"genericAPI/internal/dbops"
	"genericAPI/internal/services/marketdata/exchange_info"
	"genericAPI/internal/services/marketdata/orderbook"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

func main() {
	environment.ParseArgs()
	api_config.InitConfig()

	database_connection.InitDB(database_logger.InitDbLogger())
	app := gin.Default()
	api.ConfigureGin(app)
	api.InitRouter(app)
	dbops.Migrate() // disabled on prod env
	connection_manager.InitBinanceConnectionManager()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		exchange_info.InitBtcTurkExchangeInfo()
		exchange_info.InitBinanceExchangeInfo()
		orderbook.InitOrderbookService()
	}()
	go func() {
		defer wg.Done()
		tickers.InitTickerService()
	}()
	wg.Wait()

	log.Fatal(app.Run(":" + api_config.Config.App.Port))
}
