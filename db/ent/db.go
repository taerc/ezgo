package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	ezgo "github.com/taerc/ezgo/pkg"
)

var DB *Client = nil

func init() {
	drv, e := ezgo.EntDBDriver(ezgo.Default)
	if e != nil {
		fmt.Println(e)
		return
	}
	DB = NewClient(Driver(dialect.DebugWithContext(drv, func(ctx context.Context, a ...any) {})))
}
