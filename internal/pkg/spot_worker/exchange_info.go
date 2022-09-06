package spotworker

import (
	"sync"

	"farmer/internal/pkg/entities"
)

type exchangeInfo struct {
	mu             *sync.Mutex
	pricePrecision int
	qtyPrecision   int
	minQty         float64
	minNotional    float64
}

func newExchangeInfo() *exchangeInfo {
	return &exchangeInfo{
		mu: &sync.Mutex{},
	}
}

func (e *exchangeInfo) store(info entities.SpotExchangeInfo) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.pricePrecision = info.PricePrecision
	e.qtyPrecision = info.QtyPrecision
	e.minQty = info.MinQty
	e.minNotional = info.MinNotional
}

func (e *exchangeInfo) loadPricePrecision() int {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.pricePrecision
}

func (e *exchangeInfo) loadQtyPrecision() int {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.qtyPrecision
}

func (e *exchangeInfo) loadMinQty() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.minQty
}

func (e *exchangeInfo) loadMinNotional() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.minQty
}
