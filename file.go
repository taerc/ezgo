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

// DirtyFileWrite

func DirtyFileWrite(filename string, data []byte) error {

	var bufLen int = 1024
	if len(data) > 960 {
		bufLen = (len(data) + 127) / 64 * 64
	}
	randData := BytesRandom(bufLen)

	idx := int(HashBytesCRC32(randData[bufLen-16:]) & 0x1f)

	for i := 0; i < len(data); i += 1 {
		randData[idx+i] = data[i]
	}
	bL := UInt32Bytes(uint32(len(data)))

	for i := 0; i < 4; i += 1 {
		randData[bufLen-32+i] = bL[i]
	}
	return os.WriteFile(filename, randData, os.ModePerm)
}

// DirtyFileRead

func DirtyFileRead(filePath string) ([]byte, error) {

	if data, e := os.ReadFile(filePath); e == nil {
		idx := int(HashBytesCRC32(data[len(data)-16:]) & 0x1f)
		n := int(BytesUInt32(data[len(data)-32:]))
		return data[idx : idx+n], nil
	} else {
		return nil, e
	}

}
