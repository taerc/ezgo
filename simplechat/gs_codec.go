package simplechat

import (
	"encoding/binary"
	"fmt"

	"github.com/panjf2000/gnet/v2"
)

const (
	gsFrameDelimiter uint16 = 0xEB90
	gsRequestFrame   byte   = 0x00
	gsResponseFrame  byte   = 0x01
)

type GSFrameCodecConfig struct {
	StartDelimiterOffset int
	SendSequenceOffset   int
	RecvSequenceOffset   int
	FrameTypeOffset      int
	DataLengthOffset     int
	DataOffset           int
	EndDelimiterOffset   int
	FrameDelimiter       uint16
}

func GetGSFrameCodecConfig() *GSFrameCodecConfig {
	return &GSFrameCodecConfig{
		StartDelimiterOffset: 0,
		SendSequenceOffset:   2,
		RecvSequenceOffset:   6,
		FrameTypeOffset:      10,
		DataLengthOffset:     11,
		DataOffset:           15,
		FrameDelimiter:       0xEB90,
	}
}

func NewGSFrameCodec() *GSFrameCodec {
	return &GSFrameCodec{
		config: GetGSFrameCodecConfig(),
	}
}

type GSFrameCodec struct {
	StartTag uint16
	SendSeq  uint32 // inc
	RecvSeq  uint32 // inc
	Type     byte   // 0x00 request 0x01 response
	Length   uint32
	Data     []byte
	EndTag   uint16

	config *GSFrameCodecConfig
}

func (f *GSFrameCodec) debug() {
	fmt.Printf("startTag :%02x\n", f.StartTag)
	fmt.Printf("sendSeq :%02x\n", f.SendSeq)
	fmt.Printf("recvSeq :%02x\n", f.RecvSeq)
	fmt.Printf("session :%x\n", f.Type)
	fmt.Printf("length :%d\n", f.Length)
	fmt.Printf("endTag :%02x\n", f.EndTag)
}

func (f *GSFrameCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {

	buff := make([]byte, 0)
	f.StartTag = gsFrameDelimiter
	f.EndTag = gsFrameDelimiter
	f.SendSeq = 100
	f.RecvSeq = 101
	f.Type = gsRequestFrame
	f.Length = uint32(len(buf))

	buff = binary.LittleEndian.AppendUint16(buff, f.StartTag)
	buff = binary.LittleEndian.AppendUint32(buff, f.SendSeq)
	buff = binary.LittleEndian.AppendUint32(buff, f.RecvSeq)
	buff = append(buff, f.Type)
	buff = binary.LittleEndian.AppendUint32(buff, f.Length)
	buff = append(buff, buf...)
	buff = binary.LittleEndian.AppendUint16(buff, f.EndTag)

	return buff, nil
}

func (f *GSFrameCodec) Decode(c gnet.Conn) ([]byte, error) {

	buf := c.Read()

	f.StartTag = binary.LittleEndian.Uint16(buf)

	if f.StartTag != gsFrameDelimiter {
		fmt.Println("error start ... ")
		return nil, nil
	}

	f.Length = binary.LittleEndian.Uint32(buf[f.config.DataLengthOffset:])
	f.EndTag = binary.LittleEndian.Uint16(buf[f.config.DataOffset+int(f.Length):])
	if f.EndTag != gsFrameDelimiter {
		fmt.Println("error end ... ")
		return nil, nil
	}
	// c.ResetBuffer()
	return buf, nil
}
