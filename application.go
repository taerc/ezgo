package ezgo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
)

type Executor func(data interface{}) error

type Message struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message"`
	RequestId string      `json:"request_id"`
}

func JsonBind(ctx *gin.Context, data interface{}) error {
	if err := ctx.BindJSON(data); err != nil {
		return NewEError(CodeJsonFormatError, err)
	}
	return nil
}

func IndJsonResponse(ctx *gin.Context, er int, data interface{}) {
	ctx.IndentedJSON(Success, Message{
		Code:      er,
		Data:      data,
		Message:   GetMessageByCode(er),
		RequestId: GetRequestId(ctx),
	})
}

func OKResponse(ctx *gin.Context, data interface{}) {
	ctx.IndentedJSON(Success, Message{
		Code:      Success,
		Data:      data,
		Message:   GetMessageByCode(Success),
		RequestId: GetRequestId(ctx),
	})
}

func ErrorResponse(ctx *gin.Context, e error) {
	ctx.JSON(Success, Message{
		Code:      GetErrorCode(e),
		Data:      nil,
		Message:   e.Error(),
		RequestId: GetRequestId(ctx),
	})
}

type Application struct {
	Init   Executor
	Exec   Executor
	Done   Executor
	engine *gin.Engine
}

func (af *Application) APIGroup(ver, relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return af.engine.Group(path.Join("/api", ver, relativePath), handlers...)
}
func (af *Application) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return af.engine.Group(relativePath, handlers...)
}

func (af *Application) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return af.engine.Use(middleware...)
}

func (af *Application) Run(ipaddress ...string) error {
	return af.engine.Run(ipaddress...)
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

func (af *Application) Do(data interface{}) error {

	Info(nil, M, fmt.Sprintf("Version :%s", Version()))

	if n := af.Init(data); n != nil {
		return n
	}

	if n := af.Exec(data); n != nil {
		return n
	}

	if n := af.Done(data); n != nil {
		return n
	}

	return nil
}

var application *Application = nil

func init() {
	application = new(Application)
	application.engine = gin.Default()
	application.Use(PluginRequestId(), PluginCors(), PluginRequestSnapShot())
}

/// Application part

func DefaultApp() *Application {
	return application
}

func InitApplication(init, exec, done Executor) *Application {
	application.Init = init
	application.Exec = exec
	application.Done = done
	return application
}

func NewApplication(init, exec, done Executor) *Application {
	af := new(Application)
	af.Init = init
	af.Exec = exec
	af.Done = done
	af.engine = new(gin.Engine)
	return af
}

func Run(ipaddress ...string) error {
	return application.engine.Run(ipaddress...)
}

func Do(data interface{}) error {
	Info(nil, M, fmt.Sprintf("version: %s", Version()))
	return application.Do(data)
}
