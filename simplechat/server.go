package simplechat

import (
	"fmt"
	"log"

	"github.com/panjf2000/gnet"
	"github.com/taerc/ezgo"
)

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

	return
}

func (es *chatServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	fmt.Println("close ", c.RemoteAddr().String())
	return
}
func (es *chatServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	// Echo synchronously.
	fmt.Println(string(frame))
	out = frame
	return

	/*
		// Echo asynchronously.
		data := append([]byte{}, frame...)
		go func() {
			time.Sleep(time.Second)
			c.AsyncWrite(data)
		}()
		return
	*/
}

func StartChatServer(port int) error {
	echo := new(chatServer)
	log.Fatal(gnet.Serve(echo, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(false)))
	return nil
}
