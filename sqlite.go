package ezgo

import (
	"fmt"
	"github.com/taerc/ezgo/conf"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
)

var defSqliteDb *gorm.DB = nil
var sqliteMap sync.Map

func initSqlite(name string, conf *conf.ConfSQLite) error {

	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //使用单数表名，启用该选项时，`User` 的表名应该是 `user`而不是users
		},
	}
	if db, e := gorm.Open(sqlite.Open(conf.SQLitePath), &gormConfig); e != nil {
		return e
	} else {
		sqliteMap.Swap(name, db)
	}
	return nil
}

func SQLITE(name ...string) *gorm.DB {

	if len(name) == 0 || name[0] == Default {
		if defSqliteDb == nil {
			if e := initSqlite(Default, &conf.Config.SQLite); e != nil {
				Error(nil, M, fmt.Sprintf("unknown default db.%s (forgotten configure?)", name[0]))
			}
			return defSqliteDb
		}
		return defSqliteDb
	}

	v, ok := sqliteMap.Load(name[0])

	if !ok {
		Error(nil, M, fmt.Sprintf("unknown db.%s (forgotten configure?)", name[0]))
	}

	return v.(*gorm.DB)
}

func SqliteExists(name string) bool {
	_, ok := sqliteMap.Load(name)
	return ok
}

func WithComponentSqlite(name string, c *conf.ConfSQLite) Component {

	return func(wg *sync.WaitGroup) {
		wg.Done()
		initSqlite(name, c)
		Info(nil, M, fmt.Sprintf("Finished Load [%s]-SQLITE", name))
	}

}
