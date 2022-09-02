package wavetrendprovider

import (
	"farmer/internal/pkg/binance"
	"farmer/internal/pkg/entities"
	e "farmer/internal/pkg/errors"
	w "farmer/internal/pkg/wavetrend/worker"
	errPkg "farmer/pkg/errors"
)

type wavetrendProvider struct {
	mapping map[string]w.IWavetrendWorker
}

var provider *wavetrendProvider

func InitWavetrendProvider() {
	if provider == nil {
		provider = &wavetrendProvider{
			mapping: make(map[string]w.IWavetrendWorker),
		}
	}
}

func WavetrendProviderInstance() IWavetrendProvider {
	return provider
}

func (p *wavetrendProvider) StartService(svcName string) *errPkg.DomainError {
	if _, ok := p.mapping[svcName]; ok {
		return e.NewDomainErrorWavetrendServiceNameExisted(nil)
	}

	w := w.NewWavetrendWorker(svcName, binance.BinanceSpotClientInstance())

	start := make(chan error)
	go w.Run(start)
	if err := <-start; err != nil {
		return e.NewDomainErrorWavetrendServiceNameExisted(err)
	}

	p.mapping[svcName] = w

	return nil
}

func (p *wavetrendProvider) SetStopSignal(svcName string) {
	w, ok := p.mapping[svcName]
	if ok {
		w.SetStopSignal()
	}
}

func (p *wavetrendProvider) GetCurrentTci(svcName string) float64 {
	w, ok := p.mapping[svcName]
	if ok {
		return w.GetCurrentTci()
	}
	return 0
}

func (p *wavetrendProvider) GetCurrentDifWavetrend(svcName string) float64 {
	w, ok := p.mapping[svcName]
	if ok {
		return w.GetCurrentDifWavetrend()
	}
	return 0
}

func (p *wavetrendProvider) GetClosePrice(svcName string) float64 {
	w, ok := p.mapping[svcName]
	if ok {
		return w.GetClosePrice()
	}
	return 0
}

func (p *wavetrendProvider) GetPastWaveTrendData(svcName string) *entities.PastWavetrend {
	w, ok := p.mapping[svcName]
	if ok {
		return w.GetPastWaveTrendData()
	}
	return nil
}
