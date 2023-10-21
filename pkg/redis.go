package ezgo

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/taerc/ezgo/conf"
	"sync"
)

// redis pool
var redisMap sync.Map

func initRedis(name string, conf *conf.RedisConf) error {

	// redis://<user>:<pass>@localhost:6379/<db>
	opts, e := redis.ParseURL(conf.RedisDSN)
	if e != nil {
		return e
	}
	rdb := redis.NewClient(opts)
	redisMap.Store(name, rdb)
	return nil
}

func REDIS(name ...string) *redis.Client {

	db := Default
	if len(name) != 0 {
		db = name[0]
	}
	if v, ok := redisMap.Load(db); !ok {
		return nil
	} else {
		return v.(*redis.Client)
	}
}

func WithComponentRedis(name string, c *conf.RedisConf) Component {
	return func(wg *sync.WaitGroup) {
		wg.Done()
		initRedis(name, c)
		Info(nil, M, fmt.Sprintf("Finished Load [%s]-REDIS", name))
	}
}