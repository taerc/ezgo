package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/taerc/ezgo/conf"
	"github.com/taerc/ezgo/database/sqlmonitor/ent"
	ezgo "github.com/taerc/ezgo/pkg"
)

var ConfigPath string
var ShowVersion bool

func init() {
	flag.BoolVar(&ShowVersion, "version", false, "print program build version")
	flag.StringVar(&ConfigPath, "c", "conf/config.toml", "path of configure file.")
	flag.Parse()
}

func main() {

	conf.LoadConfigure(ConfigPath)
	ezgo.LoadComponent(
		ezgo.WithComponentMySQL(ezgo.Default, &conf.Config.SQL),
	)
	ent.InitDB()

	test()

}

func test() {
	ctx := context.Background()
	rows, e := ent.DB.QueryContext(ctx, "select TABLE_NAME, COLUMN_NAME FROM COLUMNS WHERE TABLE_SCHEMA=? ", "pd_jibei")

	if e != nil {
		fmt.Println(e.Error())
		return
	}
	defer rows.Close()

	type Columns struct {
		TableName     string `sql:"TABLE_NAME"`
		ColumnName    string `sql:"COLUMN_NAME"`
		ColumnDefault string `sql:"COLUMN_DEFAULT"`
		ColumnComment string `sql:"COLMN_COMMENT"`
	}

	provinces := make([]Columns, 0)

	for rows.Next() {

		var p Columns
		if se := rows.Scan(&p.TableName, &p.ColumnName); se != nil {
			fmt.Println(se.Error())
			continue
		} else {
			provinces = append(provinces, p)
		}
	}
	for _, v := range provinces {
		fmt.Println(v.TableName, v.ColumnName)
	}

}
