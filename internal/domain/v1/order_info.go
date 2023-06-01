package v1

import (
	"github.com/abhspatil/electronic-trading/pkg/dto/v1"
	"github.com/abhspatil/electronic-trading/pkg/utils"
)

type OrderInfo struct {
	TotalBuyQty    int64   `json:"total_buy_qty"`
	TotalBuyValue  float64 `json:"total_buy_value"`
	TotalSellValue float64 `json:"total_sell_value"`
	TotalSellQty   int64   `json:"total_sell_qty"`
	StokPrice      float64 `json:"stock_price"`
}

type OrderBook struct {
	BuyOrders     []*dto.Order
	SellOrders    []*dto.Order
	MatchedOrders []*MatchedOrders
}

type OrderBookWithHeap struct {
	BuyOrders     utils.OrderHeap
	SellOrders    utils.OrderHeap
	MatchedOrders []*MatchedOrders
}

type MatchedOrders struct {
	BuyOrder  *dto.Order
	SellOrder *dto.Order
}

func OrderBookWithHeapToResp(o OrderBookWithHeap) *OrderBook {
	var buyOrders []*dto.Order
	var sellOrders []*dto.Order

	if len(o.SellOrders) != 0 {
		for _, d := range o.SellOrders {
			sellOrders = append(sellOrders, d)
		}
	}

	if len(o.BuyOrders) != 0 {
		for _, d := range o.BuyOrders {
			d.Price = -d.Price
			buyOrders = append(buyOrders, d)
		}
	}

	return &OrderBook{BuyOrders: buyOrders, SellOrders: sellOrders, MatchedOrders: o.MatchedOrders}
}
