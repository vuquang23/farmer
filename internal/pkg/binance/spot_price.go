package binance

import (
	"context"

	"github.com/adshao/go-binance/v2"

	"farmer/internal/pkg/errors"
	"farmer/internal/pkg/utils/logger"
	"farmer/internal/pkg/utils/maths"
	errPkg "farmer/pkg/errors"
)

func GetSpotPrice(ctx context.Context, client *binance.Client, symbol string) (float64, *errPkg.DomainError) {
	ret, err := client.NewAveragePriceService().Symbol(symbol).Do(ctx)
	if err != nil {
		domainErr := errors.NewDomainErrorGetPriceFailed(err)
		logger.Error(ctx, domainErr)
		return 0, domainErr
	}

	return maths.StrToFloat(ret.Price), nil
}
