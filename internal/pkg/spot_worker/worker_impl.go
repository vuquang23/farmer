package spotworker

import (
	"farmer/internal/pkg/entities"
	"sync"
)

type spotWorker struct {
	exchangeInf *exchangeInfo
	setting     *workerSetting
}

func NewSpotWorker() ISpotWorker {
	return &spotWorker{
		exchangeInf: &exchangeInfo{
			mu: &sync.Mutex{},
		},
		setting: &workerSetting{
			mu: &sync.Mutex{},
		},
	}
}

func (w *spotWorker) SetExchangeInfo(info entities.ExchangeInfo) error {
	w.exchangeInf.set(info)
	return nil
}

func (w *spotWorker) SetWorkerSetting(setting entities.SpotWorker) error {
	w.setting.set(setting)
	return nil
}

func (w *spotWorker) Run() error {
	
	return nil
}
