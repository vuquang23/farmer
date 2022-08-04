package errors

import (
	"fmt"
	"net/http"
	"strings"
)

type InfraError struct {
	Code          string
	Message       string
	ErrorEntities []string
	RootCause     error
}

func NewInfraError(code string, message string, entities []string, rootCause error) *InfraError {
	return &InfraError{
		Code:          code,
		Message:       message,
		ErrorEntities: entities,
		RootCause:     rootCause,
	}
}

func (e *InfraError) Error() string {
	return fmt.Sprintf("INFRA ERROR: {Code: %s, Message: %s, ErrorEntities: %v, RootCause: %v}", e.Code, e.Message, e.ErrorEntities, e.RootCause)
}

func NewInfraErrorDBConnect(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeDBConnect, InfraErrMsgDBConnect, entities, rootCause)
}

func NewInfraErrorDBNotFound(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeDBNotFound, InfraErrMsgDBNotFound, entities, rootCause)
}

func NewInfraErrorDBSelect(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeDBSelect, InfraErrMsgDBSelect, entities, rootCause)
}

func NewInfraErrorDBInsert(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeDBInsert, InfraErrMsgDBInsert, entities, rootCause)
}

func NewInfraErrorDBUpdate(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeDBUpdate, InfraErrMsgDBUpdate, entities, rootCause)
}

func NewInfraErrorDBDelete(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeDBDelete, InfraErrMsgDBDelete, entities, rootCause)
}

func NewInfraErrorDBUnknown(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeDBUnknown, InfraErrMsgDBUnknown, entities, rootCause)
}

func NewInfraErrorRedisConnect(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeRedisConnect, InfraErrMsgRedisConnect, entities, rootCause)
}

func NewInfraErrorRedisNotFound(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeRedisNotFound, InfraErrMsgRedisNotFound, entities, rootCause)
}

func NewInfraErrorRedisGet(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeRedisGet, InfraErrMsgRedisGet, entities, rootCause)
}

func NewInfraErrorRedisSet(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeRedisSet, InfraErrMsgRedisSet, entities, rootCause)
}

func NewInfraErrorRedisUnknown(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeRedisUnknown, InfraErrMsgRedisUnknown, entities, rootCause)
}

func NewInfraErrorHTTPUnknown(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeHTTPUnknown, InfraErrMsgHTTPUnknown, entities, rootCause)
}

func NewInfraErrorRPCUnknown(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeRPCUnknown, InfraErrMsgRPCUnknown, entities, rootCause)
}

func NewInfraErrorElsConnect(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeElsConnect, InfraErrMsgElsConnect, entities, rootCause)
}

func NewInfraErrorElsUnknown(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeElsUnknown, InfraErrMsgElsUnknown, entities, rootCause)
}

func NewInfraErrorElsNotFound(rootCause error, entities ...string) *InfraError {
	return NewInfraError(InfraErrCodeElsNotFound, InfraErrMsgElsNotFound, entities, rootCause)
}

type DomainError struct {
	Code          string
	Message       string
	ErrorEntities []string
	RootCause     error
}

func NewDomainError(code string, message string, entities []string, rootCause error) *DomainError {
	return &DomainError{
		Code:          code,
		Message:       message,
		ErrorEntities: entities,
		RootCause:     rootCause,
	}
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("DOMAIN ERROR: {Code: %s, Message: %s, ErrorEntities: %v, RootCause: %v}", e.Code, e.Message, e.ErrorEntities, e.RootCause)
}

func NewDomainErrorRequired(rootCause error, entities ...string) *DomainError {
	return NewDomainError(DomainErrCodeRequired, DomainErrMsgRequired, entities, rootCause)
}

func NewDomainErrorInvalidFormat(rootCause error, entities ...string) *DomainError {
	return NewDomainError(DomainErrCodeInvalidFormat, DomainErrMsgInvalidFormat, entities, rootCause)
}

func NewDomainErrorInvalid(rootCause error, entities ...string) *DomainError {
	return NewDomainError(DomainErrCodeInvalid, DomainErrMsgInvalid, entities, rootCause)
}

func NewDomainErrorNotAcceptedValue(rootCause error, entities ...string) *DomainError {
	return NewDomainError(DomainErrCodeNotAcceptedValue, DomainErrMsgNotAcceptedValue, entities, rootCause)
}

func NewDomainErrorOutOfRange(rootCause error, entities ...string) *DomainError {
	return NewDomainError(DomainErrCodeOutOfRange, DomainErrMsgOutOfRange, entities, rootCause)
}

func NewDomainErrorUnauthenticated(rootCause error, entities ...string) *DomainError {
	return NewDomainError(DomainErrCodeUnauthenticated, DomainErrMsgUnauthenticated, entities, rootCause)
}

func NewDomainErrorNotFound(rootCause error, entities ...string) *DomainError {
	return NewDomainError(DomainErrCodeNotFound, DomainErrMsgNotFound, entities, rootCause)
}

func NewDomainErrorDuplicate(rootCause error, entities ...string) *DomainError {
	return NewDomainError(DomainErrCodeDuplicate, DomainErrMsgDuplicate, entities, rootCause)
}

func NewDomainErrorAlreadyExits(rootCause error, entities ...string) *DomainError {
	return NewDomainError(DomainErrCodeAlreadyExists, DomainErrMsgAlreadyExists, entities, rootCause)
}

func NewDomainErrorUnknown(rootCause error, entities ...string) *DomainError {
	return NewDomainError(DomainErrCodeUnknown, DomainErrMsgUnknown, entities, rootCause)
}

type RestAPIError struct {
	HttpStatus    int           `json:"-"`
	Code          int           `json:"code"`
	Message       string        `json:"message"`
	ErrorEntities []string      `json:"errorEntities"`
	Details       []interface{} `json:"details"`
	RootCause     error         `json:"-"`
}

func NewRestAPIError(httpStatus int, code int, message string, entities []string, rootCause error) *RestAPIError {
	return &RestAPIError{
		HttpStatus:    httpStatus,
		Code:          code,
		Message:       message,
		ErrorEntities: entities,
		RootCause:     rootCause,
	}
}

func (e *RestAPIError) Error() string {
	return fmt.Sprintf("API ERROR: {Code: %d, Message: %s, ErrorEntities: %v, RootCause: %v}", e.Code, e.Message, e.ErrorEntities, e.RootCause)
}

func AppendEntitiesToErrMsg(message string, entities []string) string {
	if len(entities) > 0 {
		message += ": "
		message += strings.Join(entities, ",")
	}
	return message
}

func NewRestAPIErrRequired(rootCause error, entities ...string) *RestAPIError {
	message := AppendEntitiesToErrMsg(ClientErrMsgRequired, entities)
	return NewRestAPIError(http.StatusBadRequest, ClientErrCodeRequired, message, entities, rootCause)
}

func NewRestAPIErrInvalidFormat(rootCause error, entities ...string) *RestAPIError {
	message := AppendEntitiesToErrMsg(ClientErrMsgInvalidFormat, entities)
	return NewRestAPIError(http.StatusBadRequest, ClientErrCodeInvalidFormat, message, entities, rootCause)
}

func NewRestAPIErrInvalid(rootCause error, entities ...string) *RestAPIError {
	message := AppendEntitiesToErrMsg(ClientErrMsgInvalid, entities)
	return NewRestAPIError(http.StatusBadRequest, ClientErrCodeInvalid, message, entities, rootCause)
}

func NewRestAPIErrNotAcceptedValue(rootCause error, entities ...string) *RestAPIError {
	message := AppendEntitiesToErrMsg(ClientErrMsgNotAcceptedValue, entities)
	return NewRestAPIError(http.StatusBadRequest, ClientErrCodeNotAcceptedValue, message, entities, rootCause)
}

func NewRestAPIErrOutOfRange(rootCause error, entities ...string) *RestAPIError {
	message := AppendEntitiesToErrMsg(ClientErrMsgOutOfRange, entities)
	return NewRestAPIError(http.StatusBadRequest, ClientErrCodeOutOfRange, message, entities, rootCause)
}

func NewRestAPIErrUnauthenticated(rootCause error, entities ...string) *RestAPIError {
	message := AppendEntitiesToErrMsg(ClientErrMsgUnauthenticated, entities)
	return NewRestAPIError(http.StatusUnauthorized, ClientErrCodeUnauthenticated, message, entities, rootCause)
}

func NewRestAPIErrNotFound(rootCause error, entities ...string) *RestAPIError {
	message := AppendEntitiesToErrMsg(ClientErrMsgNotFound, entities)
	return NewRestAPIError(http.StatusNotFound, ClientErrCodeNotFound, message, entities, rootCause)
}

func NewRestAPIErrDuplicate(rootCause error, entities ...string) *RestAPIError {
	message := AppendEntitiesToErrMsg(ClientErrMsgDuplicate, entities)
	return NewRestAPIError(http.StatusConflict, ClientErrCodeDuplicate, message, entities, rootCause)
}

func NewRestAPIErrAlreadyExits(rootCause error, entities ...string) *RestAPIError {
	message := AppendEntitiesToErrMsg(ClientErrMsgAlreadyExists, entities)
	return NewRestAPIError(http.StatusConflict, ClientErrCodeAlreadyExists, message, entities, rootCause)
}

func NewRestAPIErrInternal(rootCause error, entities ...string) *RestAPIError {
	return NewRestAPIError(http.StatusInternalServerError, ClientErrCodeInternal, ClientErrMsgInternal, entities, rootCause)
}
