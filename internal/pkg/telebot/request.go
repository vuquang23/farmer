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

type StopWorkerReq struct {
	Symbol string
}

func (r *StopWorkerReq) Normalize() *StopWorkerReq {
	r.Symbol = strings.ToUpper(r.Symbol)
	return r
}

func (r *StopWorkerReq) ToStopBotParams() *entities.StopWorkerParams {
	return &entities.StopWorkerParams{
		Symbol: r.Symbol,
	}
}

type AddCapitalReq struct {
	Symbol  string
	Capital float64
}

func (r *AddCapitalReq) Normalize() *AddCapitalReq {
	r.Symbol = strings.ToUpper(r.Symbol)
	return r
}

func (r *AddCapitalReq) ToAddCapitalParams() *entities.AddCapitalParams {
	return &entities.AddCapitalParams{
		Symbol:  r.Symbol,
		Capital: r.Capital,
	}
}

type ArchiveSpotTradingDataReq struct {
	Symbol string
}

func (r *ArchiveSpotTradingDataReq) Normalize() *ArchiveSpotTradingDataReq {
	r.Symbol = strings.ToUpper(r.Symbol)
	return r
}

func (r *ArchiveSpotTradingDataReq) ToArchiveSpotTradingDataParams() *entities.ArchiveSpotTradingDataParams {
	return &entities.ArchiveSpotTradingDataParams{
		Symbol: r.Symbol,
	}
}
