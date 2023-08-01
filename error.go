package ezgo

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	OK                  = 0
	Success             = 200
	CodeJsonFormatError = 400 + iota
	CodeInvalidDatetimeFormatString
)

const (
	_             = iota
	CodeErrorBase = 10000 * iota
	CodeParamBase
	CodeDBBase
	CodeAuthBase
	CodePermBase
	CodeResourceBase
	CodeSystemBase
	CodeServiceBase
)

type EM struct {
	C int
	M string
}

func (e EM) Error() string {
	return fmt.Sprintf("%d:%s", e.C, e.M)
}

func NewCode(c int) error {
	return &EM{
		C: c,
		M: GetMessageByCode(c),
	}
}

func NewError(c int, m string) error {
	return &EM{
		C: c,
		M: fmt.Sprintf("%s|%s", GetMessageByCode(c), m),
	}
}
func NewEError(c int, e error) error {
	return &EM{
		C: c,
		M: fmt.Sprintf("%s|%s", GetMessageByCode(c), e.Error()),
	}
}

func GetErrorCode(e error) int {
	es := e.Error()
	if idx := strings.Index(es, ":"); idx != -1 {
		if c, e := strconv.ParseInt(es[0:idx], 10, 64); e == nil {
			return int(c)
		}
	}
	return CodeErrorBase
}

var codeMessage = map[int]string{}

func RegisterCodeMessage(cm map[int]string) int {
	for code, mgs := range cm {

		if _, ok := cm[code]; ok {
			fmt.Println(fmt.Sprintf("WARN error code [%d] is used as [%s]", code, GetMessageByCode(code)))
		}
		cm[code] = mgs
	}
	return 0
}

func GetMessageByCode(c int) string {
	if _, ok := codeMessage[c]; ok {
		return codeMessage[c]
	} else {
		return "!UNK!"
	}
}

func init() {
	codeMessage[Success] = "OK"
	codeMessage[CodeJsonFormatError] = "参数格式错误"
	codeMessage[CodeInvalidDatetimeFormatString] = "无效的日期时间格式"
}
