package simplechat

import (
	"fmt"
	"log"

	"github.com/panjf2000/gnet/v2"
	"github.com/taerc/ezgo"
)

// TODO
// package split case

type HandlerFunc = func(c *connection, hd PacketHead, data interface{})

type tcpServer struct {
	*gnet.BuiltinEventEngine
	ezid          *ezgo.EZID
	connections   map[string]*connection
	connectionNum int
	paser         *PacketParser
	chat          *chatServer
	// packetEncoder   Encoder // tcp packet
	// businessEncoder Encoder // xml
}

func (ts *tcpServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	fmt.Println("onBoot")
	return
}

func (ts *tcpServer) OnShutdown(eng gnet.Engine) {
	fmt.Println("shutdown")
}

func (ts *tcpServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	cid, e := ts.ezid.NextStringID()
	if e != nil {
		fmt.Println(e.Error())
	}
	ctx := connectionContext{
		Id: cid,
	}
	c.SetContext(ctx)
	ts.connections[cid] = newConnection(cid, c)
	ts.connectionNum++
	return
}

func (ts *tcpServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	fmt.Println("close ", c.RemoteAddr().String())

	// TODO
	// valid Login and userId
	ctx := c.Context()
	if ctx != nil {
		ct := ctx.(connectionContext)
		fmt.Println(fmt.Sprintf("connId %s", ct.Id))
	}
	// close socket
	// close lock
	ts.connectionNum--
	return
}

func (ts *tcpServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	ctx := c.Context()
	if ctx != nil {
		ct := ctx.(connectionContext)
		conn := ts.connections[ct.Id]
		conn.conn = c
		ts.paser.Parse(conn)
	}
	return
}

func (ts *tcpServer) RegisterCommand(cmd string, handler HandlerFunc) {
	ts.chat.RegisterCommand(cmd, handler)
}

func StartChatServer(port int) error {
	echo := &tcpServer{
		ezid:        ezgo.NewEZID(0, 0, ezgo.ChatIDSetting()),
		paser:       NewPacketParser(),
		connections: make(map[string]*connection),
		chat:        newChatServer(),
	}
	log.Fatal(gnet.Run(echo, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(false), gnet.WithReusePort(true)))
	return nil
}
