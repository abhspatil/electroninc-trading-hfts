package http

import (
	"github.com/abhspatil/electronic-trading/pkg/controller"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitEndpoints() {
	logrus.Info("initializing endpoints")
	r := gin.Default()
	r.POST("/v1/orders", controller.MakeOfferV1)
	r.POST("/v1/match-orders", controller.MatchOrdersV1)
	r.POST("/v2/match-orders", controller.MatchOrdersV2)

	r.Run(":8080")
	r.Run()
}
