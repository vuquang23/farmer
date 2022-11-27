package spotworker

import (
	"sync"

	"farmer/internal/pkg/entities"
)

type workerSetting struct {
	mu             *sync.Mutex
	symbol         string // eg: BTCUSDT, ETHUSDT,...
	unitBuyAllowed uint64
	unitNotional   float64
}

func newWorkerSetting() *workerSetting {
	return &workerSetting{
		mu: &sync.Mutex{},
	}
}

func (s *workerSetting) store(e entities.SpotWorkerStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.symbol = e.Symbol
	s.unitBuyAllowed = e.UnitBuyAllowed
	s.unitNotional = e.UnitNotional
}

func (s *workerSetting) loadUnitBuyAllowed() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.unitBuyAllowed
}

func (s *workerSetting) loadUnitNotional() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.unitNotional
}

func (s *workerSetting) updateUnitNotional(val float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.unitNotional += val
}

func (s *workerSetting) addCapital(val float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.unitNotional += val / float64(s.unitBuyAllowed)
}
