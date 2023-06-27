package ezgo

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

var ConfigPath string
var Version string
var ShowVersion bool

type Executor func(data interface{}) int

type Application struct {
	Conf ConfParams
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
