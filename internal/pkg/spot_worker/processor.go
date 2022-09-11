package spotworker

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/adshao/go-binance/v2"

	b "farmer/internal/pkg/binance"
	c "farmer/internal/pkg/constants"
	en "farmer/internal/pkg/entities"
	e "farmer/internal/pkg/errors"
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
		w.analyzeWavetrendAndBuy()

		// check if should sell with wavetrend
		w.analyzeWavetrendAndSell()

		// check if should sell with exception
		w.analyzeExceptionsAndSell()

		// check health
		if time.Since(w.stt.loadHealth()) > time.Minute*30 {
			w.stt.storeHealth(time.Now())
		}
	}
}

func (w *spotWorker) analyzeExceptionsAndSell() {
	log := logger.WithDescription(fmt.Sprintf("%s - Analyze Exceptions And Sell", w.setting.symbol))

	sSignal, err := w.sellSignalExceptions()
	if err != nil {
		log.Sugar().Error(err)
		return
	}

	if sSignal.ShouldSell {
		res, err := w.createSellOrders(sSignal)
		if err != nil {
			log.Sugar().Error(err)
			return
		}

		if err := w.afterSell(res); err != nil {
			log.Sugar().Error(err)
			return
		}
	}
}

func (w *spotWorker) sellSignalExceptions() (*en.SpotSellSignal, error) {
	h1DiffWt := w.wavetrendProvider.GetCurrentDifWavetrend(wavetrendSvcName(w.setting.symbol, c.H1))

	currentPrice := w.wavetrendProvider.GetClosePrice(wavetrendSvcName(w.setting.symbol, c.M1))

	now := time.Now()
	buyOrders, err := w.spotTradeRepo.GetNotDoneBuyOrdersByWorkerIDAndCreatedAt(w.ID, now.Add(-time.Minute*100))
	if err != nil {
		return nil, err
	}

	orders := []*en.SpotSellOrder{}
	for _, b := range buyOrders {
		if h1DiffWt < 0 {
			if b.Price*(1+c.MinBenefit/100) <= currentPrice {
				orders = append(orders, &en.SpotSellOrder{
					Qty:        b.Qty,
					UnitBought: b.UnitBought,
					Ref:        b.ID,
					Price:      currentPrice,
				})
			}
		} else {
			distance := now.Sub(b.CreatedAt)
			ok := false

			switch {
			case distance <= 20*time.Minute:
				if b.Price*(1+(c.ExceptionBaseBenefitOnUpTrend+0*c.ExceptionStepBenefitOnUpTrend)/100) <= currentPrice {
					ok = true
				}
			case distance <= 40*time.Minute:
				if b.Price*(1+(c.ExceptionBaseBenefitOnUpTrend+1*c.ExceptionStepBenefitOnUpTrend)/100) <= currentPrice {
					ok = true
				}
			case distance <= 60*time.Minute:
				if b.Price*(1+(c.ExceptionBaseBenefitOnUpTrend+2*c.ExceptionStepBenefitOnUpTrend)/100) <= currentPrice {
					ok = true
				}
			case distance <= 80*time.Minute:
				if b.Price*(1+(c.ExceptionBaseBenefitOnUpTrend+3*c.ExceptionStepBenefitOnUpTrend)/100) <= currentPrice {
					ok = true
				}
			default:
				if b.Price*(1+(c.ExceptionBaseBenefitOnUpTrend+4*c.ExceptionStepBenefitOnUpTrend)/100) <= currentPrice {
					ok = true
				}
			}

			if ok {
				orders = append(orders, &en.SpotSellOrder{
					Qty:        b.Qty,
					UnitBought: b.UnitBought,
					Ref:        b.ID,
					Price:      currentPrice,
				})
			}
		}
	}

	if len(orders) == 0 {
		return &en.SpotSellSignal{ShouldSell: false}, nil
	}

	return &en.SpotSellSignal{
		ShouldSell: true,
		Orders:     orders,
	}, nil
}

func (w *spotWorker) analyzeWavetrendAndSell() {
	log := logger.WithDescription(fmt.Sprintf("%s - Analyze And Sell With Wavetrend", w.setting.symbol))

	sSignal, err := w.sellSignal()
	if err != nil {
		log.Sugar().Error(err)
		return
	}

	if sSignal.ShouldSell {
		res, err := w.createSellOrders(sSignal)
		if err != nil {
			log.Sugar().Error(err)
			return
		}

		if err := w.afterSell(res); err != nil {
			log.Sugar().Error(err)
			return
		}
	}
}

// TODO: re-invest benefit = 3/5?
func (w *spotWorker) createSellOrders(sSignal *en.SpotSellSignal) ([]*en.CreateSpotSellOrderResponse, error) {
	log := logger.WithDescription(fmt.Sprintf("%s - Create Sell Order", w.setting.symbol))

	m1SvcName := wavetrendSvcName(w.setting.symbol, c.M1)

	ret := []*en.CreateSpotSellOrderResponse{}
	notSold := sSignal.Orders

	// sell time lasts in 60 seconds
	ticker := time.NewTicker(time.Second)
	for i := 0; i < 60; i, _ = i+1, <-ticker.C {
		currentPrice := w.wavetrendProvider.GetClosePrice(m1SvcName)
		// down price 0.05%
		down := 0.05
		price := currentPrice * (1 - down/100)

		tempNotSold := []*en.SpotSellOrder{}
		for _, order := range notSold {
			if order.Price*(1-down*2/100) > price {
				log.Sugar().Infof("Current price is too low compared to expected. Expected: %f - Current: %f", order.Price, price)
				tempNotSold = append(tempNotSold, order)
				continue
			}

			log.Sugar().Infof(
				"Try to sell with Qty: %s - Price: %f - Current Tci: %f - DifWavetrend: %f",
				order.Qty, price, w.wavetrendProvider.GetCurrentTci(m1SvcName), w.wavetrendProvider.GetCurrentDifWavetrend(m1SvcName),
			)

			if res, err := b.CreateSpotSellOrder(
				w.bclient, w.setting.symbol, order.Qty,
				maths.RoundingUp(price, w.exchangeInf.loadPricePrecision()),
			); err == nil {
				log.Sugar().Infof(
					"Sell successfully. Qty: %s - Price: %f - Ref: %d - Current Tci: %f - DifWavetrend: %f",
					order.Qty, price, order.Ref, w.wavetrendProvider.GetCurrentTci(m1SvcName), w.wavetrendProvider.GetCurrentDifWavetrend(m1SvcName),
				)

				ret = append(ret, &en.CreateSpotSellOrderResponse{
					BinanceResponse: res,
					Order:           order,
				})
			} else {
				log.Sugar().Error(err)
				tempNotSold = append(tempNotSold, order)
			}
		}
		notSold = tempNotSold
		if len(notSold) == 0 {
			break
		}

		time.Sleep(time.Second)
	}

	return ret, nil
}

func (w *spotWorker) afterSell(res []*en.CreateSpotSellOrderResponse) error {
	log := logger.WithDescription(fmt.Sprintf("%s - After Sell", w.setting.symbol))

	buyIDs := []uint64{}
	sellTrades := []*en.SpotTrade{}
	updateUnitBought := 0

	for _, r := range res {
		updateUnitBought += int(r.Order.UnitBought)
		buyIDs = append(buyIDs, r.Order.Ref)
		sellTrades = append(sellTrades, &en.SpotTrade{
			Symbol:              w.setting.symbol,
			Side:                "SELL",
			BinanceOrderID:      uint64(r.BinanceResponse.OrderID),
			SpotWorkerID:        w.ID,
			Qty:                 r.BinanceResponse.ExecutedQuantity,
			CummulativeQuoteQty: maths.StrToFloat(r.BinanceResponse.CummulativeQuoteQuantity),
			Price:               maths.StrToFloat(r.BinanceResponse.Price),
			Ref:                 r.Order.Ref,
			IsDone:              true,
			UnitBought:          r.Order.UnitBought,
		})
	}

	w.stt.updateTotalUnitBought(int64(updateUnitBought))

	if err := w.spotTradeRepo.UpdateBuyOrders(buyIDs); err != nil {
		log.Sugar().Error(err)
	}

	if err := w.spotTradeRepo.CreateSellOrders(sellTrades); err != nil {
		log.Sugar().Error(err)
	}

	return nil
}

func (w *spotWorker) sellSignal() (*en.SpotSellSignal, error) {
	m1SvcName := wavetrendSvcName(w.setting.symbol, c.M1)

	shouldSell := w.shouldSell()
	if !shouldSell {
		return &en.SpotSellSignal{ShouldSell: false}, nil
	}

	ret := &en.SpotSellSignal{
		ShouldSell: true,
		Orders:     []*en.SpotSellOrder{},
	}
	currentPrice := w.wavetrendProvider.GetClosePrice(m1SvcName)

	trades, err := w.spotTradeRepo.GetNotDoneBuyOrdersByWorkerID(w.ID)
	if err != nil {
		return nil, err
	}
	for _, t := range trades {
		if t.Price*(1+c.MinBenefit/100) <= currentPrice { // min benefit is 0.5%
			ret.Orders = append(ret.Orders, &en.SpotSellOrder{
				Qty:        t.Qty,
				UnitBought: t.UnitBought,
				Ref:        t.ID,
				Price:      currentPrice,
			})
		}
	}

	return ret, nil
}

func (w *spotWorker) shouldSell() bool {
	m1SvcName := wavetrendSvcName(w.setting.symbol, c.M1)

	currentTci := w.wavetrendProvider.GetCurrentTci(m1SvcName)
	if currentTci < c.WavetrendOverbought {
		return false
	}

	currentDifWt := w.wavetrendProvider.GetCurrentDifWavetrend(m1SvcName)
	if currentDifWt >= 0 {
		return false
	}

	pastWtDat, isOutdated := w.wavetrendProvider.GetPastWaveTrendData(m1SvcName)
	if pastWtDat == nil { // error
		return false
	}

	if isOutdated {
		logger.Logger.Error("past wavetrend data is outdated")
		return false
	}

	for i := len(pastWtDat.PastTci) - c.OverboughtRequiredTime; i < len(pastWtDat.PastTci); i++ {
		if pastWtDat.PastTci[i] < c.WavetrendOverbought {
			return false
		}
	}

	for i := len(pastWtDat.DifWavetrend) - c.OverboughtNegativeDifWtRequiredTime - c.OverboughtPositiveDifWtRequiredTime; i < len(pastWtDat.DifWavetrend); i++ {
		if i < len(pastWtDat.DifWavetrend)-c.OverboughtNegativeDifWtRequiredTime {
			if pastWtDat.DifWavetrend[i] < 0 {
				return false
			}
		} else {
			if pastWtDat.DifWavetrend[i] >= 0 {
				return false
			}
		}
	}

	return true
}

func (w *spotWorker) analyzeWavetrendAndBuy() {
	log := logger.WithDescription(fmt.Sprintf("%s - Analyze And Buy With Wavetrend", w.setting.symbol))

	bSignal, err := w.buySignal()
	if err != nil {
		log.Sugar().Error(err)
		return
	}

	if bSignal.ShouldBuy {
		res, err := w.createBuyOrder(bSignal)
		if err != nil {
			log.Sugar().Error(err)
			return
		}

		if w.afterBuy(res, bSignal.Order.UnitBought); err != nil {
			log.Sugar().Error(err)
		}
	}
}

func (w *spotWorker) afterBuy(res *binance.CreateOrderResponse, unitBought int64) {
	log := logger.WithDescription(fmt.Sprintf("%s - After Buy", w.setting.symbol))

	w.stt.updateTotalUnitBought(-unitBought)
	w.stt.storeLastBoughtAt(time.Now())

	// update DB
	if err := w.spotTradeRepo.CreateBuyOrder(en.SpotTrade{
		Symbol:              w.setting.symbol,
		Side:                "BUY",
		BinanceOrderID:      uint64(res.OrderID),
		SpotWorkerID:        w.ID,
		Qty:                 res.ExecutedQuantity,
		CummulativeQuoteQty: maths.StrToFloat(res.CummulativeQuoteQuantity),
		Price:               maths.StrToFloat(res.Price),
		IsDone:              false,
		UnitBought:          uint64(unitBought),
	}); err != nil {
		log.Sugar().Error(err)
	}
}

func (w *spotWorker) createBuyOrder(bSignal *en.SpotBuySignal) (*binance.CreateOrderResponse, *pkgErr.DomainError) {
	log := logger.WithDescription(fmt.Sprintf("%s - Create Buy Order", w.setting.symbol))
	m1SvcName := wavetrendSvcName(w.setting.symbol, c.M1)

	// buy time lasts in 20 seconds
	ticker := time.NewTicker(time.Second)
	for i := 0; i < 20; i, _ = i+1, <-ticker.C {
		currentPrice := w.wavetrendProvider.GetClosePrice(m1SvcName)
		up := 0.05
		price := currentPrice * (1 + up/100)

		if (bSignal.Order.Price * (1 + (up*2)/100)) < price {
			log.Sugar().Info("Current price is too high compared to expected. Expected: %f - Current: %f", bSignal.Order.Price, price)
			continue
		}

		notional := w.setting.loadUnitNotional() * float64(bSignal.Order.UnitBought)
		qty := notional / price
		if notional < w.exchangeInf.loadMinNotional() || qty < w.exchangeInf.loadMinQty() {
			log.Error("Not enough notional or qty")
			continue
		}

		log.Sugar().Infof(
			"Try to buy with Notional: %f - Price: %f - Current Tci: %f - Current difwavetrend: %f",
			notional, price, w.wavetrendProvider.GetCurrentTci(m1SvcName), w.wavetrendProvider.GetCurrentDifWavetrend(m1SvcName),
		)

		if res, err := b.CreateSpotBuyOrder(
			w.bclient, w.setting.symbol,
			maths.RoundingUp(qty, w.exchangeInf.loadQtyPrecision()),
			maths.RoundingUp(price, w.exchangeInf.loadPricePrecision()),
		); err == nil {
			log.Sugar().Infof(
				"Buy successfully. Notional: %f - Price: %f - Current Tci: %f - Current difwavetrend: %f",
				notional, price, w.wavetrendProvider.GetCurrentTci(m1SvcName), w.wavetrendProvider.GetCurrentDifWavetrend(m1SvcName),
			)

			return res, nil
		} else {
			log.Sugar().Error(err)
		}
	}

	return nil, e.NewDomainErrorCreateBuyOrderFailed(nil)
}

func (w *spotWorker) buySignal() (*en.SpotBuySignal, error) {
	var unitBought int64

	shouldBuy := w.shouldBuy()
	if !shouldBuy {
		return &en.SpotBuySignal{ShouldBuy: false}, nil
	}

	h1DiffWt := w.wavetrendProvider.GetCurrentDifWavetrend(wavetrendSvcName(w.setting.symbol, c.H1))
	if h1DiffWt <= 0 {
		unitBought = int64(math.Min(c.UnitBuyOnDowntrend, float64(w.setting.loadUnitBuyAllowed())-float64(w.stt.loadTotalUnitBought())))
	} else {
		unitBought = int64(math.Min(c.UnitBuyOnUpTrend, float64(w.setting.loadUnitBuyAllowed())-float64(w.stt.loadTotalUnitBought())))
	}

	if unitBought == 0 {
		return nil, errors.New("remain 0 unit to buy")
	}

	logger.Logger.Sugar().Infof("[Buy Signal] Current h1DiffWt: %f - unitBought: %d", h1DiffWt, unitBought)

	currentPrice := w.wavetrendProvider.GetClosePrice(wavetrendSvcName(w.setting.symbol, c.M1))
	return &en.SpotBuySignal{
		ShouldBuy: true,
		Order: en.SpotBuyOrder{
			UnitBought: unitBought,
			Price:      currentPrice,
		},
	}, nil
}

func (w *spotWorker) shouldBuy() bool {
	// check status
	if w.setting.loadUnitBuyAllowed() == uint64(w.stt.loadTotalUnitBought()) {
		return false
	}

	if time.Since(w.stt.loadLastBoughtAt()) <= c.StopBuyAfterBuy {
		return false
	}

	// check wavetrend
	m1SvcName := wavetrendSvcName(w.setting.symbol, c.M1)

	currentTci := w.wavetrendProvider.GetCurrentTci(m1SvcName)
	if currentTci > c.WavetrendOversold {
		return false
	}

	currentDifWt := w.wavetrendProvider.GetCurrentDifWavetrend(m1SvcName)
	if currentDifWt <= 0 {
		return false
	}

	pastWtDat, isOutdated := w.wavetrendProvider.GetPastWaveTrendData(m1SvcName)
	if pastWtDat == nil { // get error
		return false
	}

	if isOutdated {
		logger.Logger.Error("past wavetrend data is outdated")
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
