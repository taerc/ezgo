package simplechat

import (
	"sync"

	"github.com/panjf2000/gnet/v2"
)

type Request interface {
	Head() ([]byte, error)
	Body() ([]byte, error)
}

type connection struct {
	Id         string
	UsrId      string // for demo chatting
	conn       gnet.Conn
	lastSendId uint64
	lastRecvId uint64
	request    Request
	lock       *sync.Mutex
}

func newConnection(id string, conn gnet.Conn) *connection {

	return &connection{
		Id:         id,
		UsrId:      "",
		lastSendId: 0,
		lastRecvId: 0,
		conn:       conn,
		lock:       &sync.Mutex{},
	}
}

type connectionContext struct {
	Id string
}

func (c *connection) SendMessage(data string) {
	// c.connLock.Lock()
	// defer c.connLock.Unlock()
}
