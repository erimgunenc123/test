package marketdata

import (
	"encoding/json"
	"genericAPI/internal/services/marketdata/arbitrage"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type arbitrageRequestBody struct {
	Symbols []string `json:"symbols"`
}

func ArbitrageCandlestickDataEndpoint(ctx *gin.Context) {
	body, _ := io.ReadAll(ctx.Request.Body)
	var reqBody arbitrageRequestBody
	err := json.Unmarshal(body, &reqBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON body!"})
		return
	}
	if reqBody.Symbols == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing field: symbols"})
		return
	}
	service := arbitrage.ArbitrageService{}
	ctx.JSON(http.StatusOK, service.GetCandlestickSnapshot([]string{}))
}
