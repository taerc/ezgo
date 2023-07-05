package ezgo

import (
	"github.com/taerc/ezgo/conf"
	"sync"
)

func initSqlite(conf *conf.Configure) error {

	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	//if err != nil {
	//	panic("failed to connect database")
	//}
	// Migrate the schema
	//db.AutoMigrate(&Product{})

	return nil
}

func WithCommponentSqlite(c *conf.Configure) Component {

	return func(wg *sync.WaitGroup) {
		wg.Done()
		initSqlite(c)
		Info(nil, M, "Finished Load SQLITE")
	}

}
