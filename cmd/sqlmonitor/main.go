package main

import (
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"reflect"
	"text/template"

	"github.com/taerc/ezgo/conf"
	"github.com/taerc/ezgo/database/sqlmonitor/ent"
	"github.com/taerc/ezgo/dd"
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
	ColumnType    sql.NullString `sql:"COLMN_TYPE" json:"column_type"`
	ColumnDefault sql.NullString `sql:"COLUMN_DEFAULT" json:"column_default"`
	ColumnComment sql.NullString `sql:"COLMN_COMMENT" json:"column_comment"`
}

type columsMonitor struct {
	TableSchema   string
	NewTable      []tables
	DelTable      []tables
	NewColumns    []columns
	DelColumns    []columns
	UpdateColumns []columns
}

func (c *columsMonitor) AddNewTable(table string) {
	v := tables{
		TableName: table,
	}
	c.NewTable = append(c.NewTable, v)
}

func (c *columsMonitor) AddDelTable(table string) {
	v := tables{
		TableName: table,
	}
	c.DelTable = append(c.DelTable, v)
}
func (c *columsMonitor) AddNewColumn(table, column string) {
	v := columns{
		TableName:  table,
		ColumnName: column,
	}
	c.NewColumns = append(c.NewColumns, v)
}
func (c *columsMonitor) AddDelColumn(table, column string) {
	v := columns{
		TableName:  table,
		ColumnName: column,
	}
	c.DelColumns = append(c.DelColumns, v)
}
func (c *columsMonitor) AddUpdateColumn(table, column string) {
	v := columns{
		TableName:  table,
		ColumnName: column,
	}
	c.UpdateColumns = append(c.UpdateColumns, v)
}

func (c *columsMonitor) Report() {
	fmt.Println("新增表:")
	for _, t := range c.NewTable {
		fmt.Println(t.TableName)
	}
	fmt.Println("删除表:")
	for _, t := range c.DelTable {
		fmt.Println(t.TableName)
	}

	fmt.Println("新增列:")
	for _, t := range c.NewColumns {
		fmt.Println(t.TableName, t.ColumnName)
	}
	fmt.Println("删除列:")
	for _, t := range c.DelColumns {
		fmt.Println(t.TableName, t.ColumnName)
	}
}

func (c *columsMonitor) DingMessage() string {
	tplText := `
**数据库** : {{.TableSchema}} 结构变化

请及时同步数据到 DDL 文档

**新增表**:

{{- range $i, $e := .NewTable}}
* 新表 {{$e.TableName}}
{{- end }}

**删除表**:

{{- range $i, $e := .DelTable}}
* 删表 {{$e.TableName}}
{{- end }}

**新增字段**:

{{- range $i, $e := .NewColumns}}
* 表 {{$e.TableName}} 增 {{$e.ColumnName}}
{{- end }}

**删除字段**:

{{- range $i, $e := .DelColumns}}
* 表 {{$e.TableName}} 删 {{$e.ColumnName}}
{{- end }}

**属性更新**:

{{- range $i, $e := .UpdateColumns}}
* 表 {{$e.TableName}} 改 {{$e.ColumnName}}
{{- end }}
`
	tpl, err := template.New("note").Parse(tplText)
	if err != nil {
		fmt.Printf("failed parse tpltext,err:%s\n", err.Error())
		return ""
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, c)
	if err != nil {
		fmt.Printf("failed execute tpltext,err:%s\n", err.Error())
		return ""
	}

	text := buf.String()

	if len(c.DelColumns) > 0 || len(c.NewColumns) > 0 || len(c.NewTable) > 0 || len(c.DelTable) > 0 {
		var receiver dd.Robot
		receiver.AccessToken = conf.Config.Ding.Token
		receiver.Secret = conf.Config.Ding.Secret
		webHookUrl := receiver.Signature()
		params := receiver.SendMarkdown("数据结更新", text, []string{}, []string{}, false)
		dd.SendRequest(webHookUrl, params)

	}

	return text

}

func newColumnMonitor(schema string) *columsMonitor {

	return &columsMonitor{
		TableSchema: schema,
		NewTable:    make([]tables, 0),
		DelTable:    make([]tables, 0),
		NewColumns:  make([]columns, 0),
		DelColumns:  make([]columns, 0),
	}

}

func main() {

	conf.LoadConfigure(ConfigPath)
	ezgo.LoadComponent(
		ezgo.WithComponentMySQL(ezgo.Default, &conf.Config.SQL),
	)
	ent.InitDB()

	tablePath := conf.Config.SQLMonitor.HistoryTablePath
	columnPath := conf.Config.SQLMonitor.HistoryColumnPath
	tableSchema := conf.Config.SQLMonitor.TableSchema

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

	cm := newColumnMonitor(tableSchema)

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

	// 属性更新

	for _, nc := range hisColumns {
		if isColumnUpdated(nc, currentColumns) {
			cm.AddUpdateColumn(nc.TableName, nc.ColumnName)
		}
	}

	cm.Report()
	s := cm.DingMessage()

	fmt.Println(s)

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

func isColumnUpdated(col columns, slist []columns) bool {

	for _, s := range slist {
		if s.TableName == col.TableName && s.ColumnName == col.ColumnName && reflect.DeepEqual(&col, &s) {
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
	rows, e := ent.DB.QueryContext(ctx, "select TABLE_NAME, COLUMN_NAME, COLUMN_TYPE, COLUMN_DEFAULT, COLUMN_COMMENT FROM COLUMNS WHERE TABLE_SCHEMA=? ", schema)
	if e != nil {
		fmt.Println(e.Error())
		return nil, e
	}

	defer rows.Close()

	cols := make([]columns, 0)

	for rows.Next() {

		var col columns
		if se := rows.Scan(&col.TableName, &col.ColumnName, &col.ColumnType, &col.ColumnDefault, &col.ColumnComment); se != nil {
			fmt.Println(se.Error())
			continue
		} else {
			cols = append(cols, col)
		}
	}
	return cols, nil
}
