package ezgo

import (
	"fmt"
	"testing"
)

func TestUInt32Bytes(t *testing.T) {

	b := UInt32Bytes(0x0fffff0f)
	fmt.Println(b)

	c := BytesUInt32(b)
	fmt.Println(c)
}
