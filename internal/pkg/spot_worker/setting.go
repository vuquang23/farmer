package spotworker

import (
	"farmer/internal/pkg/entities"
	"sync"
)

type workerSetting struct {
	mu              *sync.Mutex
	symbol          string
	buyCountAllowed uint64
	buyCount        int64
	buyNotional     float64
}

func (s *workerSetting) set(e entities.SpotWorker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.symbol = e.Symbol
	s.buyCountAllowed = e.BuyCountAllowed
	s.buyCount = e.BuyCount
	s.buyNotional = e.BuyNotional
}

func (s *workerSetting) getSymbol() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.symbol
}

func (s *workerSetting) getBuyCountAllowed() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buyCountAllowed
}

func (s *workerSetting) getBuyNotional() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buyNotional
}

func (s *workerSetting) getBuyCount() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buyCount
}

func (s *workerSetting) updateBuyCount(value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buyCount += value
}
