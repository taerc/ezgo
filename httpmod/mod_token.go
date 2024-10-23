package httpmod

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	ezgo "github.com/taerc/ezgo/pkg"
)

type authService struct {
}

type authRequest struct {
	DeviceId string `json:"device_id"`
}

type authResponse struct {
	Auth string `json:"auth"`
}

func (a *authService) Auth(ctx *gin.Context) {

	ar := &authRequest{}

	if e := JsonBind(ctx, ar); e != nil {
		ErrorResponse(ctx, e)
		return
	}
	devs := []string{"5YSZM3200310NX", "5YSZL8500330MZ"}

	ares := &authResponse{Auth: ""}
	if ar.DeviceId == devs[0] || ar.DeviceId == devs[1] {
		ares.Auth = ezgo.SHA256(ar.DeviceId)
		OKResponse(ctx, ares)
		return
	}
	ErrorResponse(ctx, errors.New("不合法的序列号"))
	return

}
func (a *authService) GetAuth(ctx *gin.Context) {

	devId := ctx.Query("device_id")
	devs := []string{"5YSZM3200310NX", "5YSZL8500330MZ"}

	ares := &authResponse{Auth: ""}
	if devId == devs[0] || devId == devs[1] {
		ares.Auth = ezgo.SHA256(devId)
		OKResponse(ctx, ares)
		return
	}
	ErrorResponse(ctx, errors.New("不合法的序列号"))
	return

}
func WithModuleToken() func(wg *sync.WaitGroup) {
	return func(wg *sync.WaitGroup) {
		defer wg.Done()
		s := new(authService)
		route := Group("/api/token/")
		POST(route, "/auth", s.Auth)
		GET(route, "/auth", s.GetAuth)
	}
}
