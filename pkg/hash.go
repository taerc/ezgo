package ezgo

import "hash/crc32"

// StringCRC32

func HashStringCRC32(data string) uint32 {
	return crc32.ChecksumIEEE([]byte(data))
}

// BytesCRC32

func HashBytesCRC32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}
