package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/taerc/ezgo/conf"
	"github.com/taerc/ezgo/ent/columns"
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

	test()

}

func test() {
	ctx := context.Background()

	total, err := ent.DB().Columns().Query().Select(columns.FieldTABLESCHEMA, columns.FieldCOLUMNNAME, columns.FieldCOLUMNDEFAULT, columns.FieldCOLUMNCOMMENT).All(ctx)

	for _, v := range total {
		fmt.Println(v.COLUMN_NAME, v.COLUMN_DEFAULT, v.TABLE_NAME)
		// fmt.Println(v.FieldTABLESCHEMA, v.FieldTABLENAME)

	}

}
