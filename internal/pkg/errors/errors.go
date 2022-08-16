package errors

import pkgErr "farmer/pkg/errors"

func NewDomainErrorInitExchangeInfo(rootCause error, entities ...string) *pkgErr.DomainError {
	return pkgErr.NewDomainError(DomainErrCodeInitExchangeInfo, DomainErrMsgInitExchangeInfo, entities, rootCause)
}

func NewDomainErrorInitWavetrendData(rootCause error, entities ...string) *pkgErr.DomainError {
	return pkgErr.NewDomainError(DomainErrCodeInitWavetrendData, DomainErrMsgInitWavetrendData, entities, rootCause)
}
