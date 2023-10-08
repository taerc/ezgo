package simplechat

import (
	"fmt"
	"log"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/panjf2000/gnet/v2"
	"github.com/taerc/ezgo"
)

// TODO
// package split case

type chatServer struct {
	*gnet.BuiltinEventEngine
	ezid *ezgo.EZID
}

func (es *chatServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	fmt.Println("data\n")
	// log.Printf("Echo server is listening on %s (multi-cores: %t, loops: %d)\n",
	// 	srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (es *chatServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Println("open ")
	cid, e := es.ezid.NextStringID()
	if e != nil {
		fmt.Println(e.Error())
	}
	// c.Context()
	ctx := connectionContext{
		Id: cid,
	}
	fmt.Println(cid)
	c.SetContext(ctx)
	return
}

func (es *chatServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
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

func (es *chatServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	// cmd := &Command{}
	// if e := json.Unmarshal(frame, cmd); e != nil {
	// 	fmt.Println(e.Error())
	// }
	return
}

func (es *chatServer) handlerMessage(cmd *Command, c gnet.Conn) error {

	if cmd.Cmd == CommandLogin {
		return es.commandLogin(cmd, c)
	} else if cmd.Cmd == CommandSend {
		return es.commandSend(cmd, c)
	}

	return nil
}

func (es *chatServer) commandLogin(cmd *Command, c gnet.Conn) error {

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
				Id:       ct.Id,
				conn:     c,
				connLock: &sync.Mutex{},
			},
		}
		trackUser(usr)

		fmt.Println(fmt.Sprintf("login >>usrId :%s connId %s", ct.UsrId, ct.Id))
	}

	return nil
}

func (es *chatServer) commandSend(cmd *Command, c gnet.Conn) error {

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

	// c.SendTo([]byte("OK"))

	return nil
}

func StartChatServer(port int) error {
	echo := &chatServer{
		ezid: ezgo.NewEZID(0, 0, ezgo.ChatIDSetting()),
	}
	log.Fatal(gnet.Serve(echo, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(false)))
	return nil
}
