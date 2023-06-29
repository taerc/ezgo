package ezgo

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/taerc/ezgo/conf"
	"os"
	"path"
	"sync"
	"time"
)

func init() {
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})
	log.SetOutput(os.Stdout)
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

func InitLogger(c *conf.Configure) {

	if _, err := os.Stat(c.LogDir); os.IsNotExist(err) {
		_ = os.MkdirAll(c.LogDir, os.ModePerm)
	}
	log_path := path.Join(c.LogDir, c.LogFileName)

	writer, err := rotatelogs.New(
		log_path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(log_path),                                          // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(24*c.LogMaxAge)*time.Hour),             // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(24*c.LogRotationTime)*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		Error(nil, "LOG", errors.WithStack(err))
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer},
		&nested.Formatter{
			HideKeys:    true,
			FieldsOrder: []string{"component", "category"},
			//CustomCallerFormatter: func(f *runtime.Frame) string {
			//	s := strings.Split(f.Function, ".")
			//	funcName := s[len(s)-1]
			//	return fmt.Sprintf(" %s:%d %s", path.Base(f.File), f.Line, funcName)
			//},
		},

	)
	log.AddHook(lfHook)
}

// WithComponentLogger @description: 注册日志信息
func WithComponentLogger(c *conf.Configure) Component {
	return func(wg *sync.WaitGroup) {
		wg.Done()
		InitLogger(c)
		Info(nil, M, "Finished Load Logger !")
	}

}
