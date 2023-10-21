package ezgo

import (
	"math/rand"
	"strings"
	"time"
)

var randomStrSource = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// StringRandom @description: generate string with

func StringRandom(length int) string {
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
	for i := 0; i < length; i++ {
		result[i] = randomStrSource[r.Intn(len(randomStrSource))]
	}
	return string(result)
}

// StringSplits split @str by @seps to a string array
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
