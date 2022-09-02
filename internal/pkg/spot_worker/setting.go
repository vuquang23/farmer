package spotworker

import (
	"sync"

	"farmer/internal/pkg/entities"
)

type workerSetting struct {
	mu              *sync.Mutex
	symbol          string // eg: BTCUSDT, ETHUSDT,...
	unitBuyAllowed  uint64
	totalUnitBought int64
	unitNotional    float64
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
	s.totalUnitBought = int64(e.TotalUnitBought)
	s.unitNotional = e.UnitNotional
}

func (s *workerSetting) loadSymbol() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.symbol
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

func (s *workerSetting) loadTotalUnitBought() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.totalUnitBought
}

func (s *workerSetting) updateTotalUnitBought(value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.totalUnitBought += value
}
