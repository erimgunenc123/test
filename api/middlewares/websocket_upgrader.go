package middlewares

import (
	"genericAPI/internal/common/constants"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func WebsocketMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		wsUpgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set(constants.ContextWebsocketConnectionKey, conn)
		c.Next()
	}
}
