package middlewares

import (
	"genericAPI/internal/common/constants"
	"genericAPI/internal/utils/authentication_utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidateAccessTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken_, ok := c.Request.Header["Authorization"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token not found"})
			c.Abort()
			return
		}
		if len(authToken_) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization token format"})
			c.Abort()
			return
		}
		userId, err := authentication_utils.ValidateAccessToken(authToken_[0])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set(constants.ContextUserIdKey, *userId)
		c.Next()
	}
}
