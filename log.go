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

func init() {
	cfg := zap.NewDevelopmentConfig()

	cfg.DisableCaller = true
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = MyTimeEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	logger, _ = cfg.Build()

}

func newLogger(cfg *conf.Configure) *zap.Logger {
	if len(cfg.LogFileName) == 0 {
		return logger
	}

	c := zap.NewProductionEncoderConfig()

	c.TimeKey = "time"
	c.EncodeTime = MyTimeEncoder
	c.EncodeCaller = zapcore.FullCallerEncoder

	ws := make([]zapcore.WriteSyncer, 0, 2)

	ws = append(ws, zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.LogFileName,
		MaxSize:    cfg.LogMaxSize,
		MaxAge:     cfg.LogMaxAge,
		MaxBackups: cfg.LogMaxBackups,
		Compress:   cfg.LogCompress,
		LocalTime:  true,
	}))

	if cfg.LogStderr {
		ws = append(ws, zapcore.Lock(os.Stderr))
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(c), zapcore.NewMultiWriteSyncer(ws...), zap.DebugLevel)

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
