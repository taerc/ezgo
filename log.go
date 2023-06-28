package ezgo

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})
}

func getReqId(c *gin.Context) string {

	if c == nil {
		return "-"
	}

	if s := GetRequestId(c); s == "" {
		return "-"
	} else {
		return s
	}
}

func Info(c *gin.Context, mod string, args ...interface{}) {
	log.Info(fmt.Sprintf("[%s] %s ", getReqId(c), mod), args)
}

func Error(c *gin.Context, mod string, args ...interface{}) {
	log.Error(fmt.Sprintf("[%s] %s ", getReqId(c), mod), args)
}

func Warn(c *gin.Context, mod string, args ...interface{}) {
	log.Warn(fmt.Sprintf("[%s] %s ", getReqId(c), mod), args)
}

func Panic(c *gin.Context, mod string, args ...interface{}) {
	log.Error(fmt.Sprintf("[%s] %s ", getReqId(c), mod), args)
}
