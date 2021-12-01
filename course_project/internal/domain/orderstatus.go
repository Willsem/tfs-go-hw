package domain

type OrderStatus string

const (
	Empty                      OrderStatus = ""
	Placed                     OrderStatus = "placed"
	Cancelled                  OrderStatus = "cancelled"
	InvalidOrderType           OrderStatus = "invalidOrderType"
	InvalidSide                OrderStatus = "invalidSide"
	InvalidSize                OrderStatus = "invalidSize"
	InvalidPrice               OrderStatus = "invalidPrice"
	InsufficientAvailableFunds OrderStatus = "insufficientAvailableFunds"
	SelfFill                   OrderStatus = "selfFill"
	TooManySmallOrders         OrderStatus = "tooManySmallOrders"
	MaxPositionViolation       OrderStatus = "maxPositionViolation"
	MarketSuspended            OrderStatus = "marketSuspended"
	MarketInactive             OrderStatus = "marketInactive"
	ClientOrderIdAlreadyExist  OrderStatus = "clientOrderIdAlreadyExist"
	ClientOrderIdTooLong       OrderStatus = "clientOrderIdTooLong"
	OutsidePriceCollar         OrderStatus = "outsidePriceCollar"
	PostWouldExecute           OrderStatus = "postWouldExecute"
	IocWouldNotExecute         OrderStatus = "iocWouldNotExecute"
	WouldCauseLiquidation      OrderStatus = "wouldCauseLiquidation"
	WouldNotReducePosition     OrderStatus = "wouldNotReducePosition"
)
