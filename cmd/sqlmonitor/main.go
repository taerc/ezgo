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

	total, err := ent.DB.Cols.Query().Select(cols.FieldCOLUMNNAME, cols.FieldCOLUMNCOMMENT).All(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(len(total))

	for _, v := range total {
		fmt.Println(v.COLUMNNAME, v.COLUMNCOMMENT, v.TABLENAME)
		// fmt.Println(v.FieldTABLESCHEMA, v.FieldTABLENAME)

	}

}
