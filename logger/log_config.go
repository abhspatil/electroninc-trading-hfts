package logger

import (
	"os"

	"github.com/abhspatil/electronic-trading/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func InitializeLogger() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.TraceLevel)
}

func Logger(c *gin.Context) *logrus.Entry {
	id, ok := c.Get(constants.CorrelationId)
	if !ok {
		return logrus.WithField(constants.CorrelationId, (uuid.New()).String())
	}
	return logrus.WithField(constants.CorrelationId, id)
}
