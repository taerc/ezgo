package simplechat

import (
	"fmt"
	"log"

	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/panjf2000/gnet/v2"
	ezgo "github.com/taerc/ezgo/pkg"
)

type ChatServer struct {
	commands map[string]HandlerFunc
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		commands: make(map[string]func(c *connection, hd PacketHead, data interface{})),
	}
}

func StartChatServer(port int) error {

	chatServer := NewChatServer()

	chatServer.RegisterCommand(CommandLogin, cmdLogin)
	chatServer.RegisterCommand(CommandLogout, cmdLogout)
	chatServer.RegisterCommand(CommandSend, cmdSendMsg)

	echo := &tcpServer{
		ezid:        ezgo.NewEZID(0, 0, ezgo.ChatIDSetting()),
		connections: make(map[string]*connection),
		parser:      NewPacketParser(chatServer),
	}
	log.Fatal(gnet.Run(echo, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(false), gnet.WithReusePort(true)))
	return nil
}

func (cs *ChatServer) HandlerPacket(conn *connection, hd PacketHead, packet []byte) error {

	cmd := &Command{}
	if e := json.Unmarshal(packet, &cmd); e != nil {
		fmt.Println(e)
		return e
	}
	if def, ok := cs.commands[cmd.Cmd]; ok {
		def(conn, hd, cmd.Data)
	}
	return nil
}

func (cs *ChatServer) RegisterCommand(cmd string, handler HandlerFunc) {
	cs.commands[cmd] = handler
}

func cmdLogin(c *connection, hd PacketHead, data interface{}) {
	login := &LoginMessage{}
	e := mapstructure.Decode(data, login)
	if e != nil {
		fmt.Println(e.Error())
	}
	trackUser(&ChatUser{Id: login.UsrId, conn: c})
	fmt.Println(fmt.Sprintf("login >>usrId %s connId %s", c.UsrId, c.Id))

}

func cmdLogout(c *connection, hd PacketHead, data interface{}) {

}

func cmdSendMsg(c *connection, hd PacketHead, data interface{}) {
	send := &SendMessage{}
	e := mapstructure.Decode(data, send)
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println(send.Data)
	if usr, e := getUserById(send.To); e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println(fmt.Sprintf("send >>usrId :%s connId %s", send.To, usr.conn.Id))
		usr.conn.SendMessage(hd, send)
	}
}
