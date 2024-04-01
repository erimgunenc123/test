package api

import (
	"genericAPI/api/environment"
	"genericAPI/api/middlewares"
	"genericAPI/internal/endpoints/http_endpoints/algo"
	"genericAPI/internal/endpoints/http_endpoints/login"
	http_marketdata "genericAPI/internal/endpoints/http_endpoints/marketdata"
	"genericAPI/internal/endpoints/http_endpoints/register"
	"genericAPI/internal/endpoints/websocket_endpoints/marketdata"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(app *gin.Engine) {
	app.Use(gin.Logger())
	app.Use(setCORS())
	app.NoRoute(noRoute)
	app.NoMethod(noMethod)
	addEndpointRouters(app)
}

func addEndpointRouters(app *gin.Engine) {
	// no need to handle logout
	// frontend will remove the stored jwt on logout and redirect to index page
	app.POST("/login", login.LoginEndpoint)
	app.POST("/register", register.RegisterEndpoint)
	authGroup := app.Group("/auth")
	authGroup.Use(middlewares.ValidateAccessTokenMiddleware())

	// arbitrage
	arbitrageGroup := authGroup.Group("/arbitrage")
	arbitrageGroup.POST("/candlestick", http_marketdata.ArbitrageCandlestickDataEndpoint)

	// algo
	algoGroup := authGroup.Group("/algos")
	algoGroup.GET("/get_running_algos", algo.GetRunningAlgosEndpoint)

	// ws
	wsGroup := authGroup.Group("/ws")
	wsGroup.Use(middlewares.WebsocketUpgradeMiddleware())
	wsGroup.GET("/marketdata", marketdata.MarketdataWsHandler)
}

func noRoute(c *gin.Context) {
	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found!"})
}

func noMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"code": "METHOD_NOT_ALLOWED", "message": "Method not allowed!"})
}

// todo configure later
func setCORS() gin.HandlerFunc {
	corsConf := cors.DefaultConfig()
	if environment.IsTestEnvironment() {
		corsConf.AllowAllOrigins = true
	}
	corsConf.AllowHeaders = append(corsConf.AllowHeaders, "Authorization")
	return cors.New(corsConf)
}
