package main

import (
	"flag"
	"fmt"

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
	ezgo.LoadModule(
		httpmod.WithModuleGitLab(),
	)
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

//go:generate swag init
func main() {

	httpmod.DefaultApp().Do(nil)
}
