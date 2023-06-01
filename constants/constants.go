package constants

var CorrelationId = "correlation-id"

type OrderType string

const (
	SellOrder OrderType = "SELL"
	BuyOrder  OrderType = "BUY"
)
