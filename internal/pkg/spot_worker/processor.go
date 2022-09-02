package spotworker

import (
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"

	b "farmer/internal/pkg/binance"
	c "farmer/internal/pkg/constants"
	en "farmer/internal/pkg/entities"
	"farmer/internal/pkg/errors"
	"farmer/internal/pkg/utils/logger"
	"farmer/internal/pkg/utils/maths"
	pkgErr "farmer/pkg/errors"
)

func (w *spotWorker) runMainProcessor() {
	log := logger.WithDescription(fmt.Sprintf("%s - Main Proccessor", w.setting.symbol))
	log.Sugar().Infof("Worker started")

	ticker := time.NewTicker(c.SleepAfterProcessing)
	for ; !w.getStopSignal(); <-ticker.C {
		// check if should buy
		bSignal, err := w.buySignal()
		if err != nil {
			log.Sugar().Error(err)
			continue
		}
		if bSignal.ShouldBuy {
			res, err := w.createBuyOrder(bSignal)
			if err != nil {
				log.Sugar().Error(err)
				continue
			}

			log.Sugar().Infof("%+v\n", res)

			w.afterBuy(res, bSignal.Order.UnitBought)
		}

		// check if should sell
	}
}

func (w *spotWorker) afterBuy(res *binance.CreateOrderResponse, unitBought int64) {
	w.stt.updateTotalUnitBought(-unitBought)
	w.stt.storeLastBoughtAt(uint64(time.Now().Unix()))

	// update DB
	w.spotTradeRepo.CreateBuyOrder(en.SpotTrade{
		Side:                "BUY",
		BinanceOrderID:      uint64(res.OrderID),
		SpotWorkerID:        w.ID,
		Qty:                 maths.StrToFloat(res.ExecutedQuantity),
		CummulativeQuoteQty: maths.StrToFloat(res.CummulativeQuoteQuantity),
		IsDone:              false,
		UnitBought:          uint64(unitBought),
	})
}

func (w *spotWorker) createBuyOrder(bSignal *en.BuySignal) (*binance.CreateOrderResponse, *pkgErr.DomainError) {
	log := logger.WithDescription(fmt.Sprintf("%s - Create Buy Order", w.setting.symbol))

	// buy time lasts in 10 seconds
	ticker := time.NewTicker(time.Second)
	for i := 0; i < 10; i, _ = i+1, <-ticker.C {
		currentPrice := w.wavetrendProvider.GetClosePrice(wavetrendSvcName(w.setting.symbol, c.M1))
		percentUp := 0.0

		for j := 0; j < 5; j++ {
			// up price
			percentUp += 0.05
			price := currentPrice * (1 + percentUp/100)
			notional := w.setting.loadUnitNotional() * float64(bSignal.Order.UnitBought)
			qty := notional / price
			if notional < w.exchangeInf.loadMinNotional() || qty < w.exchangeInf.loadMinQty() {
				log.Error("Not enough notional or qty")
				continue
			}

			res, err := b.CreateSpotBuyOrder(
				w.bclient, w.setting.symbol,
				maths.RoundingUp(qty, w.exchangeInf.loadQtyPrecision()),
				maths.RoundingUp(price, w.exchangeInf.loadPricePrecision()),
			)
			if err == nil {
				return res, nil
			} else {
				log.Sugar().Error(err)
			}

			time.Sleep(time.Second / 5)
		}
	}

	return nil, errors.NewDomainErrorCreateBuyOrderFailed(nil)
}

func (w *spotWorker) buySignal() (*en.BuySignal, error) {
	shouldBuy := w.shouldBuy()
	if !shouldBuy {
		return &en.BuySignal{ShouldBuy: false}, nil
	}

	ret := &en.BuySignal{
		ShouldBuy: true,
	}

	h1DiffWt := w.wavetrendProvider.GetCurrentDifWavetrend(wavetrendSvcName(w.setting.symbol, c.H1))
	if h1DiffWt <= 0 {
		ret.Order = en.BuyOrder{
			UnitBought: c.UnitBuyOnDowntrend,
		}
	} else {
		ret.Order = en.BuyOrder{
			UnitBought: c.UnitBuyOnUpTrend,
		}
	}

	return ret, nil
}

func (w *spotWorker) shouldBuy() bool {
	// buy status
	if w.setting.loadUnitBuyAllowed() == uint64(w.stt.loadTotalUnitBought()) {
		return false
	}

	now := time.Now().Unix()
	if now-int64(w.stt.loadLastBoughtAt()) <= c.StopBuyAfterBuy {
		return false
	}

	// by wavetrend
	m1SvcName := wavetrendSvcName(w.setting.symbol, c.M1)

	currentTci := w.wavetrendProvider.GetCurrentTci(m1SvcName)
	if currentTci > c.WavetrendOversold {
		return false
	}

	currentDifWt := w.wavetrendProvider.GetCurrentDifWavetrend(m1SvcName)
	if currentDifWt <= 0 {
		return false
	}

	pastWtDat := w.wavetrendProvider.GetPastWaveTrendData(m1SvcName)
	if pastWtDat == nil { // get error
		return false
	}

	for i := len(pastWtDat.PastTci) - c.OversoldRequiredTime; i < len(pastWtDat.PastTci); i++ {
		if pastWtDat.PastTci[i] > c.WavetrendOversold {
			return false
		}
	}

	for i := len(pastWtDat.DifWavetrend) - c.OversoldNegativeDifWtRequiredTime - c.OversoldPositiveDifWtRequiredTime; i < len(pastWtDat.DifWavetrend); i++ {
		if i < len(pastWtDat.DifWavetrend)-c.OversoldPositiveDifWtRequiredTime {
			if pastWtDat.DifWavetrend[i] > 0 {
				return false
			}
		} else {
			if pastWtDat.DifWavetrend[i] <= 0 {
				return false
			}
		}
	}

	return true
}
