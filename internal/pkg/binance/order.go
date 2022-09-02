package binance

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"

	"farmer/internal/pkg/errors"
)

func CreateSpotBuyOrder(client *binance.Client, symbol string, qty string, price string) (*binance.CreateOrderResponse, error) {
	order, err := client.NewCreateOrderService().Symbol(symbol).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeFOK).Quantity(qty).Price(price).
		Do(context.Background())
	if err != nil {
		return nil, errors.NewDomainErrorCreateBuyOrderFailed(err)
	}

	if order.Status != binance.OrderStatusTypeFilled {
		return nil, errors.NewDomainErrorCreateBuyOrderFailed(fmt.Errorf("status: %s", order.Status))
	}

	return order, nil
}
