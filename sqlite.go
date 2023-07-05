package ezgo

import (
	"github.com/taerc/ezgo/conf"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
)

var sqliteDb *gorm.DB = nil

func initSqlite(conf *conf.Configure) error {

	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //使用单数表名，启用该选项时，`User` 的表名应该是 `user`而不是users
		},
	}

	var e error = nil
	if sqliteDb, e = gorm.Open(sqlite.Open(conf.SQLitePath), &gormConfig); e != nil {
		return e
	}
	return nil
}

func SQLITE() *gorm.DB {
	return sqliteDb
}

func WithComponentSqlite(c *conf.Configure) Component {

	return func(wg *sync.WaitGroup) {
		wg.Done()
		initSqlite(c)
		Info(nil, M, "Finished Load SQLITE")
	}

}
