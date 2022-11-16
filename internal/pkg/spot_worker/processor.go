package spotworker

import (
	"context"
	"errors"
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

func (w *spotWorker) runMainProcessor(ctx context.Context) {
	logger.Infof(ctx, "[runMainProcessor] worker started")

	ticker := time.NewTicker(c.SleepAfterProcessing)
	for ; !w.getStopSignal(); <-ticker.C {
		// check if should buy
		w.analyzeWavetrendAndBuy(ctx)

		// check if should sell with wavetrend
		w.analyzeWavetrendAndSell(ctx)

		// check if should sell with exception
		w.analyzeExceptionsAndSell(ctx)

		// check health
		if time.Since(w.stt.loadHealth()) > time.Minute*30 {
			w.stt.storeHealth(time.Now())
		}
	}
}

func (w *spotWorker) analyzeExceptionsAndSell(ctx context.Context) {
	sSignal, err := w.sellSignalExceptions(ctx)
	if err != nil {
		return
	}

	if sSignal.ShouldSell {
		res, err := w.createSellOrders(ctx, sSignal)
		if err != nil {
			return
		}

		if err := w.afterSell(ctx, res); err != nil {
			return
		}
	}
}

func (w *spotWorker) sellSignalExceptions(ctx context.Context) (*en.SpotSellSignal, error) {
	h1DifWt, isOutdated := w.wavetrendProvider.GetCurrentDifWavetrend(wavetrendSvcName(w.setting.symbol, c.H1))
	if isOutdated {
		err := errors.New("h1DifWt is outdated")
		logger.Warnf(ctx, "[sellSignalExceptions] %s", err)
		return nil, err
	}

	currentPrice, isOutdated := w.wavetrendProvider.GetClosePrice(wavetrendSvcName(w.setting.symbol, c.M1))
	if isOutdated {
		err := errors.New("currentPrice is outdated")
		logger.Warnf(ctx, "[sellSignalExceptions] %s", err)
		return nil, err
	}

	buyOrders, err := w.spotTradeRepo.GetNotDoneBuyOrdersByWorkerIDAndCreatedAtGT(w.ID, time.Now().Add(-time.Minute*100))
	if err != nil {
		return nil, err
	}

	var orders []*en.SpotSellOrder
	for _, b := range buyOrders {
		if h1DifWt < 0 {
			if b.Price*(1+c.MinBenefit/100) <= currentPrice {
				orders = append(orders, &en.SpotSellOrder{
					Qty:        b.Qty,
					UnitBought: b.UnitBought,
					Price:      currentPrice,
					Ref:        b,
				})
			}
		} else {
			distance := time.Since(b.CreatedAt)
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
					Ref:        b,
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

func (w *spotWorker) analyzeWavetrendAndSell(ctx context.Context) {
	sSignal, err := w.sellSignal(ctx)
	if err != nil {
		return
	}

	if sSignal.ShouldSell {
		res, err := w.createSellOrders(ctx, sSignal)
		if err != nil {
			return
		}

		// not have any sucess sell orders
		if len(res) == 0 {
			return
		}

		if err := w.afterSell(ctx, res); err != nil {
			return
		}
	}
}

func (w *spotWorker) createSellOrders(ctx context.Context, sSignal *en.SpotSellSignal) ([]*en.CreateSpotSellOrderResponse, error) {
	m1SvcName := wavetrendSvcName(w.setting.symbol, c.M1)

	ret := []*en.CreateSpotSellOrderResponse{}
	notSold := sSignal.Orders

	// sell time lasts in 60 seconds
	ticker := time.NewTicker(time.Second)
	for i := 0; i < 60; i, _ = i+1, <-ticker.C {
		currentPrice, isOutdated := w.wavetrendProvider.GetClosePrice(m1SvcName)
		if isOutdated {
			logger.Warn(ctx, "[createSellOrders] closePrice is outdated")
			continue
		}

		// down price 0.05%
		down := 0.05
		price := currentPrice * (1 - down/100)

		tempNotSold := []*en.SpotSellOrder{}
		for _, order := range notSold {
			if order.Price*(1-down*2/100) > price {
				logger.Infof(ctx, "[createSellOrders] current price is too low compared to expected. expected: %f | current: %f", order.Price, price)
				tempNotSold = append(tempNotSold, order)
				continue
			}

			// ignore outdated status here.
			currentTci, _ := w.wavetrendProvider.GetCurrentTci(m1SvcName)
			currentDifWavetrend, _ := w.wavetrendProvider.GetCurrentDifWavetrend(m1SvcName)

			logger.Infof(
				ctx, "[createSellOrders] try to sell with qty: %s | price: %f | M1 currentTci: %f | M1 currentDifWavetrend: %f",
				order.Qty, price, currentTci, currentDifWavetrend,
			)

			if res, err := b.CreateSpotSellOrder(
				w.bclient, w.setting.symbol, order.Qty,
				maths.RoundingUp(price, w.exchangeInf.loadPricePrecision()),
			); err == nil {
				logger.Infof(
					ctx, "[createSellOrders] sell successfully. qty: %s | price: %f | ref: %d | M1 currentTci: %f | M1 currentDifWavetrend: %f",
					order.Qty, price, order.Ref, currentTci, currentDifWavetrend,
				)

				ret = append(ret, &en.CreateSpotSellOrderResponse{
					BinanceResponse: res,
					Order:           order,
				})
			} else {
				logger.Warnf(ctx, "[createSellOrders] %s", err)
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

func (w *spotWorker) afterSell(ctx context.Context, res []*en.CreateSpotSellOrderResponse) error {
	var (
		sellTrades       []*en.SpotTrade
		updateUnitBought = 0
		benefit          = float64(0)
	)

	for _, r := range res {
		updateUnitBought += int(r.Order.UnitBought)
		sellTrades = append(sellTrades, &en.SpotTrade{
			Symbol:              w.setting.symbol,
			Side:                "SELL",
			BinanceOrderID:      uint64(r.BinanceResponse.OrderID),
			SpotWorkerID:        w.ID,
			Qty:                 r.BinanceResponse.ExecutedQuantity,
			CummulativeQuoteQty: maths.StrToFloat(r.BinanceResponse.CummulativeQuoteQuantity),
			Price:               maths.StrToFloat(r.BinanceResponse.Price),
			Ref:                 r.Order.Ref.ID,
			IsDone:              true,
			UnitBought:          r.Order.UnitBought,
		})

		benefit += maths.StrToFloat(r.BinanceResponse.CummulativeQuoteQuantity) - r.Order.Ref.CummulativeQuoteQty
	}

	w.stt.updateTotalUnitBought(-int64(updateUnitBought))

	// update unit notional
	unitBuyAllowed := w.setting.loadUnitBuyAllowed()
	val := benefit / float64(unitBuyAllowed)
	/// in memory
	w.setting.updateUnitNotional(val)
	/// in db
	for i := 0; i < 3; i++ {
		err := w.spotWorkerRepo.UpdateUnitNotionalByID(ctx, w.ID, val)
		if err == nil {
			break
		}
		logger.Infof(ctx, "[afterSell] retry %d...", i+1)
		time.Sleep(time.Second)
	}

	// should not be failed here
	for i := 0; i < 3; i++ {
		err := w.spotTradeRepo.CreateSellOrders(ctx, sellTrades)
		if err == nil {
			break
		}
		logger.Infof(ctx, "[afterSell] retry %d...", i+1)
		time.Sleep(time.Second)
	}

	return nil
}

func (w *spotWorker) sellSignal(ctx context.Context) (*en.SpotSellSignal, error) {
	m1SvcName := wavetrendSvcName(w.setting.symbol, c.M1)

	shouldSell := w.shouldSell(ctx)
	if !shouldSell {
		return &en.SpotSellSignal{ShouldSell: false}, nil
	}

	ret := &en.SpotSellSignal{
		ShouldSell: true,
		Orders:     []*en.SpotSellOrder{},
	}

	// ignore isOutdated here because wavetrend and difWavetrend in func
	// shouldSell is updated nearly the same time with closePrice
	currentPrice, _ := w.wavetrendProvider.GetClosePrice(m1SvcName)

	trades, err := w.spotTradeRepo.GetNotDoneBuyOrdersByWorkerID(ctx, w.ID)
	if err != nil {
		return nil, err
	}

	// not exist not-sold buy orders in DB
	if len(trades) == 0 {
		return &en.SpotSellSignal{ShouldSell: false}, nil
	}

	for _, t := range trades {
		if t.Price*(1+c.MinBenefit/100) <= currentPrice { // min benefit is 0.5%
			ret.Orders = append(ret.Orders, &en.SpotSellOrder{
				Qty:        t.Qty,
				UnitBought: t.UnitBought,
				Ref:        t,
				Price:      currentPrice,
			})
		}
	}

	// all not-sold buy orders not reach min expected benefit
	if len(ret.Orders) == 0 {
		return &en.SpotSellSignal{ShouldSell: false}, nil
	}

	return ret, nil
}

func (w *spotWorker) shouldSell(ctx context.Context) bool {
	m1SvcName := wavetrendSvcName(w.setting.symbol, c.M1)

	currentTci, isOutdated := w.wavetrendProvider.GetCurrentTci(m1SvcName)
	if isOutdated {
		logger.Warn(ctx, "[shouldSell] currentTci M1 is outdated")
		return false
	}
	if currentTci < c.WavetrendOverbought {
		return false
	}

	currentDifWt, isOutdated := w.wavetrendProvider.GetCurrentDifWavetrend(m1SvcName)
	if isOutdated {
		logger.Warn(ctx, "[shouldSell] currentDifWt M1 is outdated")
		return false
	}
	if currentDifWt >= 0 {
		return false
	}

	pastWtDat, isOutdated := w.wavetrendProvider.GetPastWaveTrendData(m1SvcName)
	if pastWtDat == nil { // should not nil here
		return false
	}
	if isOutdated {
		logger.Warn(ctx, "[shouldSell] pastWtDat M1 is outdated")
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

func (w *spotWorker) analyzeWavetrendAndBuy(ctx context.Context) {
	bSignal, err := w.buySignal(ctx)
	if err != nil {
		return
	}

	if bSignal.ShouldBuy {
		res, err := w.createBuyOrder(ctx, bSignal)
		if err != nil {
			return
		}

		if w.afterBuy(ctx, res, bSignal.Order.UnitBought); err != nil {
			return
		}
	}
}

func (w *spotWorker) afterBuy(ctx context.Context, res *binance.CreateOrderResponse, unitBought int64) {
	w.stt.updateTotalUnitBought(unitBought)
	w.stt.storeLastBoughtAt(time.Now())

	// update DB
	for i := 0; i < 3; i++ {
		err := w.spotTradeRepo.CreateBuyOrder(ctx, en.SpotTrade{
			Symbol:              w.setting.symbol,
			Side:                "BUY",
			BinanceOrderID:      uint64(res.OrderID),
			SpotWorkerID:        w.ID,
			Qty:                 res.ExecutedQuantity,
			CummulativeQuoteQty: maths.StrToFloat(res.CummulativeQuoteQuantity),
			Price:               maths.StrToFloat(res.Price),
			IsDone:              false,
			UnitBought:          uint64(unitBought),
		})
		if err == nil {
			break
		}
		logger.Infof(ctx, "retry %d...", i+1)
		time.Sleep(time.Second)
	}
}

func (w *spotWorker) createBuyOrder(ctx context.Context, bSignal *en.SpotBuySignal) (*binance.CreateOrderResponse, *pkgErr.DomainError) {
	m1SvcName := wavetrendSvcName(w.setting.symbol, c.M1)

	// buy time lasts in 20 seconds
	ticker := time.NewTicker(time.Second)
	for i := 0; i < 20; i, _ = i+1, <-ticker.C {
		// ignore isOutdated here because wavetrend and difWavetrend
		// in func shouldBuy is updated in nearly the same time with closePrice
		currentPrice, _ := w.wavetrendProvider.GetClosePrice(m1SvcName)
		up := 0.05
		price := currentPrice * (1 + up/100)

		if (bSignal.Order.Price * (1 + (up*2)/100)) < price {
			logger.Infof(ctx, "[createBuyOrder] current price is too high compared to expected. expected: %f | current: %f", bSignal.Order.Price, price)
			continue
		}

		notional := w.setting.loadUnitNotional() * float64(bSignal.Order.UnitBought)
		qty := notional / price
		if notional < w.exchangeInf.loadMinNotional() || qty < w.exchangeInf.loadMinQty() {
			logger.Warn(ctx, "[createBuyOrder] not enough notional or qty")
			continue
		}

		// ignore isOutdated here
		currentTci, _ := w.wavetrendProvider.GetCurrentTci(m1SvcName)
		currentDifWavetrend, _ := w.wavetrendProvider.GetCurrentDifWavetrend(m1SvcName)

		logger.Infof(
			ctx, "[createBuyOrder] try to buy with notional: %f | price: %f | M1 current tci: %f | M1 current difwavetrend: %f",
			notional, price, currentTci, currentDifWavetrend,
		)

		if res, err := b.CreateSpotBuyOrder(
			w.bclient, w.setting.symbol,
			maths.RoundingUp(qty, w.exchangeInf.loadQtyPrecision()),
			maths.RoundingUp(price, w.exchangeInf.loadPricePrecision()),
		); err == nil {
			logger.Infof(
				ctx, "[createBuyOrder] buy successfully. notional: %f | price: %f | M1 current tci: %f | M1 current difwavetrend: %f",
				notional, price, currentTci, currentDifWavetrend,
			)

			return res, nil
		} else {
			logger.Error(ctx, err)
		}
	}

	return nil, e.NewDomainErrorCreateBuyOrderFailed(nil)
}

func (w *spotWorker) buySignal(ctx context.Context) (*en.SpotBuySignal, error) {
	shouldBuy := w.shouldBuy(ctx)
	if !shouldBuy {
		return &en.SpotBuySignal{ShouldBuy: false}, nil
	}

	unitBought, err := w.determineUnitNumberToBuy(ctx)
	if err != nil {
		return nil, err
	}
	if unitBought <= 0 {
		return &en.SpotBuySignal{ShouldBuy: false}, nil
	}

	// ignore isOutdated here because wavetrend and difWavetrend
	// in func shouldBuy is updated in nearly the same time with closePrice
	currentPrice, _ := w.wavetrendProvider.GetClosePrice(wavetrendSvcName(w.setting.symbol, c.M1))
	return &en.SpotBuySignal{
		ShouldBuy: true,
		Order: en.SpotBuyOrder{
			UnitBought: unitBought,
			Price:      currentPrice,
		},
	}, nil
}

func (w *spotWorker) determineUnitNumberToBuy(ctx context.Context) (int64, error) {
	var (
		isUptrend  bool
		unitBought int64

		h1SvcName = wavetrendSvcName(w.setting.symbol, c.H1)
	)

	pastWt, isOutdated := w.wavetrendProvider.GetPastWaveTrendData(h1SvcName)
	if isOutdated {
		err := errors.New("h1DifWt is outdated")
		logger.Warnf(ctx, "[determineUnitNumberToBuy] %s", err)
		return 0, err
	}

	isUptrend = true
	for i := len(pastWt.DifWavetrend) - c.IsUptrendOnH1RequiredTime; i < len(pastWt.DifWavetrend); i++ {
		if pastWt.DifWavetrend[i] <= 0 {
			isUptrend = false
			break
		}
	}

	if isUptrend {
		unitBought = int64(math.Min(c.UnitBuyOnUpTrend, float64(w.setting.loadUnitBuyAllowed())-float64(w.stt.loadTotalUnitBought())))
	} else {
		unitBought = int64(math.Min(c.UnitBuyOnDowntrend, float64(w.setting.loadUnitBuyAllowed())-float64(w.stt.loadTotalUnitBought())))
	}

	logger.Infof(ctx, "[determineUnitNumberToBuy] current h1DifWt slice: %v | unitBought: %d", pastWt.DifWavetrend[len(pastWt.DifWavetrend)-c.IsUptrendOnH1RequiredTime-1:], unitBought)

	return unitBought, nil
}

func (w *spotWorker) shouldBuy(ctx context.Context) bool {
	// check status
	if w.setting.loadUnitBuyAllowed() == uint64(w.stt.loadTotalUnitBought()) {
		return false
	}

	if time.Since(w.stt.loadLastBoughtAt()) <= c.StopBuyAfterBuy {
		return false
	}

	// check wavetrend
	m1SvcName := wavetrendSvcName(w.setting.symbol, c.M1)

	currentTci, isOutdated := w.wavetrendProvider.GetCurrentTci(m1SvcName)
	if isOutdated {
		logger.Warn(ctx, "[shouldBuy] M1 currentTci is outdated")
		return false
	}
	if currentTci > c.WavetrendOversold {
		return false
	}

	currentDifWt, isOutdated := w.wavetrendProvider.GetCurrentDifWavetrend(m1SvcName)
	if isOutdated {
		logger.Warn(ctx, "[shouldBuy] M1 currentDifWt is outdated")
		return false
	}
	if currentDifWt <= 0 {
		return false
	}

	pastWtDat, isOutdated := w.wavetrendProvider.GetPastWaveTrendData(m1SvcName)
	if pastWtDat == nil { // get error. not expected for err here
		return false
	}
	if isOutdated {
		logger.Warn(ctx, "[shouldBuy] M1 pastWtDat is outdated")
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
