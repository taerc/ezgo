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
	TableName  string       `sql:"TABLE_NAME" json:"table_name"`
	CreateTime int64        `sql:"CREATE_TIME" json:"create_time"`
	UpdateTime sql.NullTime `sql:"UPDATE_TIME" json:"update_time"`
}

type columns struct {
	TableName     string         `sql:"TABLE_NAME" json:"table_name"`
	ColumnName    string         `sql:"COLUMN_NAME" json:"column_name"`
	ColumnDefault sql.NullString `sql:"COLUMN_DEFAULT" json:"column_default"`
	ColumnComment sql.NullString `sql:"COLMN_COMMENT" json:"column_comment"`
}

type columsMonitor struct {
	newTable   []tables
	delTable   []tables
	newColumns []columns
	delColumns []columns
}

func (c *columsMonitor) AddNewTable(table string) {
	v := tables{
		TableName: table,
	}
	c.newTable = append(c.newTable, v)
}

func (c *columsMonitor) AddDelTable(table string) {
	v := tables{
		TableName: table,
	}
	c.delTable = append(c.delTable, v)
}
func (c *columsMonitor) AddNewColumn(table, column string) {
	v := columns{
		TableName:  table,
		ColumnName: column,
	}
	c.newColumns = append(c.newColumns, v)
}
func (c *columsMonitor) AddDelColumn(table, column string) {
	v := columns{
		TableName:  table,
		ColumnName: column,
	}
	c.delColumns = append(c.delColumns, v)
}

func (c *columsMonitor) Report() {
	fmt.Println("新增表:")
	for _, t := range c.newTable {
		fmt.Println(t.TableName)
	}
	fmt.Println("删除表:")
	for _, t := range c.delTable {
		fmt.Println(t.TableName)
	}

	fmt.Println("新增列:")
	for _, t := range c.newColumns {
		fmt.Println(t.TableName, t.ColumnName)
	}
	fmt.Println("删除列:")
	for _, t := range c.delColumns {
		fmt.Println(t.TableName, t.ColumnName)
	}

}

func newColumnMonitor() *columsMonitor {

	return &columsMonitor{
		newTable:   make([]tables, 0),
		delTable:   make([]tables, 0),
		newColumns: make([]columns, 0),
		delColumns: make([]columns, 0),
	}

}

func main() {

	conf.LoadConfigure(ConfigPath)
	ezgo.LoadComponent(
		ezgo.WithComponentMySQL(ezgo.Default, &conf.Config.SQL),
	)
	ent.InitDB()

	tablePath := "tables.json"
	columnPath := "columns.json"
	tableSchema := "pd_jibei"

	ctx := context.Background()
	tbls, e := queryTables(ctx, tableSchema)

	if e != nil {
		fmt.Println(e)
		return
	}

	hisTables := make([]tables, 0)
	if ezgo.PathExists(tablePath) {
		ezgo.LoadJson(tablePath, &hisTables)
	}
	ezgo.SaveJson(tablePath, tbls)

	cm := newColumnMonitor()

	// 新增表
	for _, nt := range tbls {
		if !inTableList(nt.TableName, hisTables) {
			cm.AddNewTable(nt.TableName)
		}
	}

	// 删除表
	for _, nt := range hisTables {
		if !inTableList(nt.TableName, tbls) {
			cm.AddDelTable(nt.TableName)
		}
	}

	currentColumns, e := queryColumns(ctx, tableSchema)

	if e != nil {
		fmt.Println(e)
		return
	}
	hisColumns := make([]columns, 0)
	if ezgo.PathExists(columnPath) {
		ezgo.LoadJson(columnPath, &hisColumns)
	}
	ezgo.SaveJson(columnPath, currentColumns)
	// 新增字段
	for _, nc := range currentColumns {
		if !inColumnList(nc.TableName, nc.ColumnName, hisColumns) {
			cm.AddNewColumn(nc.TableName, nc.ColumnName)
		}
	}

	// 删除字段
	for _, nc := range hisColumns {
		if !inColumnList(nc.TableName, nc.ColumnName, currentColumns) {
			cm.AddDelColumn(nc.TableName, nc.ColumnName)
		}
	}

	cm.Report()

}

func inTableList(target string, slist []tables) bool {

	for _, s := range slist {
		if s.TableName == target {
			return true
		}
	}
	return false
}

func inColumnList(table_name, column string, slist []columns) bool {

	for _, s := range slist {
		if s.TableName == table_name && s.ColumnName == column {
			return true
		}
	}
	return false
}

func queryTables(ctx context.Context, schema string) ([]tables, error) {
	rows, e := ent.DB.QueryContext(ctx, "select TABLE_NAME, CREATE_TIME, UPDATE_TIME FROM TABLES WHERE TABLE_SCHEMA=? ", schema)

	totalTables := make([]tables, 0)

	if e != nil {
		return nil, e
	}
	defer rows.Close()

	for rows.Next() {
		var t tables

		if e = rows.Scan(&t.TableName, &t.CreateTime, &t.UpdateTime); e != nil {
			totalTables = append(totalTables, t)

		} else {
			fmt.Println(e)
		}
	}

	return totalTables, nil
}

func queryColumns(ctx context.Context, schema string) ([]columns, error) {
	rows, e := ent.DB.QueryContext(ctx, "select TABLE_NAME, COLUMN_NAME, COLUMN_DEFAULT, COLUMN_COMMENT FROM COLUMNS WHERE TABLE_SCHEMA=? ", schema)
	if e != nil {
		fmt.Println(e.Error())
		return nil, e
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
	return cols, nil
}
