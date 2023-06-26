package ezgo

import (
	"fmt"
	"testing"
)

func TestStringSplits(t *testing.T) {

	s := "this is a book; 中国共产党；你好吗;ss,ss; DDDD ; adkfafla；iskskdidi"
	sa := StringSplits(s, []string{";", "；", ","})
	for _, v := range sa {
		fmt.Println(v)
	}
}
