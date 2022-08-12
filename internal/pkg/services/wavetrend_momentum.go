package services

import (
	"github.com/gin-gonic/gin"

	"farmer/internal/pkg/entities"
	"farmer/internal/pkg/enum"
	"farmer/pkg/errors"
)

type IWavetrendMomentumService interface {
	Calculate(ctx *gin.Context, market enum.Market, symbolList []string, interval string) ([]*entities.WavetrendMomentum, *errors.DomainError)
}
