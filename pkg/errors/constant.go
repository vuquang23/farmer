package errors

const (
	// HTTP 200 - OK
	ClientErrCodeOK = 0
	ClientErrMsgOK  = "Succeeded"

	//HTTP 400 - Bad Request
	ClientErrCodeRequired = 4000
	ClientErrMsgRequired  = "Missing required fields"

	ClientErrCodeNotAcceptedValue = 4001
	ClientErrMsgNotAcceptedValue  = "Input is not in the accepted values"

	ClientErrCodeOutOfRange = 4002
	ClientErrMsgOutOfRange  = "Input is out of range"

	ClientErrCodeInvalidFormat = 4003
	ClientErrMsgInvalidFormat  = "Input has an invalid format"

	ClientErrCodeInvalid = 4004
	ClientErrMsgInvalid  = "Input is invalid"

	//HTTP 401 - Unauthorized
	ClientErrCodeUnauthenticated = 4010
	ClientErrMsgUnauthenticated  = "Unauthenticated"

	//HTTP 404 - Not found
	ClientErrCodeNotFound = 4040
	ClientErrMsgNotFound  = "Not found"

	//HTTP 409 - Duplicate
	ClientErrCodeDuplicate = 4090
	ClientErrMsgDuplicate  = "Duplicate data"

	ClientErrCodeAlreadyExists = 4091
	ClientErrMsgAlreadyExists  = "Data already exists"

	//HTTP 500 - Internal Server Error
	ClientErrCodeInternal = 5000
	ClientErrMsgInternal  = "Internal Server Error"
)

const (
	DomainErrCodeRequired = "DOMAIN:REQUIRED"
	DomainErrMsgRequired  = "Domain error: Missing required fields"

	DomainErrCodeNotAcceptedValue = "DOMAIN:NOT_ACCEPTED_VALUE"
	DomainErrMsgNotAcceptedValue  = "Domain error: Input is not in the accepted values"

	DomainErrCodeOutOfRange = "DOMAIN:OUT_OF_RANGE"
	DomainErrMsgOutOfRange  = "Domain error: Input is out of range"

	DomainErrCodeInvalidFormat = "DOMAIN:INVALID_FORMAT"
	DomainErrMsgInvalidFormat  = "Domain error: Input has an invalid format"

	DomainErrCodeInvalid = "DOMAIN:INVALID"
	DomainErrMsgInvalid  = "Domain error: Input is invalid "

	DomainErrCodeUnauthenticated = "DOMAIN:UNAUTHENTICATED"
	DomainErrMsgUnauthenticated  = "Domain error: Unauthenticated"

	DomainErrCodeNotFound = "DOMAIN:NOT_FOUND"
	DomainErrMsgNotFound  = "Domain error: Not found"

	DomainErrCodeDuplicate = "DOMAIN:DUPLICATE"
	DomainErrMsgDuplicate  = "Domain error: Duplicate data"

	DomainErrCodeAlreadyExists = "DOMAIN:ALREADY_EXISTS"
	DomainErrMsgAlreadyExists  = "Domain error: Data already exists"

	DomainErrCodeUnknown = "DOMAIN:UNKNOWN"
	DomainErrMsgUnknown  = "Domain error: Unknown Domain Error"
)

const (
	InfraErrCodeDBConnect = "INFRA:DATABASE:CONNECT"
	InfraErrMsgDBConnect  = "Infra error: Failed to connect to database"

	InfraErrCodeDBNotFound = "INFRA:DATABASE:NOT_FOUND"
	InfraErrMsgDBNotFound  = "Infra error: Not found resource in database"

	InfraErrCodeDBSelect = "INFRA:DATABASE:SELECT"
	InfraErrMsgDBSelect  = "Infra error: Failed to select resources from database"

	InfraErrCodeDBInsert = "INFRA:DATABASE:INSERT"
	InfraErrMsgDBInsert  = "Infra error: Failed to insert into database"

	InfraErrCodeDBUpdate = "INFRA:DATABASE:UPDATE"
	InfraErrMsgDBUpdate  = "Infra error: Failed to update resources in database"

	InfraErrCodeDBDelete = "INFRA:DATABASE:DELETE"
	InfraErrMsgDBDelete  = "Infra error: Failed to delete resources from database"

	InfraErrCodeDBUnknown = "INFRA:DATABASE:UNKNOWN"
	InfraErrMsgDBUnknown  = "Infra error: Unknown database error"

	InfraErrCodeRedisConnect = "INFRA:REDIS:CONNECT"
	InfraErrMsgRedisConnect  = "Infra error: Failed to connect to Redis"

	InfraErrCodeRedisNotFound = "INFRA:REDIS:NOT_FOUND"
	InfraErrMsgRedisNotFound  = "Infra error: Not found resource in Redis"

	InfraErrCodeRedisGet = "INFRA:REDIS:GET"
	InfraErrMsgRedisGet  = "Infra error: Failed to get data from Redis"

	InfraErrCodeRedisSet = "INFRA:REDIS:SET"
	InfraErrMsgRedisSet  = "Infra error: Failed to write data to Redis"

	InfraErrCodeRedisUnknown = "INFRA:REDIS:UNKNOWN"
	InfraErrMsgRedisUnknown  = "Infra error: Unknown redis error"

	InfraErrCodeHTTPUnknown = "INFRA:HTTP:UNKNOWN"
	InfraErrMsgHTTPUnknown  = "Infra error: Unknown HTTP error"

	InfraErrCodeRPCUnknown = "INFRA:RPC:UNKNOWN"
	InfraErrMsgRPCUnknown  = "Infra error: Unknown RPC error"

	InfraErrCodeElsConnect = "INFRA:ELS:CONNECT"
	InfraErrMsgElsConnect  = "Infra error: Failed to connect to Els"

	InfraErrCodeElsUnknown = "INFRA:ELS:UNKNOWN"
	InfraErrMsgElsUnknown  = "Infra error: Unknown Els error"

	InfraErrCodeElsNotFound = "INFRA:ELS:NOT_FOUND"
	InfraErrMsgElsNotFound  = "Infra error: Not found resource in Els"
)
