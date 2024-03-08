package algo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRunningAlgosEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"running_algos": []string{}})
}
