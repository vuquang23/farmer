package spotworker

import "farmer/internal/pkg/entities"

type ISpotWorker interface {
	SetExchangeInfo(info entities.ExchangeInfo)
}
