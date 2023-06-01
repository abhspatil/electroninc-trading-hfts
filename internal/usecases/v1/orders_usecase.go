package v1

import (
	"container/heap"
	"fmt"
	"github.com/abhspatil/electronic-trading/constants"
	domainv1 "github.com/abhspatil/electronic-trading/internal/domain/v1"
	"github.com/abhspatil/electronic-trading/logger"
	dtov1 "github.com/abhspatil/electronic-trading/pkg/dto/v1"
	"github.com/abhspatil/electronic-trading/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var OrderInfo domainv1.OrderInfo
var orderBook *domainv1.OrderBook

var orderBookWithHeap *domainv1.OrderBookWithHeap

func MakeOffer(c *gin.Context, order dtov1.Order) (*domainv1.OrderInfo, error) {
	logctx := logger.Logger(c).WithFields(logrus.Fields{"method": "MakeOffer", "payload": order})
	logctx.Info("received order")

	if order.OrderType == constants.BuyOrder {
		OrderInfo.TotalBuyQty += order.Quantity
		OrderInfo.TotalBuyValue += (float64(order.Quantity) * order.Price)
	} else if order.OrderType == constants.SellOrder {
		OrderInfo.TotalSellQty += order.Quantity
		OrderInfo.TotalSellValue += (float64(order.Quantity) * order.Price)
	} else {
		logctx.Error("invalid order type")
		return nil, fmt.Errorf("invalid order type")
	}

	CalculateStockPrice(&OrderInfo)
	return &OrderInfo, nil
}

func MatchOffer(c *gin.Context, order dtov1.Order) (*domainv1.OrderBook, error) {
	logctx := logger.Logger(c).WithFields(logrus.Fields{"method": "MatchOffer", "payload": order})
	logctx.Info("received match order")

	if orderBook == nil {
		orderBook = &domainv1.OrderBook{BuyOrders: []*dtov1.Order{}, SellOrders: []*dtov1.Order{}, MatchedOrders: []*domainv1.MatchedOrders{}}
	}

	if order.OrderType == constants.BuyOrder {
		orderBook.BuyOrders = append(orderBook.BuyOrders, &order)
	} else if order.OrderType == constants.SellOrder {
		orderBook.SellOrders = append(orderBook.SellOrders, &order)
	} else {
		logctx.Error("invalid order type")
		return nil, fmt.Errorf("invalid order type")
	}

	MatchOrdersOrderN2(orderBook)

	return orderBook, nil
}

func MatchOfferWithHeap(c *gin.Context, order dtov1.Order) (*domainv1.OrderBook, error) {
	logctx := logger.Logger(c).WithFields(logrus.Fields{"method": "MatchOrdersWithHeap", "payload": order})
	logctx.Info("received match order")

	if orderBookWithHeap == nil {
		orderBookWithHeap = &domainv1.OrderBookWithHeap{BuyOrders: make(utils.OrderHeap, 0), SellOrders: make(utils.OrderHeap, 0),
			MatchedOrders: []*domainv1.MatchedOrders{}}
	}

	if order.OrderType == constants.BuyOrder {
		order.Price = -order.Price
		heap.Push(&orderBookWithHeap.BuyOrders, &order)
	} else if order.OrderType == constants.SellOrder {
		heap.Push(&orderBookWithHeap.SellOrders, &order)
	} else {
		logctx.Error("invalid order type")
		return nil, fmt.Errorf("invalid order type")
	}

	MatchOrdersWithHeapNLogN(orderBookWithHeap)

	return domainv1.OrderBookWithHeapToResp(*orderBookWithHeap), nil
}

func MatchOrdersWithHeapNLogN(orderBookWithHeap *domainv1.OrderBookWithHeap) {
	if len(orderBookWithHeap.BuyOrders) == 0 || len(orderBookWithHeap.SellOrders) == 0 {
		return
	}

	fmt.Println(orderBookWithHeap.BuyOrders)
	fmt.Println(orderBookWithHeap.SellOrders)

	for (len(orderBookWithHeap.BuyOrders) != 0 && len(orderBookWithHeap.SellOrders) != 0) && -(orderBookWithHeap.BuyOrders[0]).Price >= (*orderBookWithHeap.SellOrders[0]).Price {
		fmt.Println(orderBookWithHeap.BuyOrders[0], orderBookWithHeap.SellOrders[0])

		if orderBookWithHeap.BuyOrders[0].Quantity == orderBookWithHeap.SellOrders[0].Quantity {
			orderBookWithHeap.MatchedOrders = append(orderBookWithHeap.MatchedOrders,
				&domainv1.MatchedOrders{BuyOrder: orderBookWithHeap.BuyOrders[0], SellOrder: orderBookWithHeap.SellOrders[0]})
			heap.Pop(&orderBookWithHeap.BuyOrders)
			heap.Pop(&orderBookWithHeap.SellOrders)
		} else if orderBookWithHeap.BuyOrders[0].Quantity < orderBookWithHeap.SellOrders[0].Quantity {
			orderBookWithHeap.MatchedOrders = append(orderBookWithHeap.MatchedOrders,
				&domainv1.MatchedOrders{BuyOrder: orderBookWithHeap.BuyOrders[0], SellOrder: orderBookWithHeap.SellOrders[0]})

			sellOrder := orderBookWithHeap.SellOrders[0]
			sellOrder.Quantity -= orderBookWithHeap.BuyOrders[0].Quantity

			heap.Pop(&orderBookWithHeap.SellOrders)
			heap.Push(&orderBookWithHeap.SellOrders, sellOrder)

			heap.Pop(&orderBookWithHeap.BuyOrders)
		} else {
			orderBookWithHeap.MatchedOrders = append(orderBookWithHeap.MatchedOrders,
				&domainv1.MatchedOrders{BuyOrder: orderBookWithHeap.BuyOrders[0], SellOrder: orderBookWithHeap.SellOrders[0]})

			buyOrder := orderBookWithHeap.BuyOrders[0]
			buyOrder.Quantity -= orderBookWithHeap.SellOrders[0].Quantity

			heap.Pop(&orderBookWithHeap.BuyOrders)
			heap.Push(&orderBookWithHeap.BuyOrders, buyOrder)

			heap.Pop(&orderBookWithHeap.SellOrders)
		}
	}
}

func MatchOrdersOrderN2(orderBook *domainv1.OrderBook) {
	// buyOrders := []*dtov1.Order{}
	// sellOrders := []*dtov1.Order{}

	if len(orderBook.SellOrders) == 0 || len(orderBook.BuyOrders) == 0 {
		return
	}

	for i := 0; i < len(orderBook.BuyOrders); i++ {
		for j := 0; j < len(orderBook.SellOrders); j++ {

			if orderBook.BuyOrders[i].Price >= orderBook.SellOrders[j].Price && orderBook.SellOrders[j].Quantity != 0 {

				if orderBook.BuyOrders[i].Quantity == orderBook.SellOrders[j].Quantity {
					orderBook.MatchedOrders = append(orderBook.MatchedOrders, &domainv1.MatchedOrders{BuyOrder: orderBook.BuyOrders[i], SellOrder: orderBook.SellOrders[j]})
					orderBook.BuyOrders = append(orderBook.BuyOrders[:i], orderBook.BuyOrders[i+1:]...)
					orderBook.SellOrders = append(orderBook.SellOrders[:j], orderBook.SellOrders[j+1:]...)
					i--
					j--
				} else if orderBook.BuyOrders[i].Quantity > orderBook.SellOrders[j].Quantity {
					orderBook.BuyOrders[i].Quantity -= orderBook.SellOrders[j].Quantity
					orderBook.MatchedOrders = append(orderBook.MatchedOrders, &domainv1.MatchedOrders{BuyOrder: orderBook.BuyOrders[i], SellOrder: orderBook.SellOrders[j]})
					orderBook.SellOrders = append(orderBook.SellOrders[:j], orderBook.SellOrders[j+1:]...)
					j--
				} else {
					orderBook.SellOrders[j].Quantity -= orderBook.BuyOrders[j].Quantity
					orderBook.MatchedOrders = append(orderBook.MatchedOrders, &domainv1.MatchedOrders{BuyOrder: orderBook.BuyOrders[i], SellOrder: orderBook.SellOrders[j]})
					orderBook.BuyOrders = append(orderBook.BuyOrders[:i], orderBook.BuyOrders[i+1:]...)
					i--
				}
			}
		}
	}

	// orderBook.BuyOrders = buyOrders
	// orderBook.SellOrders = sellOrders
}

func CalculateStockPrice(orderInfo *domainv1.OrderInfo) {
	if orderInfo.TotalBuyQty == 0 || orderInfo.TotalSellQty == 0 {
		return
	}

	orderInfo.StokPrice = (orderInfo.TotalBuyValue + orderInfo.TotalSellValue) / (float64(orderInfo.TotalBuyQty) + float64(orderInfo.TotalSellQty))
}
