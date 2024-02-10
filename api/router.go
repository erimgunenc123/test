package api

import (
	"genericAPI/api/environment"
	"genericAPI/api/middlewares"
	"genericAPI/internal/endpoints/http_endpoints/login"
	"genericAPI/internal/endpoints/http_endpoints/register"
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

	// ws
	wsGroup := app.Group("/ws")
	wsGroup.Use(middlewares.WebsocketMiddleware())
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
	return cors.New(corsConf)
}
