package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/taerc/ezgo/conf"
	"github.com/taerc/ezgo/db/ent/columns"
	ezgo "github.com/taerc/ezgo/pkg"

	"github.com/taerc/ezgo/ent"
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

	total, err := ent.DB.Columns.Query().Select(columns.FieldCOLUMNNAME, columns.FieldCOLUMNCOMMENT).All(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range total {
		fmt.Println(v.COLUMNNAME, v.COLUMNCOMMENT, v.TABLENAME)
		// fmt.Println(v.FieldTABLESCHEMA, v.FieldTABLENAME)

	}

}
