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

func ChatIDSetting() EZIDSetting {
	return EZIDSetting{
		groupIdWidth:  8,
		objectIdWidth: 8,
		sequenceWidth: 12,
	}
}

func (iw *EZID) NewEZID(objectId, groupId int64, setting EZIDSetting) error {

	var baseValue int64 = -1
	iw.startTime = 1463834116272
	iw.setting = setting
	iw.setting.objectIdWidth = 5
	iw.setting.groupIdWidth = 5
	iw.setting.sequenceWidth = 12
	iw.objectIdLimit = baseValue ^ (baseValue << iw.setting.objectIdWidth)
	iw.groupIdLimit = baseValue ^ (baseValue << iw.setting.groupIdWidth)
	iw.objectIdOffset = iw.setting.sequenceWidth
	iw.groupIdOffset = iw.setting.objectIdWidth + iw.objectIdOffset
	iw.timestampOffset = iw.setting.groupIdWidth + iw.groupIdOffset
	iw.sequenceMask = baseValue ^ (baseValue << iw.setting.sequenceWidth)
	iw.sequence = 0
	iw.lastTimestamp = -1
	iw.signMask = ^baseValue + 1
	iw.idMutex = &sync.Mutex{}

	if iw.objectId < 0 || iw.objectId > iw.objectIdLimit {
		return errors.New(fmt.Sprintf("objectId[%v] is less than 0 or greater than objectIdLimit[%v].", objectId, groupId))
	}
	if iw.groupId < 0 || iw.groupId > iw.groupIdLimit {
		return errors.New(fmt.Sprintf("groupId[%d] is less than 0 or greater than groupIdLimit[%d].", objectId, groupId))
	}
	iw.objectId = objectId
	iw.groupId = groupId
	return nil
}

func (iw *EZID) NextId() (int64, error) {
	iw.idMutex.Lock()
	defer iw.idMutex.Unlock()

	timestamp := time.Now().UnixNano()
	if timestamp < iw.lastTimestamp {
		return -1, errors.New(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", iw.lastTimestamp-timestamp))
	}

	if timestamp == iw.lastTimestamp {
		iw.sequence = (iw.sequence + 1) & iw.sequenceMask
		if iw.sequence == 0 {
			timestamp = iw.tilNextMillis()
			iw.sequence = 0
		}
	} else {
		iw.sequence = 0
	}

	iw.lastTimestamp = timestamp

	id := ((timestamp - iw.startTime) << iw.timestampOffset) |
		(iw.groupId << iw.groupIdOffset) |
		(iw.objectId << iw.objectIdOffset) |
		iw.sequence

	if id < 0 {
		id = -id
	}

	return id, nil
}

func (iw *EZID) NextStringID() (string, error) {
	id, e := iw.NextId()
	if e != nil {
		return "", e
	}
	sid := strconv.FormatInt(id, 10)
	return sid, nil
}

func (iw *EZID) tilNextMillis() int64 {
	timestamp := time.Now().UnixNano()
	if timestamp <= iw.lastTimestamp {
		timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return timestamp
}
