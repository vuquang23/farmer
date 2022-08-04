package errors

import (
	"encoding/json"
	errs "errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
)

type IDomainTransformer interface {
	InfraErrToDomainErr(infraErr *InfraError) *DomainError
	RegisterTransformFunc(infraErrCode string, transformFunc domainTransformFunc)
}

type domainTransformFunc func(rootCause error, entities ...string) *DomainError

type domainTransformer struct {
	mapping map[string]domainTransformFunc
}

var domainTransformerInstance *domainTransformer

func initDomainTransformerInstance() {
	if domainTransformerInstance == nil {
		domainTransformerInstance = &domainTransformer{}
		domainTransformerInstance.mapping = map[string]domainTransformFunc{}
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeDBConnect, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeDBNotFound, NewDomainErrorNotFound)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeDBSelect, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeDBInsert, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeDBUpdate, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeDBDelete, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeDBUnknown, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeRedisConnect, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeRedisNotFound, NewDomainErrorNotFound)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeRedisGet, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeRedisSet, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeRedisUnknown, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeHTTPUnknown, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeRPCUnknown, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeElsConnect, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeElsUnknown, NewDomainErrorUnknown)
		domainTransformerInstance.RegisterTransformFunc(InfraErrCodeElsNotFound, NewDomainErrorNotFound)
	}
}

func DomainTransformerInstance() IDomainTransformer {
	return domainTransformerInstance
}

// InfraErrToDomainErr transforms InfraError to DomainError
func (t *domainTransformer) InfraErrToDomainErr(infraError *InfraError) *DomainError {
	f := t.mapping[infraError.Code]
	if f == nil {
		return NewDomainErrorUnknown(fmt.Errorf("can not transform error, InfraError: %v", infraError))
	}
	return f(infraError, infraError.ErrorEntities...)
}

// RegisterTransformFunc is used to add new function to transform InternalError to DomainError
// if the infraErrCode is already registered, the old transform function will be overridden
func (t *domainTransformer) RegisterTransformFunc(infraErrCode string, function domainTransformFunc) {
	t.mapping[infraErrCode] = function
}

type IRestTransformer interface {
	DomainErrToRestAPIErr(domainErr *DomainError) *RestAPIError
	ValidationErrToRestAPIErr(err error) *RestAPIError
	RegisterTransformFunc(domainErrCode string, transformFunc restTransformFunc)
}

type restTransformFunc func(rootCause error, entities ...string) *RestAPIError

type restTransformer struct {
	mapping map[string]restTransformFunc
}

var restTransformerInstance *restTransformer

func initRestTransformerInstance() {
	if restTransformerInstance == nil {
		restTransformerInstance = &restTransformer{}
		restTransformerInstance.mapping = map[string]restTransformFunc{}
		restTransformerInstance.RegisterTransformFunc(DomainErrCodeRequired, NewRestAPIErrRequired)
		restTransformerInstance.RegisterTransformFunc(DomainErrCodeNotAcceptedValue, NewRestAPIErrNotAcceptedValue)
		restTransformerInstance.RegisterTransformFunc(DomainErrCodeOutOfRange, NewRestAPIErrOutOfRange)
		restTransformerInstance.RegisterTransformFunc(DomainErrCodeInvalidFormat, NewRestAPIErrInvalidFormat)
		restTransformerInstance.RegisterTransformFunc(DomainErrCodeInvalid, NewRestAPIErrInvalid)
		restTransformerInstance.RegisterTransformFunc(DomainErrCodeUnauthenticated, NewRestAPIErrUnauthenticated)
		restTransformerInstance.RegisterTransformFunc(DomainErrCodeNotFound, NewRestAPIErrNotFound)
		restTransformerInstance.RegisterTransformFunc(DomainErrCodeDuplicate, NewRestAPIErrDuplicate)
		restTransformerInstance.RegisterTransformFunc(DomainErrCodeAlreadyExists, NewRestAPIErrAlreadyExits)
		restTransformerInstance.RegisterTransformFunc(DomainErrCodeUnknown, NewRestAPIErrInternal)
	}
}

func RestTransformerInstance() IRestTransformer {
	return restTransformerInstance
}

// ValidationErrToRestAPIErr transforms ValidationError to RestAPIError
// this function will be used when bind JSON request to DTO in gin framework
func (t *restTransformer) ValidationErrToRestAPIErr(err error) *RestAPIError {
	var validationErrs validator.ValidationErrors
	var unmarshalTypeErr *json.UnmarshalTypeError
	var jsonSynTaxErr *json.SyntaxError
	var numErr *strconv.NumError
	if errs.As(err, &validationErrs) {
		validationErr := validationErrs[0]
		return apiErrForTag(validationErr.Tag(), err, validationErr.Field())
	}
	if errs.As(err, &unmarshalTypeErr) {
		field := unmarshalTypeErr.Field
		fieldArr := strings.Split(field, ".")
		return NewRestAPIErrInvalidFormat(err, fieldArr[len(fieldArr)-1])
	}
	if errs.As(err, &jsonSynTaxErr) {
		return NewRestAPIErrInvalidFormat(err)
	}
	if errs.As(err, &numErr) {
		return NewRestAPIErrInvalidFormat(err)
	}
	return NewRestAPIErrInternal(err)
}

// DomainErrToRestAPIErr transforms DomainError to RestAPIError
func (t *restTransformer) DomainErrToRestAPIErr(domainErr *DomainError) *RestAPIError {
	f := t.mapping[domainErr.Code]
	if f == nil {
		return NewRestAPIErrInternal(fmt.Errorf("can not transform error, DomainError: %v", domainErr))
	}
	return f(domainErr, domainErr.ErrorEntities...)
}

// RegisterTransformFunc is used to add new function to transform DomainError to RestAPIError
// if the domainErrCode is already registered, the old transform function will be overridden
func (t *restTransformer) RegisterTransformFunc(domainErrCode string, function restTransformFunc) {
	t.mapping[domainErrCode] = function
}

// apiErrForTag return RestAPIError which corresponds to the validation tag
func apiErrForTag(tag string, err error, fields ...string) *RestAPIError {
	switch tag {
	case "required":
		return NewRestAPIErrRequired(err, fields...)
	case "oneof":
		return NewRestAPIErrNotAcceptedValue(err, fields...)
	case "min", "max":
		return NewRestAPIErrOutOfRange(err, fields...)
	default:
		return NewRestAPIErrInternal(err)
	}
}
