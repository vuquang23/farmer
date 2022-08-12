package services

import (
	"github.com/gin-gonic/gin"

	"farmer/internal/pkg/entities"
	"farmer/pkg/errors"
)

type IWavetrendMomentumService interface {
	Calculate(ctx *gin.Context, symbolList []string, interval string) ([]*entities.WavetrendMomentum, *errors.DomainError)
}
