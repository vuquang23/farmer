package errors

const (
	DomainErrCodeInitExchangeInfo = "DOMAIN:INIT_EXCHANGE_INFO"
	DomainErrMsgInitExchangeInfo  = "Domain error: Manager can not init exchange info"

	DomainErrCodeInitWavetrendData = "DOMAIN:INIT_WAVETREND_DATA"
	DomainErrMsgInitWavetrendData  = "Domain error: Worker can not init wavetrend data"

	DomainErrCodeWavetrendServiceNameExisted = "DOMAIN:WAVETREND_SERVICE_NAME_EXISTED"
	DomainErrMsgWavetrendServiceNameExisted  = "Domain error: Wavetrend service name is already registered with provider"

	DomainErrCodeCreateBuyOrderFailed = "DOMAIN:CREATE_BUY_ORDER_FAILED"
	DomainErrMsgCreateBuyOrderFailed  = "Domain error: Create buy order failed"
)
