package ezgo

import (
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"sync"
)

func WithModuleSwagger() func(wg *sync.WaitGroup) {
	return func(wg *sync.WaitGroup) {
		wg.Done()
		r := Group("/docs")
		GET(r, "/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler), ginSwagger.Config{
			Title: "EZGO 接口平台",
		})
		Info(nil, M, "Load swagger Done!")
	}
}
