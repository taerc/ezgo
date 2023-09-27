package simplechat

import (
	"sync"

	"github.com/panjf2000/gnet"
)

type connection struct {
	Id       string
	conn     gnet.Conn
	connLock *sync.Mutex
}

type connectionContext struct {
	Id    string
	UsrId string
}

func (c *connection) SendMessage(data string) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	// fmt.Println(c.connId, data)
}
