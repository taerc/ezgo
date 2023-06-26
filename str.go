package ezgo

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

var randomStrSource = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetRandomStr(length int) string {
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63())) //增大随机性
	for i := 0; i < length; i++ {
		result[i] = randomStrSource[r.Intn(len(randomStrSource))]
	}
	return string(result)
}

func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Get16MD5Encode(data string) string {
	return GetMD5Encode(data)[8:24]
}

func StringSplits(str string, seps []string) []string {

	strs := []string{str}
	strsBuf := make([]string, 0)
	for _, sp := range seps {
		for _, s := range strs {
			sa := strings.Split(s, sp)
			strsBuf = append(strsBuf, sa...)
		}
		strs = strsBuf
		strsBuf = append([]string{}) // clear slice
	}
	return strs
}
