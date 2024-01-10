package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/taerc/ezgo/conf"
	"github.com/taerc/ezgo/httpmod"
	ezgo "github.com/taerc/ezgo/pkg"
	_ "gorm.io/driver/mysql"
)

var M string = "MAIN"

var ConfigPath string
var ShowVersion bool

func init() {
	flag.BoolVar(&ShowVersion, "version", false, "print program build version")
	flag.StringVar(&ConfigPath, "c", "conf/config.toml", "path of configure file.")
	flag.Parse()
}

func Init(data interface{}) error {

	conf.LoadConfigure(ConfigPath)
	ezgo.LoadComponent(
		// ezgo.WithComponentResource(conf.Config),
		//ezgo.WithComponentLogger(conf.Config),
		ezgo.WithComponentMySQL(ezgo.Default, &conf.Config.SQL),
	)

	ezgo.LoadModule()
	return nil
}

func Exec(data interface{}) error {

	httpmod.Run(fmt.Sprintf("%s:%s", conf.Config.Host, conf.Config.Port))
	return nil

}

func Done(data interface{}) error {

	return nil
}

func init() {
	httpmod.InitGinApplication(Init, Exec, Done)
}

type nginxMirror struct {
}

func (nm *nginxMirror) AddLog(ctx *gin.Context) {

	fmt.Println("Log ing ...  data data ...")

}

func WithModuleNginxMirror() func(wg *sync.WaitGroup) {
	return func(wg *sync.WaitGroup) {
		defer wg.Done()
		s := new(nginxMirror)
		route := httpmod.Group("/nginx/mirror/")
		httpmod.POST(route, "/add", s.AddLog)
	}
}

func main() {

	httpmod.Do(nil)
}
