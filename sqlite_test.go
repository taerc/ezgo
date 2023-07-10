package ezgo

import (
	"fmt"
	"github.com/taerc/ezgo/conf"
	"testing"
)

func Test_initSqlite(t *testing.T) {
	c := &conf.ConfSQLite{SQLitePath: "test.db"}
	initSqlite("def", c)
	// create databases
	db := SQLITE()
	db.AutoMigrate(requestRecord{})

	db.Create(&requestRecord{RequestId: "1dd2345677",
		Method: "meth",
		Url:    "......",
		Body:   []byte{'A', 'B', 'C'}})

	var rr requestRecord
	db.Find(&rr)
	fmt.Println(rr.RequestId)
	fmt.Println(rr.Method)
	fmt.Println(string(rr.Body))

}
