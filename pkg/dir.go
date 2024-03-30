package ezgo

import (
	"fmt"
	"os"
	"path/filepath"
)

// contains 检查 @sclie 中是否包含元素 @elem, 如果存在就返回 true
func contains(slice []string, elem string) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

// WorkDir 从 @dir 中遍历 @exts 扩展名的文件, 然后返回全路径的文件名
func WorkDir(dir string, exts []string) []string {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if !info.IsDir() && contains(exts, ext) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("WorkDir: %v\n", err)
	}
	return files
}
