package ezgo

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
)

var ConfigPath string
var Version string
var ShowVersion bool

type Executor func(data interface{}) int

type Application struct {
	HTTP *GinContext
}

type AppFlow struct {
	Init Executor
	Exec Executor
	Done Executor
	Application
}

func (af *AppFlow) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return af.HTTP.Engine.Group(relativePath, handlers...)
}

func (af *AppFlow) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return af.HTTP.Engine.Use(middleware...)
}

func (af *AppFlow) Run(ipaddress ...string) error {
	return af.HTTP.Engine.Run(ipaddress...)
}
func (af *AppFlow) SetWhiteList(basePath, relativePath string) {
	af.HTTP.SetWhiteList(basePath, relativePath)
}

func (af *AppFlow) Do(data interface{}) int {

	// show version

	if ShowVersion {
		fmt.Println("version : ", Version)
		os.Exit(0)
	}

	if n := af.Init(data); n != Success {
		return n
	}

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
	appFlow.HTTP = NewGinContext()
	appFlow.Use(PluginRequestId(), PluginCors())
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
	appFlow.Init = init
	appFlow.Exec = exec
	appFlow.Done = done
	return appFlow
}

func Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return appFlow.HTTP.Engine.Group(relativePath, handlers...)
}

func Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return appFlow.HTTP.Engine.Use(middleware...)
}

func Run(ipaddress ...string) error {
	return appFlow.HTTP.Engine.Run(ipaddress...)
}

func Do(data interface{}) int {
	return appFlow.Do(data)
}

/// Default Plugin part
var headerXRequestID string = "X-Request-ID"

func PluginRequestId() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Get id from request
		rid := c.GetHeader(headerXRequestID)
		if rid == "" {
			rid = uuid.New().String()
			c.Request.Header.Add(headerXRequestID, rid)
		}
		// Set the id to ensure that the request-id is in the response
		c.Header(headerXRequestID, rid)
		c.Next()
	}
}

func GetRequestId(c *gin.Context) string {
	return c.Writer.Header().Get(headerXRequestID)
}

//

func PluginCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,Authorization-Token,AuthorizationToken")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			//c.Header("Access-Control-Max-Age", "172800")          // 缓存请求信息 单位为秒[48小时]
			//c.Header("Access-Control-Allow-Credentials", "false") //  跨域请求是否需要带cookie信息 默认设置为true
			//c.Set("content-type", "application/json") // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}
