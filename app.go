package ezgo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
)

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
		Message:   GetErrorMessage(er),
		RequestId: GetRequestId(ctx),
	})
}
func (gf *GinFlow) ResponseIndJson(ctx *gin.Context, er int, data interface{}) {
	ctx.IndentedJSON(Success, Response{
		Code:      er,
		Data:      data,
		Message:   GetErrorMessage(er),
		RequestId: GetRequestId(ctx),
	})
}

type AppFlow struct {
	Init Executor
	Exec Executor
	Done Executor
	Application
}

func (af *AppFlow) APIGroup(ver, relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return af.Engine.Group(path.Join("/api", ver, relativePath), handlers...)
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

func ProcPOST(route *gin.RouterGroup, relativePath string, processor Processor) {
	POST(route, relativePath, processor.PreProc, processor.Proc, processor.PreProc)
}

func ProcGET(route *gin.RouterGroup, relativePath string, processor Processor) {
	GET(route, relativePath, processor.PreProc, processor.Proc, processor.PostProc)
}
func ProcDELETE(route *gin.RouterGroup, relativePath string, processor Processor) {
	DELETE(route, relativePath, processor.PreProc, processor.Proc, processor.PostProc)
}
func ProcPATCH(route *gin.RouterGroup, relativePath string, processor Processor) {
	PATCH(route, relativePath, processor.PreProc, processor.Proc, processor.PostProc)
}
func ProcPUT(route *gin.RouterGroup, relativePath string, processor Processor) {
	PUT(route, relativePath, processor.PreProc, processor.Proc, processor.PostProc)
}
func ProcOPTIONS(route *gin.RouterGroup, relativePath string, processor Processor) {
	OPTIONS(route, relativePath, processor.PreProc, processor.Proc, processor.PostProc)
}
func ProcHEAD(route *gin.RouterGroup, relativePath string, processor Processor) {
	HEAD(route, relativePath, processor.PreProc, processor.Proc, processor.PostProc)
}

func POST(route *gin.RouterGroup, relativePath string, handlrs ...gin.HandlerFunc) {
	route.POST(relativePath, handlrs...)
}
func GET(route *gin.RouterGroup, relativePath string, handlrs ...gin.HandlerFunc) {
	route.GET(relativePath, handlrs...)
}
func DELETE(route *gin.RouterGroup, relativePath string, handlrs ...gin.HandlerFunc) {
	route.DELETE(relativePath, handlrs...)
}
func PATCH(route *gin.RouterGroup, relativePath string, handlrs ...gin.HandlerFunc) {
	route.PATCH(relativePath, handlrs...)
}
func PUT(route *gin.RouterGroup, relativePath string, handlrs ...gin.HandlerFunc) {
	route.PUT(relativePath, handlrs...)
}
func OPTIONS(route *gin.RouterGroup, relativePath string, handlrs ...gin.HandlerFunc) {
	route.OPTIONS(relativePath, handlrs...)
}
func HEAD(route *gin.RouterGroup, relativePath string, handlrs ...gin.HandlerFunc) {
	route.HEAD(relativePath, handlrs...)
}

func Version() string {
	return version
}

func (af *AppFlow) Do(data interface{}) int {

	Info(nil, M, fmt.Sprintf("Version :%s", Version()))

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
