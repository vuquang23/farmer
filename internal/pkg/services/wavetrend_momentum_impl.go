package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/gin-gonic/gin"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/enum"
	"farmer/internal/pkg/utils/indicators"
	"farmer/internal/pkg/utils/logger"
	"farmer/pkg/errors"
)

type wtMomentumService struct {
	spot        *binance.Client
	future      *futures.Client
	mapMomentum map[string]float64
	mu          *sync.Mutex
	n1          uint64
	n2          uint64
}

var wtMomentumSvc *wtMomentumService

func InitWaveTrendMomentumService(spot *binance.Client, future *futures.Client) {
	if wtMomentumSvc == nil {
		wtMomentumSvc = &wtMomentumService{
			spot:        spot,
			future:      future,
			mapMomentum: make(map[string]float64),
			mu:          &sync.Mutex{},
			n1:          10,
			n2:          21,
		}
	}
}

func WaveTrendMomentumServiceInstance() IWavetrendMomentumService {
	return wtMomentumSvc
}

func (s *wtMomentumService) Calculate(
	ctx *gin.Context, market enum.Market, symbolList []string, interval string,
) ([]*entities.WavetrendMomentum, *errors.DomainError) {
	batch := 30
	ret := []*entities.WavetrendMomentum{}

	for i := 0; i < len(symbolList); i += batch {
		wg := &sync.WaitGroup{}

		r := int(math.Min(float64(len(symbolList)), float64(i+batch)))
		for j := i; j < r; j++ {
			wg.Add(1)
			go s.calcForSymbol(ctx, wg, market, symbolList[j], interval)
		}

		wg.Wait()
		time.Sleep(time.Second * time.Duration(2)) // sleep 2s
	}

	for _, sym := range symbolList {
		ret = append(ret, &entities.WavetrendMomentum{
			Symbol: sym,
			Value:  s.mapMomentum[sym],
		})
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Value > ret[j].Value
	})
	return ret, nil
}

func (s *wtMomentumService) setMap(symbol string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mapMomentum[symbol] = value
}

func (s *wtMomentumService) calcForSymbol(
	ctx *gin.Context, wg *sync.WaitGroup, market enum.Market, symbol string, interval string,
) {
	defer wg.Done()

	candles, err := s.queryKline(ctx, market, symbol, interval)
	if err != nil {
		logger.FromGinCtx(ctx).Sugar().Error(err)
		return
	}

	momentum, _ := indicators.WaveTrendMomentumValue(
		candles, s.n1, s.n2,
	)
	s.setMap(symbol, momentum)
}

func (s *wtMomentumService) queryKline(
	ctx *gin.Context, market enum.Market, symbol string, interval string,
) ([]*entities.MinimalKline, *errors.DomainError) {
	pass := uint64(600) // 600 candles til now
	switch market {
	case enum.Spot:
		candles, err := s.spot.NewKlinesService().
			Symbol(symbol + "USDT").
			Interval(interval).
			Limit(int(pass)).
			Do(context.Background())
		if err != nil {
			return nil, errors.NewDomainErrorUnknown(err)
		}
		return indicators.SpotKlineToMinimalKline(candles), nil
	case enum.Future:
		candles, err := s.future.NewKlinesService().
			Symbol((symbol + "USDT")).
			Interval(interval).
			Limit(int(pass)).
			Do(context.Background())
		if err != nil {
			return nil, errors.NewDomainErrorUnknown(err)
		}
		return indicators.FutureKlineToMinimalKline(candles), nil
	default:
		return nil, errors.NewDomainErrorUnknown(fmt.Errorf("unknown market enum: %s", market))
	}
}
