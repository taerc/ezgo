package httpmod

import (
	"fmt"
	"path"

	"github.com/gin-gonic/gin"
	ezgo "github.com/taerc/ezgo/pkg"
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
		return ezgo.NewEError(ezgo.CodeJsonFormatError, err)
	}
	return nil
}

func IndJsonResponse(ctx *gin.Context, er int, data interface{}) {
	ctx.IndentedJSON(ezgo.Success, Message{
		Code:      er,
		Data:      data,
		Message:   ezgo.GetMessageByCode(er),
		RequestId: ezgo.GetRequestId(ctx),
	})
}

func OKResponse(ctx *gin.Context, data interface{}) {
	ctx.IndentedJSON(ezgo.Success, Message{
		Code:      ezgo.Success,
		Data:      data,
		Message:   ezgo.GetMessageByCode(ezgo.Success),
		RequestId: ezgo.GetRequestId(ctx),
	})
}

func ErrorResponse(ctx *gin.Context, e error) {
	ctx.JSON(ezgo.Success, Message{
		Code:      ezgo.GetErrorCode(e),
		Data:      nil,
		Message:   e.Error(),
		RequestId: ezgo.GetRequestId(ctx),
	})
}

type GinApplication struct {
	engine *gin.Engine
	Init   Executor
	Exec   Executor
	Done   Executor
}

func (af *GinApplication) APIGroup(ver, relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return af.engine.Group(path.Join("/api", ver, relativePath), handlers...)
}
func (af *GinApplication) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return af.engine.Group(relativePath, handlers...)
}

func (af *GinApplication) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return af.engine.Use(middleware...)
}

func (af *GinApplication) Run(ipaddress ...string) error {
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
	return ""
	// return ezgo.version
}

func (af *GinApplication) Do(data interface{}) error {
	// ezgo.Info(nil, M, fmt.Sprintf("Version :%s", Version()))

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

var application *GinApplication = nil

func init() {
	application = new(GinApplication)
	application.engine = gin.Default()
	application.Init = func(data interface{}) error {
		fmt.Println("weclome ginapplication")
		return nil
	}
	application.Exec = func(data interface{}) error {
		fmt.Println("exec ginapplication")
		return nil
	}
	application.Done = func(data interface{}) error {
		fmt.Println("done ginapplication")
		return nil
	}
	application.Use(PluginRequestId(), PluginCors())
}

/// GinApplication part

func DefaultApp() *GinApplication {
	return application
}

func InitGinApplication(init, exec, done Executor) *GinApplication {
	application.Init = init
	application.Exec = exec
	application.Done = done
	return application
}

// default GinApplication method
func Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return application.engine.Group(relativePath, handlers...)
}

func Run(ipaddress ...string) error {
	return application.engine.Run(ipaddress...)
}

func Do(data interface{}) error {
	// Info(nil, M, fmt.Sprintf("version: %s", Version()))
	return application.Do(data)
}

func NewGinApplication(init, exec, done Executor) *GinApplication {
	af := new(GinApplication)
	af.Init = init
	af.Exec = exec
	af.Done = done
	af.engine = new(gin.Engine)
	return af
}
