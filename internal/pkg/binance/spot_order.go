package binance

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"

	"farmer/internal/pkg/errors"
	"farmer/internal/pkg/utils/logger"
	errPkg "farmer/pkg/errors"
)

func CreateSpotBuyOrder(ctx context.Context, client *binance.Client, symbol string, qty string, price string) (*binance.CreateOrderResponse, *errPkg.DomainError) {
	logger.Infof(ctx, "[CreateSpotBuyOrder] qty: %s | price: %s", qty, price)

	order, err := client.NewCreateOrderService().Symbol(symbol).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeFOK).Quantity(qty).Price(price).
		Do(ctx)
	if err != nil {
		logger.Error(ctx, err)
		domainErr := errors.NewDomainErrorCreateBuyOrderFailed(err)
		return nil, domainErr
	}

	if order.Status != binance.OrderStatusTypeFilled {
		err := errors.NewDomainErrorCreateBuyOrderFailed(fmt.Errorf("status: %s", order.Status))
		logger.Error(ctx, err)
		return nil, err
	}

	return order, nil
}

func CreateSpotSellOrder(ctx context.Context, client *binance.Client, symbol string, qty string, price string) (*binance.CreateOrderResponse, *errPkg.DomainError) {
	logger.Infof(ctx, "[CreateSpotSellOrder] qty: %s | price: %s", qty, price)

	order, err := client.NewCreateOrderService().Symbol(symbol).
		Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeFOK).Quantity(qty).Price(price).
		Do(ctx)
	if err != nil {
		logger.Error(ctx, err)
		domainErr := errors.NewDomainErrorCreateSellOrderFailed(err)
		return nil, domainErr
	}

	if order.Status != binance.OrderStatusTypeFilled {
		err := errors.NewDomainErrorCreateSellOrderFailed(fmt.Errorf("status: %s", order.Status))
		logger.Error(ctx, err)
		return nil, err
	}

	return order, nil
}
