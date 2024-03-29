package spotmanager

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/errors"
	"farmer/internal/pkg/utils/logger"
	"farmer/internal/pkg/utils/maths"
)

func (m *spotManager) updateExchangeInfoPeriodically(ctx context.Context, doneC chan<- error) {
	logger.Info(ctx, "[updateExchangeInfoPeriodically] start updating exchange info")

	once := &sync.Once{}
	ticker := time.NewTicker(time.Hour)
	isInit := true
	for ; true; <-ticker.C {
		if err := m.updateExchangeInfo(ctx); err != nil {
			domainErr := errors.NewDomainErrorInitExchangeInfo(err)
			if isInit {
				doneC <- domainErr
				return
			}
			continue
		}

		once.Do(func() {
			isInit = false
			doneC <- nil
		})
	}
}

func (m *spotManager) updateExchangeInfo(ctx context.Context) error {
	exchangeInfo, err := m.bclient.NewExchangeInfoService().Do(ctx)
	if err != nil {
		logger.Error(ctx, err)
		return err
	}
	for _, s := range exchangeInfo.Symbols {
		var (
			pricePrecision, qtyPrecision int
			minNotional, minQty          float64
			err                          error
		)
		for _, f := range s.Filters {
			switch f["filterType"].(string) {
			case "PRICE_FILTER":
				tickSize, ok := f["tickSize"].(string)
				if !ok {
					err := fmt.Errorf("[updateExchangeInfo] tickSize: can not cast to string")
					logger.Error(ctx, err)
					return err
				}
				pricePrecision, err = maths.GetPrecision(tickSize)
				if err != nil {
					logger.Error(ctx, err)
					return err
				}
			case "LOT_SIZE":
				stepSize, ok := f["stepSize"].(string)
				if !ok {
					err := fmt.Errorf("[updateExchangeInfo] stepSize: can not cast to string")
					logger.Error(ctx, err)
					return err
				}
				qtyPrecision, err = maths.GetPrecision(stepSize)
				if err != nil {
					logger.Error(ctx, err)
					return err
				}

				minQtyStr, ok := f["minQty"].(string)
				if !ok {
					err := fmt.Errorf("[updateExchangeInfo] minQty: can not cast to string")
					logger.Error(ctx, err)
					return err
				}
				minQty, err = strconv.ParseFloat(minQtyStr, 64)
				if err != nil {
					logger.Error(ctx, err)
					return err
				}
			case "MIN_NOTIONAL":
				minNotionalStr, ok := f["minNotional"].(string)
				if !ok {
					err := fmt.Errorf("[updateExchangeInfo] minNotional: can not cast to string")
					logger.Error(ctx, err)
					return err
				}
				minNotional, err = strconv.ParseFloat(minNotionalStr, 64)
				if err != nil {
					logger.Error(ctx, err)
					return err
				}
			}
		}

		exchangeInfo := entities.SpotExchangeInfo{
			PricePrecision: pricePrecision,
			QtyPrecision:   qtyPrecision,
			MinQty:         minQty,
			MinNotional:    minNotional,
		}
		m.SetExchangeInfo(s.Symbol, exchangeInfo)

		worker, ok := m.mapSymbolWorker[s.Symbol]
		if ok {
			worker.SetExchangeInfo(ctx, exchangeInfo)
		}

	}
	return nil
}

func (m *spotManager) SetExchangeInfo(symbol string, info entities.SpotExchangeInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.mapExchangeInfo[symbol] = info
}

func (m *spotManager) GetExchangeInfo(symbol string) (entities.SpotExchangeInfo, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ret, ok := m.mapExchangeInfo[symbol]
	return ret, ok
}
