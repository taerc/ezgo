package ezgo

import (
	"fmt"
	"testing"
)

func TestDirtyFileWrite(t *testing.T) {
	data := []byte("This is a boy testing")

	DirtyFileWrite("1.txt", data)

	bu, e := DirtyFileRead("1.txt")

	if e == nil {
		fmt.Println(string(bu))
	}
}
