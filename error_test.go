package ezgo

import (
	"fmt"
	"testing"
)

func TestErrorRegister(t *testing.T) {

	ErrorRegister(
		EM{C: Success, M: "OK"},
		EM{C: 401, M: "用户信息错误"},
		EM{C: 502, M: "服务端错误"})

	fmt.Println(GetErrorMessage(401))
	fmt.Println(GetErrorMessage(404))
}
