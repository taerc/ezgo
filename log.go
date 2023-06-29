package ezgo

import (
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

var zapConf = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
	Development:      true,
	Encoding:         "console",
	OutputPaths:      []string{"stdout", "./zap.log"},
	ErrorOutputPaths: []string{"stderr"},
}

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
	c.OutputPaths = []string{"stdout", "zap.log"}
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
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.NewMultiWriteSyncer(ws...), zap.DebugLevel)
	return zap.New(core)
}

func initLogger(cfg *conf.Configure) {
	logger = newLogger(cfg)
}

// MyTimeEncoder 自定义时间格式化
func MyTimeEncoder(t time.Time, e zapcore.PrimitiveArrayEncoder) {
	e.AppendString(t.Format("2006-01-02 01:01:01"))
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

func Info(c *gin.Context, mod string, msg string) {
	reqId := zap.String("req", getReqId(c))
	m := zap.String("mod", mod)
	logger.Info(msg, reqId, m)
}

func Error(c *gin.Context, mod string, msg string) {
	reqId := zap.String("req", getReqId(c))
	m := zap.String("mod", mod)
	logger.Info(msg, reqId, m)
}

func Warn(c *gin.Context, mod string, msg string) {
	reqId := zap.String("req", getReqId(c))
	m := zap.String("mod", mod)
	logger.Info(msg, reqId, m)
}

func Panic(c *gin.Context, mod string, msg string) {
	reqId := zap.String("req", getReqId(c))
	m := zap.String("mod", mod)
	logger.Info(msg, reqId, m)
}

// WithComponentLogger @description: 注册日志信息
func WithComponentLogger(c *conf.Configure) Component {
	return func(wg *sync.WaitGroup) {
		wg.Done()
		initLogger(c)
		Info(nil, M, "Finished Load Logger !")
	}

}
