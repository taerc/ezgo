package ezgo

import (
	"encoding/json"
	"os"
)

// LoadJson @description: Load @jsonPath to @data
//
func LoadJson(jsonPath string, data interface{}) (int, error) {
	fd, e := os.Open(jsonPath)
	if e != nil {
		return Failed, e
	}
	defer func() {
		fd.Close()
	}()

	dec := json.NewDecoder(fd)
	if e = dec.Decode(&data); e != nil {
		return Failed, e
	}
	return OK, nil
}

// SaveJson @description: Save @data to @dstPath
func SaveJson(dstPath string, data interface{}) (int, error) {

	// 创建文件
	fd, e := os.Create(dstPath)
	if e != nil {
		return Failed, e
	}
	defer func() {
		fd.Close()
	}()

	// 带JSON缩进格式写文件
	if data, ei := json.MarshalIndent(data, "  ", "  "); e != nil {
		return Failed, ei
	} else {
		if _, e = fd.Write(data); e != nil {
			return Failed, e
		}
	}
	return OK, nil
}
