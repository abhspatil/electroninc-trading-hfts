package controller

import (
	usecasesv1 "github.com/abhspatil/electronic-trading/internal/usecases/v1"
	"github.com/abhspatil/electronic-trading/logger"
	dtov1 "github.com/abhspatil/electronic-trading/pkg/dto/v1"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func MakeOfferV1(c *gin.Context) {
	ctxlog := logger.Logger(c).WithFields(logrus.Fields{"method": "MakeOfferV1"})

	ctxlog.Info("got request for making offer")

	body := dtov1.Order{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid payload, Please try again, Err:" + err.Error(),
			"data":    nil,
		})
		return
	}

	resp, err := usecasesv1.MakeOffer(c, body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(200, resp)
	return
}

func MatchOrdersV1(c *gin.Context) {
	ctxlog := logger.Logger(c).WithFields(logrus.Fields{"method": "MatchOrdersV1"})

	ctxlog.Info("got request for matching orders")

	body := dtov1.Order{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid payload, Please try again, Err:" + err.Error(),
			"data":    nil,
		})
		return
	}

	resp, err := usecasesv1.MatchOffer(c, body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(200, resp)
	return
}

func MatchOrdersV2(c *gin.Context) {
	ctxlog := logger.Logger(c).WithFields(logrus.Fields{"method": "MatchOrdersV2"})

	ctxlog.Info("got request for matching orders")

	body := dtov1.Order{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid payload, Please try again, Err:" + err.Error(),
			"data":    nil,
		})
		return
	}

	resp, err := usecasesv1.MatchOfferWithHeap(c, body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(200, resp)
	return
}
