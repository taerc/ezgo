package ezgo

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"strings"
)

var ConfigPath string
var ShowVersion bool

type Executor func(data interface{}) int

type Application struct {
	Engine    *gin.Engine
	whiteList map[string]bool
	urlRole   map[string]int
}

type Response struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message"`
	RequestId string      `json:"request_id"`
}

// gin.context typedef 处理一下，后面可以自由替换
// 控制器的处理
type GinFlow struct {
}

type Processor interface {
	PreProc(ctx *gin.Context)
	Proc(ctx *gin.Context)
	PostProc(ctx *gin.Context)
}

// 路由分组的管理
// 模块化路由注册的管理
func (gf *GinFlow) PreProc(ctx *gin.Context) {

}
func (gf *GinFlow) PostProc(ctx *gin.Context) {

}

func (gf *GinFlow) BindJson(ctx *gin.Context, data interface{}) int {
	if err := ctx.BindJSON(data); err != nil {

		if !strings.HasPrefix(err.Error(), "json: invalid use of ,string struct tag,") {
			gf.ResponseJson(ctx, ErrInvalidJsonFormat, nil)
			return ErrInvalidJsonFormat
		} else {
			return Success
		}
	}
	return Success
}
func (gf *GinFlow) ResponseJson(ctx *gin.Context, er int, data interface{}) {
	ctx.JSON(Success, Response{
		Code:      er,
		Data:      data,
		Message:   "",
		RequestId: GetRequestId(ctx),
	})
}

type AppFlow struct {
	Init Executor
	Exec Executor
	Done Executor
	Application
}

func (af *AppFlow) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return af.Engine.Group(relativePath, handlers...)
}

func (af *AppFlow) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return af.Engine.Use(middleware...)
}

func (af *AppFlow) Run(ipaddress ...string) error {
	return af.Engine.Run(ipaddress...)
}
func (af *AppFlow) SetWhiteList(basePath, relativePath string) {
	af.whiteList[path.Join(basePath, relativePath)] = true
}

func SetPostProc(route *gin.RouterGroup, relativePath string, processor Processor) {
	route.POST(relativePath, processor.PreProc, processor.Proc, processor.PreProc)
}

func SetGetProc(route *gin.RouterGroup, relativePath string, processor Processor) {
	route.GET(relativePath, processor.PreProc, processor.Proc, processor.PostProc)
}

func POST(route *gin.RouterGroup, relativePath string, handlrs ...gin.HandlerFunc) {
	route.POST(relativePath, handlrs...)
}
func GET(route *gin.RouterGroup, relativePath string, handlrs ...gin.HandlerFunc) {
	route.GET(relativePath, handlrs...)
}

func Version() string {
	return version
}

func (af *AppFlow) Do(data interface{}) int {

	// show version

	if ShowVersion {
		fmt.Println("version : ", Version())
		os.Exit(0)
	}

	if n := af.Init(data); n != Success {
		return n
	}

	// fixed

	if n := af.Exec(data); n != Success {
		return n
	}

	if n := af.Done(data); n != Success {
		return n
	}

	return Success
}

var appFlow *AppFlow = nil

func init() {
	flag.BoolVar(&ShowVersion, "version", false, "print program build version")
	flag.StringVar(&ConfigPath, "c", "conf/config.toml", "path of configure file.")
	flag.Parse()
	appFlow = new(AppFlow)
	appFlow.Engine = gin.Default()
	appFlow.Use(PluginRequestId(), PluginCors(), PluginRequestSnapShot())
}

/// Application part

func DefaultApp() *AppFlow {
	return appFlow
}

func InitAppFlow(init, exec, done Executor) *AppFlow {
	appFlow.Init = init
	appFlow.Exec = exec
	appFlow.Done = done
	return appFlow
}

func NewAppFlow(init, exec, done Executor) *AppFlow {
	af := new(AppFlow)
	af.Init = init
	af.Exec = exec
	af.Done = done
	af.Engine = new(gin.Engine)
	return af
}

func Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return appFlow.Engine.Group(relativePath, handlers...)
}

func Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return appFlow.Engine.Use(middleware...)
}

func Run(ipaddress ...string) error {
	return appFlow.Engine.Run(ipaddress...)
}

func Do(data interface{}) int {
	Info(nil, M, fmt.Sprintf("version: %s", Version()))
	return appFlow.Do(data)
}
