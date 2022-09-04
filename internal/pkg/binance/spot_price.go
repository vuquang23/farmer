package binance

import (
	"context"

	"github.com/adshao/go-binance/v2"

	"farmer/internal/pkg/errors"
	"farmer/internal/pkg/utils/maths"
	pkgErr "farmer/pkg/errors"
)

func GetSpotPrice(client *binance.Client, symbol string) (float64, *pkgErr.DomainError) {
	ret, err := client.NewAveragePriceService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return 0, errors.NewDomainErrorGetPriceFailed(err)
	}

	return maths.StrToFloat(ret.Price), nil
}
