package spotworker

import (
	"farmer/internal/pkg/entities"
	"sync"
)

type exchangeInfo struct {
	mu             *sync.Mutex
	pricePrecision int
	qtyPrecision   int
	minQty         float64
	minNotional    float64
}

func (e *exchangeInfo) set(info entities.ExchangeInfo) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.pricePrecision = info.PricePrecision
	e.qtyPrecision = info.QtyPrecision
	e.minQty = info.MinQty
	e.minNotional = info.MinNotional
}

func (e *exchangeInfo) getPricePrecision() int {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.pricePrecision
}

func (e *exchangeInfo) getQtyPrecision() int {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.qtyPrecision
}

func (e *exchangeInfo) getMinQty() float64 {
	e.mu.Lock()
	defer e.mu.Lock()

	return e.minQty
}

func (e *exchangeInfo) getMinNotional() float64 {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.minQty
}
