package errors

import pkgErr "farmer/pkg/errors"

func NewDomainErrorInitExchangeInfo(rootCause error, entities ...string) *pkgErr.DomainError {
	return pkgErr.NewDomainError(DomainErrCodeInitExchangeInfo, DomainErrMsgInitExchangeInfo, entities, rootCause)
}

func NewDomainErrorInitWavetrendData(rootCause error, entities ...string) *pkgErr.DomainError {
	return pkgErr.NewDomainError(DomainErrCodeInitWavetrendData, DomainErrMsgInitWavetrendData, entities, rootCause)
}

func NewDomainErrorWavetrendServiceNameExisted(rootCause error, entities ...string) *pkgErr.DomainError {
	return pkgErr.NewDomainError(DomainErrCodeWavetrendServiceNameExisted, DomainErrMsgWavetrendServiceNameExisted, entities, rootCause)
}

func NewDomainErrorCreateBuyOrderFailed(rootCause error, entities ...string) *pkgErr.DomainError {
	return pkgErr.NewDomainError(DomainErrCodeCreateBuyOrderFailed, DomainErrMsgCreateBuyOrderFailed, entities, rootCause)
}

func NewDomainErrorCreateSellOrderFailed(rootCause error, entities ...string) *pkgErr.DomainError {
	return pkgErr.NewDomainError(DomainErrCodeCreateSellOrderFailed, DomainErrMsgCreateSellOrderFailed, entities, rootCause)
}

func NewDomainErrorGetPriceFailed(rootCause error, entities ...string) *pkgErr.DomainError {
	return pkgErr.NewDomainError(DomainErrCodeGetPriceFailed, DomainErrMsgGetPriceFailed, entities, rootCause)
}
