package simplechat

import (
	"fmt"
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

func (c *connection) SendMessage(head PacketHead, v interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	data, e := EncodePacket(head).Marshal(v)
	if e != nil {
		fmt.Println(e)
		return e
	}
	c.conn.Write(data)
	return nil
}

func (c *connection) Send(v interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	data, e := EncodePacket(defaultPacketHead()).Marshal(v)
	if e != nil {
		fmt.Println(e)
		return e
	}
	c.conn.Write(data)
	return nil
}
