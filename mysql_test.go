package ezgo

import (
	"fmt"
	"github.com/taerc/ezgo/conf"
	"testing"
)

func Test_initMySQL(t *testing.T) {

	c := &conf.MySQLConf{
		MySQLHostname:          "172.10.50.239",
		MySQLPort:              "3306",
		MySQLUserName:          "maicro",
		MySQLPass:              "Maicro@2314",
		MySQLDBName:            "cloudbusi",
		MySQLMaxOpenConnection: 256,
		MySQLMaxIdleConnection: 30,
		Charset:                "utf8mb4",
		Loc:                    "Local",
		ParseTime:              "true",
		MaxLifeTime:            "1h",
	}

	initMySQL("def", c)
}

func Test_initEntDb(t *testing.T) {
	c := conf.MySQLConf{
		MySQLHostname: "127.0.0.1",
		MySQLPort:     "3306",
		MySQLUserName: "wp",
		MySQLPass:     "wORd@2314",
		MySQLDBName:   "buckets",
		Charset:       "utf8mb4",
		ParseTime:     "true",
		Loc:           "Local",
	}

	name := "mysql"

	if err := initEntDb(name, name, &c); err != nil {
		fmt.Println(err)
	}
}
