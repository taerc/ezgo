package ezgo

// steal from https://github.com/gitstliu/go-id-object/blob/master/idobject.go
// base from snowflake
import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// timestamp | group | object | seq

type EZIDSetting struct {
	groupIdWidth  uint
	objectIdWidth uint
	sequenceWidth uint
}

type EZID struct {
	// part
	lastTimestamp int64
	groupId       int64
	objectId      int64
	sequence      int64

	startTime       int64
	objectIdLimit   int64
	groupIdLimit    int64
	timestampOffset uint
	groupIdOffset   uint
	objectIdOffset  uint
	sequenceMask    int64
	signMask        int64

	setting EZIDSetting
	idMutex *sync.Mutex
}

func DataCenterSetting() EZIDSetting {
	return EZIDSetting{
		groupIdWidth:  5,
		objectIdWidth: 5,
		sequenceWidth: 12,
	}
}

// ChaitIDSetting
// 聊天服务器用于生成 connection id
func ChatIDSetting() EZIDSetting {
	return EZIDSetting{
		groupIdWidth:  8,
		objectIdWidth: 8,
		sequenceWidth: 12,
	}
}

func NewEZID(groupId, objectId int64, setting EZIDSetting) *EZID {
	var maxLimit int64 = -1
	ezid := &EZID{
		startTime:      time.Now().UnixNano(),
		setting:        setting,
		objectIdOffset: setting.sequenceWidth,
		sequence:       0,
		lastTimestamp:  -1,
		signMask:       ^maxLimit + 1,
		idMutex:        &sync.Mutex{},
		objectId:       objectId,
		groupId:        groupId,
	}
	ezid.objectIdLimit = maxLimit ^ (maxLimit << int64(setting.objectIdWidth))
	ezid.groupIdLimit = maxLimit ^ (maxLimit << int64(setting.groupIdWidth))
	ezid.sequenceMask = maxLimit ^ (maxLimit << int64(setting.sequenceWidth))
	ezid.groupIdOffset = setting.objectIdWidth + ezid.objectIdOffset
	ezid.timestampOffset = setting.groupIdWidth + ezid.groupIdOffset

	if ezid.objectId < 0 || ezid.objectId > ezid.objectIdLimit {
		// return errors.New(fmt.Sprintf("objectId[%v] is less than 0 or greater than objectIdLimit[%v].", objectId, groupId))
		return nil
	}
	if ezid.groupId < 0 || ezid.groupId > ezid.groupIdLimit {
		// return errors.New(fmt.Sprintf("groupId[%d] is less than 0 or greater than groupIdLimit[%d].", objectId, groupId))
		return nil
	}
	return ezid
}

func (ezid *EZID) NextId() (int64, error) {
	ezid.idMutex.Lock()
	defer ezid.idMutex.Unlock()

	timestamp := time.Now().UnixNano()
	if timestamp < ezid.lastTimestamp {
		return -1, errors.New(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", ezid.lastTimestamp-timestamp))
	}

	if timestamp == ezid.lastTimestamp {
		ezid.sequence = (ezid.sequence + 1) & ezid.sequenceMask
		if ezid.sequence == 0 {
			timestamp = ezid.tilNextMillis()
			ezid.sequence = 0
		}
	} else {
		ezid.sequence = 0
	}

	ezid.lastTimestamp = timestamp

	id := ((timestamp - ezid.startTime) << ezid.timestampOffset) |
		(ezid.groupId << ezid.groupIdOffset) |
		(ezid.objectId << ezid.objectIdOffset) |
		ezid.sequence

	if id < 0 {
		id = -id
	}

	return id, nil
}

func (ezid *EZID) NextStringID() (string, error) {
	id, e := ezid.NextId()
	if e != nil {
		return "", e
	}
	sid := strconv.FormatInt(id, 10)
	return sid, nil
}

func (ezid *EZID) tilNextMillis() int64 {
	timestamp := time.Now().UnixNano()
	if timestamp <= ezid.lastTimestamp {
		timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return timestamp
}
