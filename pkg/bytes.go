package ezgo

import (
	"math/rand"
	"time"
)

// BytesRandom 生成随机的字节流, 0-255范围
func BytesRandom(length int) []byte {
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
	for i := 0; i < length; i++ {
		result[i] = byte(r.Intn(255))
	}
	return result
}

func UInt32Bytes(u uint32) []byte {
	buf := make([]byte, 4)

	buf[0] = byte(u & 0xff)
	buf[1] = byte(u >> 8 & 0xff)
	buf[2] = byte(u >> 16 & 0xff)
	buf[3] = byte(u >> 24 & 0xff)

	return buf
}

func BytesUInt32(b []byte) uint32 {
	var d uint32 = 0

	d |= uint32(b[0])
	d |= uint32(b[1]) << 8
	d |= uint32(b[2]) << 16
	d |= uint32(b[3]) << 24

	return d
}
