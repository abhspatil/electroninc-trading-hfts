package main

import (
	"github.com/abhspatil/electronic-trading/logger"
	"github.com/abhspatil/electronic-trading/pkg/http"
	"github.com/gin-gonic/gin"
)

func init() {
	logger.InitializeLogger()
	http.InitEndpoints()
}

func main() {
	logger.Logger(&gin.Context{}).Info("Welcome to Patils World :) ")
}
