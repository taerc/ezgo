package ezgo

import "time"

// GetUnixTimestamp @description: 获取当前时间戳
// @return int64
func GetUnixTimeStamp() int64 {
	return time.Now().Unix()
}

func UnixToTime(t int64) time.Time {
	return time.Unix(t, 0)
}

func DateFromTime(t time.Time) string {
	return t.Format(time.DateOnly)
}

func DatetimeFromTime(t time.Time) string {
	return t.Format(time.DateTime)
}

func TimeParse(strFmtTime string) (time.Time, error) {
	if t, e := time.Parse(time.DateTime, strFmtTime); e != nil {
		return time.Time{}, NewEError(CodeInvalidDatetimeFormatString, e)
	} else {
		return t, nil
	}
}

func GetLocalDate() string {
	return time.Now().Format(time.DateOnly)
}
