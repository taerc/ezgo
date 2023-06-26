package ezgo

import "gopkg.in/ini.v1"

func LoadConfigure(source interface{}, others ...interface{}) (*ini.File, error) {
	return ini.Load(source, others...)
}
