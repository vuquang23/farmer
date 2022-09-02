package spotworker

import (
	"sync"

	"farmer/internal/pkg/entities"
)

type workerSetting struct {
	mu              *sync.Mutex
	symbol          string // eg: BTCUSDT, ETHUSDT,...
	buyCountAllowed uint64
	buyCount        int64
	buyNotional     float64
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
	s.buyCountAllowed = e.UnitBuyAllowed
	s.buyCount = int64(e.TotalUnitBought)
	s.buyNotional = e.UnitNotional
}

func (s *workerSetting) loadSymbol() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.symbol
}

func (s *workerSetting) loadBuyCountAllowed() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buyCountAllowed
}

func (s *workerSetting) loadBuyNotional() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buyNotional
}

func (s *workerSetting) loadBuyCount() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buyCount
}

func (s *workerSetting) updateBuyCount(value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buyCount += value
}
