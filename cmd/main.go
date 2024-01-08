package main

import (
	"genericAPI/api"
	"genericAPI/api/api_config"
	"genericAPI/api/environment"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	environment.ParseArgs()
	api_config.InitConfig()
	api.InitDB()
	api.ConfigureGinLogger()
	app := gin.Default()
	api.InitRouter(app)
	log.Fatal(app.Run(":" + api_config.Config.App.Port))
}
