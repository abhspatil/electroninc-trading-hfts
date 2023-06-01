package dto

import "github.com/abhspatil/electronic-trading/constants"

type Order struct {
	Id        int64               `json:"id"`
	OrderType constants.OrderType `json:"order_type"`
	Quantity  int64               `json:"quantity"`
	Price     float64             `json:"price"`
}
