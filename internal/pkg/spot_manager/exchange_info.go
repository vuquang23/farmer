package spotmanager

import (
	"context"
	"strconv"
	"sync"
	"time"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/utils/logger"
	"farmer/internal/pkg/utils/maths"
)

func (m *spotManager) updateExchangeInfoPeriodically(doneC chan<- struct{}) {
	logger := logger.WithDescription("Manager updates exchange info periodically")

	once := &sync.Once{}
	ticker := time.NewTicker(time.Hour)
	for ; true; <-ticker.C {
		if err := m.updateExchangeInfo(); err != nil {
			logger.Sugar().Error(err)
			// FIXME: now assume always update successfully on the first time
			continue
		}
		once.Do(func() {
			doneC <- struct{}{}
		})
	}
}

func (m *spotManager) updateExchangeInfo() error {
	exchangeInfo, err := m.bclient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return err
	}
	for _, s := range exchangeInfo.Symbols {
		worker, ok := m.mapSymbolWorker[s.Symbol]
		if !ok {
			continue
		}

		var (
			pricePrecision, qtyPrecision int
			minNotional, minQty          float64
			err                          error
		)
		for _, f := range s.Filters {
			switch f["filterType"].(string) {
			case "PRICE_FILTER":
				tickSize := f["tickSize"].(string)
				pricePrecision, err = maths.GetPrecision(tickSize)
				if err != nil {
					return err
				}
			case "LOT_SIZE":
				stepSize := f["stepSize"].(string)
				qtyPrecision, err = maths.GetPrecision(stepSize)
				if err != nil {
					return err
				}
			case "MIN_NOTIONAL":
				minNotional, err = strconv.ParseFloat(f["minNotional"].(string), 64)
				if err != nil {
					return err
				}
			}
		}

		worker.SetExchangeInfo(entities.ExchangeInfo{
			PricePrecision: pricePrecision,
			QtyPrecision:   qtyPrecision,
			MinQty:         minQty,
			MinNotional:    minNotional,
		})
	}
	return nil
}
