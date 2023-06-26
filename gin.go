package ezgo

// ResponseTemplate
import (
	"github.com/gin-gonic/gin"
	"strings"
)

type ResponseTemplate struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"响应信息"`
}

// 控制器的处理
type GinFlow struct {
}

type Processor interface {
	PreProc(ctx *gin.Context)
	Proc(ctx *gin.Context)
	PostProc(ctx *gin.Context)
}

func AddPostProc(route *gin.RouterGroup, relativePath string, processor Processor) {
	route.POST(relativePath, processor.PreProc, processor.Proc, processor.PreProc)
}

func AddGetProc(route *gin.RouterGroup, relativePath string, processor Processor) {
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

// 路由分组的管理
// 模块化路由注册的管理

//
