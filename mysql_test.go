package ezgo

import (
	"github.com/taerc/ezgo/conf"
	"testing"
)

func Test_initMySQL(t *testing.T) {

	c := &conf.Configure{
		MySQLHostname:          "172.10.50.239",
		MySQLPort:              "3306",
		MySQLUserName:          "maicro",
		MySQLPass:              "Maicro@2314",
		MySQLDBName:            "cloudbusi",
		MySQLMaxOpenConnection: 256,
		MySQLMaxIdleConnection: 30,
		Charset:                "utf8mb4",
		Loc:                    "Loc",
		ParseTime:              "true",
	}

	initMySQL(c)
}
