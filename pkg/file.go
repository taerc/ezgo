package ezgo

import (
	"os"
)

// 文件操作的基础库都在 os 包里面, 根据需要时再补充
// Mkdirs @description

func Mkdirs(dirPath string) (string, error) {
	if _, e := os.Stat(dirPath); os.IsNotExist(e) { //如果不存在该目录，那么创建该目录
		e = os.MkdirAll(dirPath, os.ModePerm)
		return dirPath, e
	}
	return dirPath, nil
}

// PathExists @description check path

func PathExists(pth string) bool {
	if _, e := os.Stat(pth); os.IsNotExist(e) { //如果不存在该目录，那么创建该目录
		return false
	}
	return true
}
