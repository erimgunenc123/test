package logout

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LogoutEndpoint(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out."})
}
