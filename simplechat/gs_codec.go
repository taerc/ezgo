package simplechat

import (
	"encoding/binary"
	"fmt"

	"github.com/panjf2000/gnet"
)

const (
	gsFrameDelimiter         uint16 = 0xEB90
	gsRequestFrame           byte   = 0x00
	gsResponseFrame          byte   = 0x01
	gsFrameFixedLength       int    = 15
	gsFrameFixedLengthOffset int    = 11
)

type GSFrame struct {
	StartTag    uint16
	SendSeq     uint32 // inc
	RecvSeq     uint32 // inc
	SessionFlag byte   // 0x00 request 0x01 response
	Length      uint32
	Data        []byte
	EndTag      uint16
}

func (f *GSFrame) debug() {
	fmt.Printf("startTag :%02x\n", f.StartTag)
	fmt.Printf("sendSeq :%02x\n", f.SendSeq)
	fmt.Printf("recvSeq :%02x\n", f.RecvSeq)
	fmt.Printf("session :%x\n", f.SessionFlag)
	fmt.Printf("length :%d\n", f.Length)
	fmt.Printf("endTag :%02x\n", f.EndTag)
}

func (f *GSFrame) Encode(c gnet.Conn, buf []byte) ([]byte, error) {

	buff := make([]byte, 1024)
	f.StartTag = gsFrameDelimiter
	f.EndTag = gsFrameDelimiter
	f.SendSeq = 100
	f.RecvSeq = 101
	f.SessionFlag = gsRequestFrame
	f.Length = uint32(len(buf))

	buff = binary.LittleEndian.AppendUint16(buff, f.StartTag)
	buff = binary.LittleEndian.AppendUint32(buff, f.SendSeq)
	buff = binary.LittleEndian.AppendUint32(buff, f.RecvSeq)
	buff = append(buff, f.SessionFlag)
	buff = binary.LittleEndian.AppendUint32(buff, f.Length)
	buff = append(buff, buf...)
	buff = binary.LittleEndian.AppendUint16(buff, f.EndTag)

	return buff, nil
}

func (f *GSFrame) Decode(c gnet.Conn) ([]byte, error) {

	buf := c.Read()

	f.StartTag = binary.LittleEndian.Uint16(buf)

	if f.StartTag != gsFrameDelimiter {
		fmt.Println("error start ... ")
		return nil, nil
	}

	f.Length = binary.LittleEndian.Uint32(buf[gsFrameFixedLengthOffset:])
	f.EndTag = binary.LittleEndian.Uint16(buf[gsFrameFixedLength+int(f.Length):])
	if f.EndTag != gsFrameDelimiter {
		fmt.Println("error end ... ")
		return nil, nil
	}
	// c.ResetBuffer()
	return buf, nil
}
