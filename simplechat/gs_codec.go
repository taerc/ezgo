package simplechat

import "github.com/panjf2000/gnet"

type GSFrame struct {
	StartTag    string // EB90
	SendSeq     int    // inc
	RecvSeq     int    // inc
	SessionFlag byte   // 0x00 request 0x01 response
	Length      int
	Data        string
	EndTag      string //EB90
}

func (f *GSFrame) Encode(c gnet.Conn) {

}

func (f *GSFrame) Decode() {

}
