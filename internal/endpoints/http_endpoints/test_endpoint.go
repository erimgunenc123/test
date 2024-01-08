package http_endpoints

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestEndpoint(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, world!",
	})

}
