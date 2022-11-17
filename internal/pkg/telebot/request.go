package telebot

import (
	"strings"

	"farmer/internal/pkg/entities"
)

type CreateNewSpotWorkerReq struct {
	Symbol         string
	UnitBuyAllowed uint64
	UnitNotional   float64
}

func (r *CreateNewSpotWorkerReq) ToCreateNewSpotWorkerParams() *entities.CreateNewSpotWorkerParams {
	return &entities.CreateNewSpotWorkerParams{
		Symbol:         r.Symbol,
		UnitBuyAllowed: r.UnitBuyAllowed,
		UnitNotional:   r.UnitNotional,
	}
}

func (r *CreateNewSpotWorkerReq) Normalize() *CreateNewSpotWorkerReq {
	r.Symbol = strings.ToUpper(r.Symbol)
	return r
}
