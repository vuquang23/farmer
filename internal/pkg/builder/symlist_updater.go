package builder

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"

	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/utils/logger"
)

type ISymlistUpdater interface {
	Run(ctx *gin.Context, filePath string) error
}

type symlistUpdater struct{}

func NewSymlistUpdater() ISymlistUpdater {
	return &symlistUpdater{}
}

// Run update symbol list to file.
// Simple logic so not need to split to service or repo.
func (updater *symlistUpdater) Run(ctx *gin.Context, filePath string) error {
	fo, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fo.Close()

	res, err := binance.BinanceSpotClientInstance().
		NewGetAllCoinsInfoService().Do(context.Background())
	if err != nil {
		return err
	}

	for _, coin := range res {
		// check whether pair coin-usdt is existed
		_, err := binance.BinanceSpotClientInstance().
			NewKlinesService().Symbol(coin.Coin + "USDT").
			Interval("1d").Limit(1).
			Do(context.Background())
		if err != nil {
			logger.FromGinCtx(ctx).Error(err.Error())
			continue
		}
		fo.Write([]byte(coin.Coin + "\n"))
	}

	logger.FromGinCtx(ctx).Info("update symbol list successfully")

	return nil
}
