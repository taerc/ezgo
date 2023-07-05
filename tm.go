package ezgo

import "time"

// GetUnixTimestamp @description: 获取当前时间戳
// @return int64
func GetUnixTimeStamp() int64 {
	return time.Now().Unix()
}

func GetLocalDate() string {
	return time.Time{}.Format(time.DateOnly)
}
