package wavetrendprovider

import (
	goctx "context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/adshao/go-binance/v2"

	b "farmer/internal/pkg/binance"
	"farmer/internal/pkg/entities"
	e "farmer/internal/pkg/errors"
	"farmer/internal/pkg/utils/context"
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
func (p *wavetrendProvider) StartService(ctx goctx.Context, svcName string) *errPkg.DomainError {
	if _, ok := p.mapSymbolWorker[svcName]; ok {
		return e.NewDomainErrorWavetrendServiceNameExisted(nil)
	}

	// init ws to receive realtime data from binance and push data to wavetrend worker
	initC := make(chan error)
	stopConnC := make(chan struct{})
	go p.startKlineWSConnection(context.Child(ctx, fmt.Sprintf("[kline-ws] %s", svcName)), svcName, initC, stopConnC)
	if err := <-initC; err != nil {
		return errPkg.NewDomainErrorUnknown(err)
	}
	p.mapSymbolStopWsChan[svcName] = stopConnC

	// wavetrend worker subscribe to receive kline data from wavetrend provider
	c, cancel := goctx.WithCancel(ctx)
	klineMsgChan, err := p.klineChannel.Subscribe(context.Child(c, fmt.Sprintf("[kline-channel] %s", svcName)), svcName)
	if err != nil {
		cancel()
		return errPkg.NewDomainErrorUnknown(err)
	}

	worker := w.NewWavetrendWorker(svcName, b.BinanceSpotClientInstance(), klineMsgChan, cancel)
	start := make(chan error)
	go worker.Run(context.Child(ctx, fmt.Sprintf("[wavetrend-worker] %s", svcName)), start)
	if err := <-start; err != nil {
		return errPkg.NewDomainErrorUnknown(err)
	}
	p.mapSymbolWorker[svcName] = worker

	return nil
}

func (p *wavetrendProvider) startKlineWSConnection(ctx goctx.Context, svcName string, initC chan<- error, stopConnC chan struct{}) {
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
			logger.Error(ctx, err)
			return
		}

		if err := p.klineChannel.Publish(svcName, &message.Message{
			Payload: marshedPayload,
		}); err != nil {
			logger.Error(ctx, err)
			return
		}
	}

	var errHandler = func(err error) {
		logger.Error(ctx, err)
	}

	once := &sync.Once{}
	for {
		logger.Info(ctx, "[startKlineWSConnection] connect WS Kline")

		doneC, stopC, err := binance.WsKlineServe(symbol, timeFrame, handler, errHandler)
		if err != nil {
			logger.Error(ctx, err)
			continue
		}

		once.Do(func() {
			initC <- nil
		})

		logger.Info(ctx, "[startKlineWSConnection] start polling...")
		// polling
		select {
		case <-stopConnC:
			logger.Info(ctx, "[startKlineWSConnection] in stopConnC...")
			stopC <- struct{}{}
			return
		case <-doneC:
			logger.Info(ctx, "[startKlineWSConnection] in doneC...")
			time.Sleep(2 * time.Second)
		}

		logger.Info(ctx, "[startKlineWSConnection] reset Kline WS connection")
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

func (p *wavetrendProvider) GetClosePrice(svcName string) (float64, bool) {
	w, ok := p.mapSymbolWorker[svcName]
	if ok {
		return w.GetClosePrice()
	}
	return 0, true
}

func (p *wavetrendProvider) GetPastWaveTrendData(svcName string) (*entities.PastWavetrend, bool) {
	w, ok := p.mapSymbolWorker[svcName]
	if ok {
		return w.GetPastWaveTrendData()
	}
	return nil, true
}
