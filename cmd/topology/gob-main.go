package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"
)

func load(fd string) map[string]int {

	// 打开文件
	file, err := os.Open(fd)
	if err != nil {
		fmt.Println("Error opening file:", err)
		data := make(map[string]int, 256)
		return data
	}
	defer file.Close()

	// 创建 Gob 解码器
	decoder := gob.NewDecoder(file)

	// 创建一个空的 map 对象，用于存储解码后的数据
	var data map[string]int

	// 解码数据
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding data:", err)
		data := make(map[string]int, 256)
		return data
	}
	return data
}

var DataPath string
func init() {
	flag.StringVar(&DataPath, "p", "xxx.gob", "gob")
	flag.Parse()
}
func main() {

	data := load(DataPath)

	for k, v := range data {
		fmt.Println(k, v)
	}

}