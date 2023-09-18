package ezgo

import (
	"fmt"
	"strconv"
	"strings"
)

const (
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

type zerror struct {
	c int
	m string
}

func (e zerror) Error() string {
	return fmt.Sprintf("%d:%s", e.c, e.m)
}

func NewCode(c int) error {
	return &zerror{
		c: c,
		m: GetMessageByCode(c),
	}
}

func NewError(c int, m string) error {
	return &zerror{
		c: c,
		m: fmt.Sprintf("%s|%s", GetMessageByCode(c), m),
	}
}
func NewEError(c int, e error) error {
	return &zerror{
		c: c,
		m: fmt.Sprintf("%s|%s", GetMessageByCode(c), e.Error()),
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
