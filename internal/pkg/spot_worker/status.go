package spotworker

import "sync"

type status struct {
	mu              *sync.Mutex
	totalUnitBought int64
	lastBoughtAt    uint64 // FIXME: always reset to 0 after restart bot. in second.
}

func newStatus() *status {
	return &status{
		mu: &sync.Mutex{},
	}
}

func (s *status) updateTotalUnitBought(value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.totalUnitBought = value
}

func (s *status) loadTotalUnitBought() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.totalUnitBought
}

func (s *status) storeLastBoughtAt(value uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastBoughtAt = value
}

func (s *status) loadLastBoughtAt() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.lastBoughtAt
}
