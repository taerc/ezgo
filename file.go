package ezgo

import "os"

// Mkdirs @description

func Mkdirs(dirPath string) (string, error) {
	if _, e := os.Stat(dirPath); os.IsNotExist(e) { //如果不存在该目录，那么创建该目录
		e = os.MkdirAll(dirPath, os.ModePerm)
		return dirPath, e
	}
	return dirPath, nil
}
