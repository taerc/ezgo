package ezgo

import (
	"flag"
	"fmt"
	"os"
)

type Executor func(conf *Application, data interface{}) int

type AppFlow struct {
	Init Executor
	Exec Executor
	Done Executor
}

type Application struct {
	ConfigPath  string
	Version     string
	ShowVersion bool
	Conf        ConfParams
	HTTP        *GinContext
}

func (af *AppFlow) Run(conf *Application, data interface{}) int {

	// show version

	if conf.ShowVersion {
		fmt.Println("version : ", conf.Version)
		os.Exit(0)
	}

	if n := af.Init(conf, data); n != Success {
		return n
	}

	if n := af.Exec(conf, data); n != Success {
		return n
	}

	if n := af.Done(conf, data); n != Success {
		return n
	}

	return Success
}

var APP *Application = nil

func init() {
	APP = new(Application)
	APP.HTTP = NewGinContext()

	flag.BoolVar(&APP.ShowVersion, "version", false, "print program build version")
	flag.StringVar(&APP.ConfigPath, "c", "conf/config.toml", "path of configure file.")
	flag.Parse()
}
