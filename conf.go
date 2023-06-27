package ezgo

import "gopkg.in/ini.v1"


type ConfParams struct {
	Host        string
	Port        string
}

func LoadConfigure(source interface{}, others ...interface{}) (*ini.File, error) {
	return ini.Load(source, others...)
}
