package ezgo

// ResponseTemplate
import (
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

type ResponseTemplate struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"响应信息"`
}

type GContext struct {
	gin.Context
}

// gin.context typedef 处理一下，后面可以自由替换
// 控制器的处理
type GinFlow struct {
}

type Processor interface {
	PreProc(ctx *gin.Context)
	Proc(ctx *gin.Context)
	PostProc(ctx *gin.Context)
}

// 路由分组的管理
// 模块化路由注册的管理
type GinContext struct {
	Engine    *gin.Engine
	whiteList map[string]bool
	urlRole   map[string]int
}

func NewGinContext() *GinContext {
	gc := new(GinContext)
	gc.Engine = gin.New()
	return gc
}

// Group 如果给了 HandlerFunc 会是怎样的处理流程
func (gc *GinContext) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return gc.Engine.Group(relativePath, handlers...)
}

func (gc *GinContext) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return gc.Engine.Use(middleware...)
}

func (gc *GinContext) Run(ipaddress ...string) error {
	return gc.Engine.Run(ipaddress...)
}

func (gc *GinContext) SetPostProc(route *gin.RouterGroup, relativePath string, processor Processor) {
	gc.SetWhiteList(route.BasePath(), relativePath)
	route.POST(relativePath, processor.PreProc, processor.Proc, processor.PreProc)
}
func (gc *GinContext) SetGetProc(route *gin.RouterGroup, relativePath string, processor Processor) {
	gc.SetWhiteList(route.BasePath(), relativePath)
	route.GET(relativePath, processor.PreProc, processor.Proc, processor.PreProc)
}
func (gc *GinContext) SetWhiteList(basePath, relativePath string) {
	gc.whiteList[filepath.Join(basePath, relativePath)] = true
}

func SetPostProc(route *gin.RouterGroup, relativePath string, processor Processor) {
	route.POST(relativePath, processor.PreProc, processor.Proc, processor.PreProc)
}

func SetGetProc(route *gin.RouterGroup, relativePath string, processor Processor) {
	route.GET(relativePath, processor.PreProc, processor.Proc, processor.PostProc)
}

func (gf *GinFlow) PreProc(ctx *gin.Context) {

}
func (gf *GinFlow) PostProc(ctx *gin.Context) {

}

func (gf *GinFlow) BindJson(ctx *gin.Context, data interface{}) int {
	if err := ctx.BindJSON(data); err != nil {

		if !strings.HasPrefix(err.Error(), "json: invalid use of ,string struct tag,") {
			gf.ResponseJson(ctx, ErrInvalidJsonFormat, nil)
			return ErrInvalidJsonFormat
		} else {
			return Success
		}
	}
	return Success
}
func (gf *GinFlow) ResponseJson(ctx *gin.Context, er int, data interface{}) {
	ctx.JSON(Success, ResponseTemplate{
		Code:    er,
		Data:    data,
		Message: "",
	})
}
//
