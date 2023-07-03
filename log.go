package ezgo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/taerc/ezgo/conf"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger = nil

func init() {
	c := zap.NewDevelopmentConfig()

	c.DisableCaller = true
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	c.EncoderConfig.EncodeTime = MyTimeEncoder
	c.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	c.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	c.Development = true
	c.Encoding = "console"
	c.ErrorOutputPaths = []string{"stderr"}
	c.OutputPaths = []string{"stdout"}
	logger, _ = c.Build()

}

func newLogger(c *conf.Configure) *zap.Logger {
	if len(c.LogFileName) == 0 {
		return logger
	}

	cfg := zap.NewProductionEncoderConfig()

	cfg.TimeKey = "time"
	cfg.EncodeTime = MyTimeEncoder
	cfg.EncodeCaller = zapcore.FullCallerEncoder

	ws := make([]zapcore.WriteSyncer, 0, 2)

	ws = append(ws, zapcore.AddSync(&lumberjack.Logger{
		Filename:   c.LogFileName,
		MaxSize:    c.LogMaxSize,
		MaxAge:     c.LogMaxAge,
		MaxBackups: c.LogMaxBackups,
		Compress:   c.LogCompress,
		LocalTime:  true,
	}))

	if c.LogStderr {
		ws = append(ws, zapcore.Lock(os.Stderr))
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.NewMultiWriteSyncer(ws...), zap.DebugLevel)
	return zap.New(core)
}

func initLogger(cfg *conf.Configure) {
	logger = newLogger(cfg)
}

// MyTimeEncoder 自定义时间格式化
func MyTimeEncoder(t time.Time, e zapcore.PrimitiveArrayEncoder) {
	e.AppendString(t.Format("2006-01-02 15:04:05"))
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

func InfoReq(reqId, mod, msg string, fields ...zap.Field) {
	logger.Info(fmt.Sprintf("%s %s %s", reqId, mod, msg), fields...)
}

func Info(c *gin.Context, mod string, msg string, fields ...zap.Field) {
	InfoReq(getReqId(c), mod, msg, fields...)
}

func ErrorReq(reqId, mod, msg string, fields ...zap.Field) {
	logger.Error(fmt.Sprintf("%s %s %s", reqId, mod, msg), fields...)
}

func Error(c *gin.Context, mod string, msg string, fields ...zap.Field) {
	ErrorReq(getReqId(c), mod, msg, fields...)
}
func WarnReq(reqId, mod, msg string, fields ...zap.Field) {
	logger.Warn(fmt.Sprintf("%s %s %s", reqId, mod, msg), fields...)
}
func Warn(c *gin.Context, mod string, msg string, fields ...zap.Field) {
	WarnReq(getReqId(c), mod, msg, fields...)
}

func PanicReq(reqId, mod, msg string, fields ...zap.Field) {
	logger.Panic(fmt.Sprintf("%s %s %s", reqId, mod, msg), fields...)
}

func Panic(c *gin.Context, mod string, msg string, fields ...zap.Field) {
	PanicReq(getReqId(c), mod, msg, fields...)
}

// WithComponentLogger @description: 注册日志信息
func WithComponentLogger(c *conf.Configure) Component {
	return func(wg *sync.WaitGroup) {
		wg.Done()
		initLogger(c)
		Info(nil, M, "Finished Load Logger !")
	}

}
