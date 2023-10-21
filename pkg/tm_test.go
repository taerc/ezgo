package ezgo

import (
	"fmt"
	"testing"
)

func TestUnixToTime(t *testing.T) {
	var ti int64 = 1690803707
	tm := UnixToTime(ti)
	fmt.Println(DateFromTime(tm))
	fmt.Println(DatetimeFromTime(tm))

	if tm, e := TimeParse("2023-07-32 00:00:00"); e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(DateFromTime(tm))

	}

}
