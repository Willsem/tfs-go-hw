package domain

type OrderStatus string

const (
	Empty                      = ""
	Placed                     = "placed"
	Cancelled                  = "cancelled"
	InvalidOrderType           = "invalidOrderType"
	InvalidSide                = "invalidSide"
	InvalidSize                = "invalidSize"
	InvalidPrice               = "invalidPrice"
	InsufficientAvailableFunds = "insufficientAvailableFunds"
	SelfFill                   = "selfFill"
	TooManySmallOrders         = "tooManySmallOrders"
	MaxPositionViolation       = "maxPositionViolation"
	MarketSuspended            = "marketSuspended"
	MarketInactive             = "marketInactive"
	ClientOrderIdAlreadyExist  = "clientOrderIdAlreadyExist"
	ClientOrderIdTooLong       = "clientOrderIdTooLong"
	OutsidePriceCollar         = "outsidePriceCollar"
	PostWouldExecute           = "postWouldExecute"
	IocWouldNotExecute         = "iocWouldNotExecute"
	WouldCauseLiquidation      = "wouldCauseLiquidation"
	WouldNotReducePosition     = "wouldNotReducePosition"
)
