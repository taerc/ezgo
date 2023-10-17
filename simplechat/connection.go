package simplechat

import (
	"sync"

	"github.com/panjf2000/gnet"
)

// type Request struct {
// 	ProtocolHeader
// 	body []byte
// }

type Request interface {
	Head() ([]byte, error)
	Body() ([]byte, error)
}

type connection struct {
	Id         string
	lastSendId uint64
	lastRecvId uint64
	conn       gnet.Conn
	lock       *sync.Mutex
	request    Request
}

func newConnection(id string, conn gnet.Conn) *connection {

	return &connection{
		Id:         id,
		lastSendId: 0,
		lastRecvId: 0,
		conn:       conn,
		lock:       &sync.Mutex{},
	}

}

type connectionContext struct {
	Id    string
	UsrId string
}

func (c *connection) SendMessage(data string) {
	// c.connLock.Lock()
	// defer c.connLock.Unlock()
}
