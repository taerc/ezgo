package simplechat

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/panjf2000/gnet"
	"github.com/taerc/ezgo"
)

// TODO
// package split case

type chatServer struct {
	*gnet.EventServer
	ezid *ezgo.EZID
}

func (es *chatServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Echo server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (es *chatServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
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

func (es *chatServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	fmt.Println("close ", c.RemoteAddr().String())
	return
}

func (es *chatServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	out = frame
	fmt.Println(len(frame))
	cmd := &Command{}
	if e := json.Unmarshal(frame, cmd); e != nil {
		fmt.Println(e.Error())
	}
	go es.handlerMessage(cmd, c)
	// TODO
	// frame decoding
	return
}

func (es *chatServer) handlerMessage(cmd *Command, c gnet.Conn) error {

	if cmd.Cmd == CommandLogin {
		return es.commandLogin(cmd, c)
	} else if cmd.Cmd == CommandSend {

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
		fmt.Println(ct.Id)
		ct.UsrId = login.UsrId
		c.SetContext(ct)
	}

	return nil
}

func (es *chatServer) commandSend(cmd *Command, c gnet.Conn) error {

	send := &SendMessage{}
	e := mapstructure.Decode(cmd.Data, send)
	if e != nil {
		fmt.Println(e.Error())
	}

	return nil
}

func StartChatServer(port int) error {
	echo := &chatServer{
		ezid: ezgo.NewEZID(0, 0, ezgo.ChatIDSetting()),
	}
	log.Fatal(gnet.Serve(echo, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(false)))
	return nil
}
