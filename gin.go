package ezgo

// Response
import (
	"github.com/gin-gonic/gin"
	"strings"
)

type Response struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message"`
	RequestId string      `json:"request_id"`
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
	ctx.JSON(Success, Response{
		Code:      er,
		Data:      data,
		Message:   "",
		RequestId: GetRequestId(ctx),
	})
}

//
