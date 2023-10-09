package simplechat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"

	"github.com/panjf2000/gnet/v2"
)

const (
	gsFrameDelimiter  uint16 = 0xEB90
	gsRequestFrame    byte   = 0x00
	gsResponseFrame   byte   = 0x01
	gsFrameHeaderSize        = 23
)

const (
	GSFRAME_DECODE_STATE_INIT = 0 + iota
	GSFRAME_DECODE_STATE_90
	GSFRAME_DECODE_STATE_EB
	GSFRAME_DECODE_STATE_START
	GSFRAME_DECODE_STATE_FRAME
	GSFRAME_DECODE_STATE_END
	GSFRAME_DECODE_STATE_FRAME_BROKEN
)

type GSFrameHeader struct {
	StartTag uint16
	SendSeq  uint64 // inc
	RecvSeq  uint64 // inc
	Type     byte   // 0x00 request 0x01 response
	Length   uint32
}

func NewGSFrameHeader(sendSeq uint64, recvSeq uint64, ty byte) *GSFrameHeader {
	return &GSFrameHeader{
		StartTag: gsFrameDelimiter,
		SendSeq:  sendSeq,
		RecvSeq:  recvSeq,
		Type:     ty,
		Length:   0,
	}

}

type GSFrameEncoder struct {
}

func (f *GSFrameEncoder) debug() {
	// fmt.Printf("state :%02x\n", f.StartTag)
	// fmt.Printf("sendSeq :%02x\n", f.SendSeq)
	// fmt.Printf("recvSeq :%02x\n", f.RecvSeq)
	// fmt.Printf("session :%x\n", f.Type)
	// fmt.Printf("length :%d\n", f.Length)
	// fmt.Printf("endTag :%02x\n", f.EndTag)
}
func NewGSFrameEncoder() *GSFrameEncoder {
	return &GSFrameEncoder{}
}

func (f *GSFrameEncoder) Encode(header *GSFrameHeader, data []byte) ([]byte, error) {

	buff := &bytes.Buffer{}

	header.Length = uint32(len(data))
	e := binary.Write(buff, binary.LittleEndian, header)
	if e != nil {
		return nil, e
	}
	n, e := buff.Write(data)
	if n != int(header.Length) {
		return nil, errors.New("not match")
	}
	e = binary.Write(buff, binary.LittleEndian, header.StartTag)
	if e != nil {
		return nil, e
	}
	return buff.Bytes(), e
}

type GSFrameDecoder struct {
	header GSFrameHeader
	state  int
	endTag uint16
	mutex  *sync.Mutex
}

func NewGSFrameDecoder() *GSFrameDecoder {
	return &GSFrameDecoder{
		mutex: &sync.Mutex{},
		state: GSFRAME_DECODE_STATE_INIT,
	}
}

func (gs *GSFrameDecoder) Decode(c gnet.Conn) (action gnet.Action) {

	headBuff, e := c.Next(gsFrameHeaderSize)
	if e != nil {
		fmt.Println("decode ", e)
	}
	headRd := bytes.NewReader(headBuff)
	binary.Read(headRd, binary.LittleEndian, &gs.header)
	fmt.Printf("magic %04x\n", gs.header.StartTag)
	fmt.Printf("length %d\n", gs.header.Length)
	if gs.header.StartTag == gsFrameDelimiter {
		data, e := c.Next(int(gs.header.Length + 2))
		if e != nil {
			fmt.Println("rd data ", e.Error())
		}
		gs.endTag = binary.LittleEndian.Uint16(data[gs.header.Length : gs.header.Length+2])
		fmt.Printf("%04x \n", gs.endTag)
		fmt.Println("limit bytes :", c.InboundBuffered())
		// bufio.NewReadWriter()
	}

	return 0
}
