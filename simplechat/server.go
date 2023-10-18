package simplechat

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/panjf2000/gnet/v2"
	"github.com/taerc/ezgo"
)

// TODO
// package split case

type HandlerFunc interface {
	Handler(c *connection, v interface{})
}

type Encoder interface {
	Marshal(v interface{}) ([]byte, error)
	UnMashal([]byte) (interface{}, error)
}

type tcpServer struct {
	*gnet.BuiltinEventEngine
	ezid            *ezgo.EZID
	decoder         *GSFrameDecoder
	connections     map[string]*connection
	routers         map[string]HandlerFunc
	packetEncoder   Encoder // tcp packet
	businessEncoder Encoder // xml
}

func (es *tcpServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	fmt.Println("onBoot")
	return
}

func (cs *tcpServer) OnShutdown(eng gnet.Engine) {

}

func (es *tcpServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	cid, e := es.ezid.NextStringID()
	if e != nil {
		fmt.Println(e.Error())
	}
	// c.Context()
	ctx := connectionContext{
		Id: cid,
	}
	c.SetContext(ctx)
	return
}

func (es *tcpServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	fmt.Println("close ", c.RemoteAddr().String())

	// TODO
	// valid Login and userId
	ctx := c.Context()
	if ctx != nil {
		ct := ctx.(connectionContext)
		fmt.Println(fmt.Sprintf("login >>usrId :%s connId %s", ct.UsrId, ct.Id))
	}
	return
}

func (cs *tcpServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	cs.decoder.Decode(c)
	return
}

func (es *tcpServer) handlerMessage(cmd *Command, c gnet.Conn) error {

	if cmd.Cmd == CommandLogin {
		return es.commandLogin(cmd, c)
	} else if cmd.Cmd == CommandSend {
		return es.commandSend(cmd, c)
	}

	return nil
}

func (es *tcpServer) commandLogin(cmd *Command, c gnet.Conn) error {

	login := &LoginMessage{}
	e := mapstructure.Decode(cmd.Data, login)
	if e != nil {
		fmt.Println(e.Error())
	}
	// TODO
	// valid Login and userId
	ctx := c.Context()
	if ctx != nil {
		ct := ctx.(connectionContext)
		ct.UsrId = login.UsrId
		c.SetContext(ct)

		usr := &ChatUser{
			Id: ct.UsrId,
			conn: connection{
				Id: ct.Id,
				// conn:     c,
			},
		}
		trackUser(usr)

		fmt.Println(fmt.Sprintf("login >>usrId :%s connId %s", ct.UsrId, ct.Id))
	}

	return nil
}

func (es *tcpServer) commandSend(cmd *Command, c gnet.Conn) error {

	send := &SendMessage{}
	e := mapstructure.Decode(cmd.Data, send)
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println(send.Data)
	ctx := c.Context()
	if ctx != nil {
		ct := ctx.(connectionContext)
		fmt.Println(ct.Id)
		fmt.Println(ct.UsrId)
	}

	if usr, e := getUserById(send.To); e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println(fmt.Sprintf("send >>usrId :%s connId %s", send.To, usr.conn.Id))
		usr.conn.conn.AsyncWrite([]byte(send.Data))
	}
	return nil
}

func StartChatServer(port int) error {
	echo := &tcpServer{
		ezid:    ezgo.NewEZID(0, 0, ezgo.ChatIDSetting()),
		decoder: NewGSFrameDecoder(),
	}
	log.Fatal(gnet.Run(echo, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(false), gnet.WithReusePort(true)))
	return nil
}
