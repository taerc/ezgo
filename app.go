package ezgo

import (
	"flag"
	"fmt"
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

func (af *AppFlow) Run(data interface{}) int {

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

func init() {
	flag.BoolVar(&ShowVersion, "version", false, "print program build version")
	flag.StringVar(&ConfigPath, "c", "conf/config.toml", "path of configure file.")
	flag.Parse()
}

func NewAppFlow(init, exec, done Executor) *AppFlow {

	af := new(AppFlow)
	af.Init = init
	af.Exec = exec
	af.Done = done
	af.HTTP = NewGinContext()

	return af
}
