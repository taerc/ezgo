package ezgo

// steal from https://github.com/gitstliu/go-id-worker/blob/master/idworker.go
import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type IdWorker struct {
	startTime             int64
	workerIdBits          uint
	datacenterIdBits      uint
	maxWorkerId           int64
	maxDatacenterId       int64
	sequenceBits          uint
	workerIdLeftShift     uint
	datacenterIdLeftShift uint
	timestampLeftShift    uint
	sequenceMask          int64
	workerId              int64
	datacenterId          int64
	sequence              int64
	lastTimestamp         int64
	signMask              int64
	idLock                *sync.Mutex
}

func (iw *IdWorker) InitIdWorker(workerId, datacenterId int64) error {

	var baseValue int64 = -1
	iw.startTime = 1463834116272
	iw.workerIdBits = 5
	iw.datacenterIdBits = 5
	iw.maxWorkerId = baseValue ^ (baseValue << iw.workerIdBits)
	iw.maxDatacenterId = baseValue ^ (baseValue << iw.datacenterIdBits)
	iw.sequenceBits = 12
	iw.workerIdLeftShift = iw.sequenceBits
	iw.datacenterIdLeftShift = iw.workerIdBits + iw.workerIdLeftShift
	iw.timestampLeftShift = iw.datacenterIdBits + iw.datacenterIdLeftShift
	iw.sequenceMask = baseValue ^ (baseValue << iw.sequenceBits)
	iw.sequence = 0
	iw.lastTimestamp = -1
	iw.signMask = ^baseValue + 1

	iw.idLock = &sync.Mutex{}

	if iw.workerId < 0 || iw.workerId > iw.maxWorkerId {
		return errors.New(fmt.Sprintf("workerId[%v] is less than 0 or greater than maxWorkerId[%v].", workerId, datacenterId))
	}
	if iw.datacenterId < 0 || iw.datacenterId > iw.maxDatacenterId {
		return errors.New(fmt.Sprintf("datacenterId[%d] is less than 0 or greater than maxDatacenterId[%d].", workerId, datacenterId))
	}
	iw.workerId = workerId
	iw.datacenterId = datacenterId
	return nil
}

func (iw *IdWorker) NextId() (int64, error) {
	iw.idLock.Lock()
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

	iw.idLock.Unlock()

	id := ((timestamp - iw.startTime) << iw.timestampLeftShift) |
		(iw.datacenterId << iw.datacenterIdLeftShift) |
		(iw.workerId << iw.workerIdLeftShift) |
		iw.sequence

	if id < 0 {
		id = -id
	}

	return id, nil
}

func (iw *IdWorker) tilNextMillis() int64 {
	timestamp := time.Now().UnixNano()
	if timestamp <= iw.lastTimestamp {
		timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return timestamp
}
