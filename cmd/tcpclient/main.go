package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/panjf2000/gnet/v2"
	"github.com/taerc/ezgo/simplechat"
)

// @user message
// Q
// ^group message
// message

type commandProto struct {
	UsrId   string
	GroupId string
	Cmd     string
	Type    string
	Data    string
}

func (c *commandProto) Scan(input string) error {
	c.clean()
	trimSpace := strings.TrimSpace(input)

	if trimSpace == "Q" {
		c.Cmd = "Q"
	} else {
		c.Cmd = "Message"
	}

	if strings.HasPrefix(trimSpace, "@") {
		c.Type = "User"
		idx := strings.Index(trimSpace, " ")
		if idx != -1 {
			c.Data = trimSpace[idx:]
			c.UsrId = trimSpace[1:idx]
		} else {
			c.UsrId = trimSpace[1:]
		}

	} else if strings.HasPrefix(trimSpace, "^") {
		c.Type = "Group"
		idx := strings.Index(trimSpace, " ")
		if idx != -1 {
			c.Data = trimSpace[idx:]
			c.GroupId = trimSpace[1:idx]
		} else {
			c.UsrId = trimSpace[1:]
		}
	} else {
		c.Type = "invalid"
	}

	return nil
}

func (c *commandProto) clean() {
	c.Cmd = ""
	c.Data = ""
	c.UsrId = ""
	c.GroupId = ""
	c.Type = ""
}

func (c *commandProto) Debug() {
	fmt.Println("usrId : ", c.UsrId)
	fmt.Println("groupId : ", c.GroupId)
	fmt.Println("command : ", c.Cmd)
	fmt.Println("type : ", c.Type)
	fmt.Println("data : ", c.Data)
}

func login(usrId string) (net.Conn, error) {

	fmt.Println("chat server testing")
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return nil, err
	}
	header := simplechat.NewPacketHead(1, 2, 1)
	frame := simplechat.EncodePacket(*header)

	login := simplechat.LoginMessage{
		UsrId: usrId,
	}
	cmd := simplechat.NewCommand(simplechat.CommandLogin, login)
	byteData, e := frame.Marshal(cmd)
	// byteData, e := json.Marshal(cmd)
	if e != nil {
		fmt.Println(e.Error())
	}
	// fmt.Println(len(byteData))

	n, e := conn.Write(byteData)
	fmt.Println("bytes ", n)

	return conn, nil
}

func connect() (net.Conn, error) {
	return net.Dial("tcp", "localhost:9999")
}

func pureTCPConnection() {

	var usrId string
	//项目配置文件路径获取
	flag.StringVar(&usrId, "user", "wangfangming", "path of configure file.")
	flag.Parse()

	// conn, e := connect()
	conn, e := login(usrId)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	defer conn.Close()

	inputReader := bufio.NewReader(os.Stdin)
	cmd := &commandProto{}
	time.Sleep(2 * time.Second)

	header := simplechat.NewPacketHead(1, 2, 1)
	frame := simplechat.EncodePacket(*header)

	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("read from console failed, err: %v\n", err)
			break
		}
		cmd.Scan(input)

		if cmd.Cmd == "Q" {
			break
		}

		if cmd.Type == "invalid" || cmd.Type == "Group" {
			fmt.Println("only support to user")
			continue
		}

		cmd.Debug()

		send := simplechat.SendMessage{
			From: usrId,
			To:   cmd.UsrId,
			Data: cmd.Data,
		}

		go func() {
			time.Sleep(2 * time.Second)
			for {
				recvMessage := make([]byte, 1024)
				n, _ := conn.Read(recvMessage)
				fmt.Println("bytes ", n)
				fmt.Println(string(recvMessage))
				time.Sleep(1 * time.Second)
			}
		}()

		message := simplechat.NewCommand(simplechat.CommandSend, send)
		byteMessage, e := json.Marshal(message)
		if e != nil {
			fmt.Println(e.Error())
		}
		fmt.Println(len(byteMessage))
		packet, e := frame.Marshal(message)
		dirtyData := append(packet, packet...)
		n, e := conn.Write(dirtyData)
		if e != nil {
			fmt.Println(e.Error())
		}
		fmt.Println("bytes ", n)
	}

}

type clientEvents struct {
	*gnet.BuiltinEventEngine
}

func (ev *clientEvents) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	// c.SetContext([]byte{})
	// rspCh := make(chan []byte, 1)
	// ev.rspChMap.Store(c.LocalAddr().String(), rspCh)
	return nil, gnet.None
}

func (ev *clientEvents) OnClose(c gnet.Conn, err error) gnet.Action {
	// if ev.svr != nil {
	// 	if atomic.AddInt32(&ev.svr.clientActive, -1) == 0 {
	// 		return Shutdown
	// 	}
	// }
	return gnet.None
}

func (ev *clientEvents) React(packet []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Println("")
	// ctx := c.Context()
	// var p []byte
	// if ctx != nil {
	// 	p = ctx.([]byte)
	// } else { // UDP
	// 	ev.packetLen = 1024
	// }
	// p = append(p, packet...)
	// if len(p) < ev.packetLen {
	// 	c.SetContext(p)
	// 	return
	// }
	// v, _ := ev.rspChMap.Load(c.LocalAddr().String())
	// rspCh := v.(chan []byte)
	// rspCh <- p
	// c.SetContext([]byte{})
	return
}

func gnetTcpClientConnection() {

	// codec := simplechat.NewGSFrameCodec()
	cve := new(clientEvents)
	client, e := gnet.NewClient(cve)

	if e != nil {
		fmt.Println(e.Error())
	}
	conn, e := client.Dial("tcp", "127.0.0.1:9999")
	client.Start()
	defer client.Stop()

	for {
		conn.AsyncWrite([]byte("This is gnet/v2 "), func(c gnet.Conn, err error) error {
			// fmt.Println("sync %v", err)
			return err
		})
		time.Sleep(1 * time.Second)

	}

}

func main() {
	// gnetTcpClientConnection()
	pureTCPConnection()

	// fmt.Printf("%c", 0xEB)
	// fmt.Printf("%c", 0x90)

}
