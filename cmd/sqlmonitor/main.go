package main

import (
	"context"
	"database/sql"
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

type tables struct {
	TableName  string       `sql:"TABLE_NAME"`
	TableType  string       `sql:"TABLE_TYPE"`
	CreateTime int64        `sql:"CREATE_TIME"`
	UpdateTime sql.NullTime `sql:"UPDATE_TIME"`
}

type columns struct {
	TableName     string         `sql:"TABLE_NAME"`
	ColumnName    string         `sql:"COLUMN_NAME"`
	ColumnDefault sql.NullString `sql:"COLUMN_DEFAULT"`
	ColumnComment sql.NullString `sql:"COLMN_COMMENT"`
}

func main() {

	conf.LoadConfigure(ConfigPath)
	ezgo.LoadComponent(
		ezgo.WithComponentMySQL(ezgo.Default, &conf.Config.SQL),
	)
	ent.InitDB()

	test()

}

func queryTables(ctx context.Context, schema string) ([]tables, error) {
	rows, e := ent.DB.QueryContext(ctx, "select TABLE_NAME, TABLE_TYPE, CREATE_TIME, UPDATE_TIME FROM TABLES WHERE TABLE_SCHEMA=? ", schema)

	totalTables := make([]tables, 0)

	if e != nil {
		return nil, e
	}
	defer rows.Close()

	for rows.Next() {

	}

	return totalTables, nil
}

func test() {
	ctx := context.Background()

	rows, e := ent.DB.QueryContext(ctx, "select TABLE_NAME, COLUMN_NAME, COLUMN_DEFAULT, COLUMN_COMMENT FROM COLUMNS WHERE TABLE_SCHEMA=? ", "pd_jibei")
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	defer rows.Close()

	cols := make([]columns, 0)

	for rows.Next() {

		var col columns
		if se := rows.Scan(&col.TableName, &col.ColumnName, &col.ColumnDefault, &col.ColumnComment); se != nil {
			fmt.Println(se.Error())
			continue
		} else {
			cols = append(cols, col)
		}
	}
	for _, v := range cols {
		fmt.Println(v.TableName, v.ColumnName, v.ColumnComment.String, v.ColumnDefault)
	}

}
