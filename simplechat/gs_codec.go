package simplechat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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

func NewGSFrameCodec() *GSFrameCodec {
	return &GSFrameCodec{
		StartTag: gsFrameDelimiter,
	}
}

type GSFrameCodec struct {
	StartTag uint16
	SendSeq  uint64 // inc
	RecvSeq  uint64 // inc
	Type     byte   // 0x00 request 0x01 response
	Length   uint32
}

func (f *GSFrameCodec) debug() {
	fmt.Printf("startTag :%02x\n", f.StartTag)
	fmt.Printf("sendSeq :%02x\n", f.SendSeq)
	fmt.Printf("recvSeq :%02x\n", f.RecvSeq)
	fmt.Printf("session :%x\n", f.Type)
	fmt.Printf("length :%d\n", f.Length)
	// fmt.Printf("endTag :%02x\n", f.EndTag)
}

func (f *GSFrameCodec) Encode(sendSeq uint64, recvSeq uint64, ty byte, data []byte) ([]byte, error) {

	buff := &bytes.Buffer{}
	f.SendSeq = sendSeq
	f.RecvSeq = recvSeq
	f.Type = ty
	f.Length = uint32(len(data))
	e := binary.Write(buff, binary.LittleEndian, f)
	if e != nil {
		return nil, e
	}
	n, e := buff.Write(data)
	if n != int(f.Length) {
		return nil, errors.New("not match")
	}
	e = binary.Write(buff, binary.LittleEndian, f.StartTag)
	if e != nil {
		return nil, e
	}
	return buff.Bytes(), e
}
