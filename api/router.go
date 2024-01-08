package api

import (
	"genericAPI/internal/endpoints/http_endpoints"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(app *gin.Engine) {
	app.Use(gin.Logger())
	app.NoRoute(noRoute)
	app.NoMethod(noMethod)
	indexRoute := app.GET("/", http_endpoints.TestEndpoint)
	indexRoute.Use()
}

func noRoute(c *gin.Context) {
	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found!"})
}

func noMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"code": "METHOD_NOT_ALLOWED", "message": "Method not allowed!"})
}
