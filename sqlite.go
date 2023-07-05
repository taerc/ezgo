package ezgo

import (
	"github.com/taerc/ezgo/conf"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"sync"
)

var sqliteDb *gorm.DB = nil

func initSqlite(conf *conf.Configure) error {

	var e error = nil
	if sqliteDb, e = gorm.Open(sqlite.Open(conf.SQLitePath), &gorm.Config{}); e != nil {
		return e
	}
	return nil
}

func SQLITE() *gorm.DB {
	return sqliteDb
}

func WithCommponentSqlite(c *conf.Configure) Component {

	return func(wg *sync.WaitGroup) {
		wg.Done()
		initSqlite(c)
		Info(nil, M, "Finished Load SQLITE")
	}

}
