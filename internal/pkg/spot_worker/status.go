package spotworker

import (
	"sync"
	"time"
)

type status struct {
	mu              *sync.Mutex
	totalUnitBought int64
	lastBoughtAt    time.Time
	health          time.Time
}

func newStatus() *status {
	return &status{
		mu:           &sync.Mutex{},
		lastBoughtAt: time.Unix(0, 0),
		health:       time.Unix(0, 0),
	}
}

func (s *status) updateTotalUnitBought(value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.totalUnitBought += value
}

func (s *status) loadTotalUnitBought() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.totalUnitBought
}

func (s *status) storeLastBoughtAt(value time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastBoughtAt = value
}

func (s *status) loadLastBoughtAt() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.lastBoughtAt
}

func (s *status) storeHealth(value time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.health = value
}

func (s *status) loadHealth() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.health
}
