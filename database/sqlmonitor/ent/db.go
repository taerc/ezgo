package ent

import (
	"context"

	"entgo.io/ent/dialect"
	"github.com/taerc/ezgo"
)

var DB *Client = nil

func InitDB() *Client {

	if DB != nil {
		return DB
	}
	drv, e := ezgo.EntDBDriver(ezgo.Default)
	if e != nil {
		return nil
	}
	DB = NewClient(Driver(dialect.DebugWithContext(drv, func(ctx context.Context, a ...any) {
	})))
	return DB

}
