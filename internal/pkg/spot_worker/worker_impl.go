package spotworker

import (
	"farmer/internal/pkg/entities"
	"sync"
)

type spotWorker struct {
	symbol      string
	exchangeInf *exchangeInfo
}

func NewSpotWorker(symbol string) ISpotWorker {
	return &spotWorker{
		symbol: symbol,
		exchangeInf: &exchangeInfo{
			mu: &sync.Mutex{},
		},
	}
}

func (w *spotWorker) SetExchangeInfo(info entities.ExchangeInfo) {
	w.exchangeInf.set(info)
}
