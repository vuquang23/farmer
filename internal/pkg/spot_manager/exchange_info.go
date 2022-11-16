package spotmanager

import (
	"context"
	"strconv"
	"sync"
	"time"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/errors"
	"farmer/internal/pkg/utils/logger"
	"farmer/internal/pkg/utils/maths"
)

func (m *spotManager) updateExchangeInfoPeriodically(doneC chan<- error) {
	log := logger.WithDescription("[updateExchangeInfoPeriodically] spotManager")

	once := &sync.Once{}
	ticker := time.NewTicker(time.Hour)
	isInit := true
	for ; true; <-ticker.C {
		if err := m.updateExchangeInfo(); err != nil {
			domainErr := errors.NewDomainErrorInitExchangeInfo(err)
			if isInit {
				doneC <- domainErr
				return
			}
			log.Sugar().Warn(domainErr)
			continue
		}

		once.Do(func() {
			isInit = false
			doneC <- nil
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

		worker.SetExchangeInfo(entities.SpotExchangeInfo{
			PricePrecision: pricePrecision,
			QtyPrecision:   qtyPrecision,
			MinQty:         minQty,
			MinNotional:    minNotional,
		})
	}
	return nil
}
