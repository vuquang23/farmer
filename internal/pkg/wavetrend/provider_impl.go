package wavetrendprovider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/adshao/go-binance/v2"

	b "farmer/internal/pkg/binance"
	"farmer/internal/pkg/entities"
	e "farmer/internal/pkg/errors"
	"farmer/internal/pkg/utils/logger"
	w "farmer/internal/pkg/wavetrend/worker"
	errPkg "farmer/pkg/errors"
)

type wavetrendProvider struct {
	mapSymbolWorker     map[string]w.IWavetrendWorker
	mapSymbolStopWsChan map[string]chan<- struct{}

	klineChannel *gochannel.GoChannel
}

var provider *wavetrendProvider

func InitWavetrendProvider() {
	if provider == nil {
		provider = &wavetrendProvider{
			mapSymbolWorker:     make(map[string]w.IWavetrendWorker),
			mapSymbolStopWsChan: make(map[string]chan<- struct{}),
			klineChannel: gochannel.NewGoChannel(gochannel.Config{
				OutputChannelBuffer:            50,
				Persistent:                     false,
				BlockPublishUntilSubscriberAck: false,
			}, nil),
		}
	}
}

func WavetrendProviderInstance() IWavetrendProvider {
	return provider
}

// StartService start connect websocket to binance server and start worker.
//
// svcName: BTCUSDT:1h, ETHUSDT:1m, future:ETHUSDT:1m...
func (p *wavetrendProvider) StartService(svcName string) *errPkg.DomainError {
	if _, ok := p.mapSymbolWorker[svcName]; ok {
		return e.NewDomainErrorWavetrendServiceNameExisted(nil)
	}

	// init ws to receive realtime data from binance and push data to wavetrend worker
	initC := make(chan error)
	stopConnC := make(chan struct{})
	go p.startKlineWSConnection(svcName, initC, stopConnC)
	if err := <-initC; err != nil {
		return errPkg.NewDomainErrorUnknown(err)
	}
	p.mapSymbolStopWsChan[svcName] = stopConnC

	// wavetrend worker subscribe to receive kline data from wavetrend provider
	ctx, cancel := context.WithCancel(context.Background())
	klineMsgChan, err := p.klineChannel.Subscribe(ctx, svcName)
	if err != nil {
		cancel()
		return errPkg.NewDomainErrorUnknown(err)
	}

	worker := w.NewWavetrendWorker(svcName, b.BinanceSpotClientInstance(), klineMsgChan, cancel)
	start := make(chan error)
	go worker.Run(start)
	if err := <-start; err != nil {
		return e.NewDomainErrorWavetrendServiceNameExisted(err)
	}
	p.mapSymbolWorker[svcName] = worker

	return nil
}

func (p *wavetrendProvider) startKlineWSConnection(svcName string, initC chan<- error, stopConnC chan struct{}) {
	log := logger.WithDescription(fmt.Sprintf("%s - Start WS Connection", svcName))
	strs := strings.Split(svcName, ":")

	// TODO: case future.
	if len(strs) > 2 {
		initC <- errors.New("future is not supported now")
		return
	}

	symbol := strs[0]
	timeFrame := strs[1]

	var handler = func(event *binance.WsKlineEvent) {
		wsKline := event.Kline
		kline := &binance.Kline{
			OpenTime:                 wsKline.StartTime,
			Open:                     wsKline.Open,
			High:                     wsKline.High,
			Low:                      wsKline.Low,
			Close:                    wsKline.Close,
			Volume:                   wsKline.Volume,
			CloseTime:                wsKline.EndTime,
			QuoteAssetVolume:         wsKline.QuoteVolume,
			TradeNum:                 wsKline.TradeNum,
			TakerBuyBaseAssetVolume:  wsKline.ActiveBuyVolume,
			TakerBuyQuoteAssetVolume: wsKline.ActiveBuyQuoteVolume,
		}

		marshedPayload, err := json.Marshal(kline)
		if err != nil {
			log.Sugar().Error(err)
			return
		}

		if err := p.klineChannel.Publish(svcName, &message.Message{
			Payload: marshedPayload,
		}); err != nil {
			log.Sugar().Error(err)
		}
	}

	resetC := make(chan struct{})
	var errHandler = func(err error) {
		log.Sugar().Error(err)

		resetC <- struct{}{}
	}

	once := &sync.Once{}
	for {
		_, stopC, err := binance.WsKlineServe(symbol, timeFrame, handler, errHandler)
		if err != nil {
			log.Sugar().Error()
			stopC <- struct{}{}
			continue
		}

		once.Do(func() {
			initC <- nil
		})

		// polling
		select {
		case <-stopConnC:
			return
		case <-resetC:
			stopC <- struct{}{}
		}
	}
}

func (p *wavetrendProvider) SetStopSignal(svcName string) {
	w, ok := p.mapSymbolWorker[svcName]
	if ok {
		// stop worker
		w.Stop()

		// unscribe to ws
		stopConnC := p.mapSymbolStopWsChan[svcName]
		stopConnC <- struct{}{}
	}
}

func (p *wavetrendProvider) GetCurrentTci(svcName string) (float64, bool) {
	w, ok := p.mapSymbolWorker[svcName]
	if ok {
		return w.GetCurrentTci()
	}
	return 0, true
}

func (p *wavetrendProvider) GetCurrentDifWavetrend(svcName string) (float64, bool) {
	w, ok := p.mapSymbolWorker[svcName]
	if ok {
		return w.GetCurrentDifWavetrend()
	}
	return 0, true
}

func (p *wavetrendProvider) GetClosePrice(svcName string) float64 {
	w, ok := p.mapSymbolWorker[svcName]
	if ok {
		return w.GetClosePrice()
	}
	return 0
}

func (p *wavetrendProvider) GetPastWaveTrendData(svcName string) (*entities.PastWavetrend, bool) {
	w, ok := p.mapSymbolWorker[svcName]
	if ok {
		return w.GetPastWaveTrendData()
	}
	return nil, false
}
