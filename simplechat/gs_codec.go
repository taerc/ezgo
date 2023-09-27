package simplechat

import (
	"encoding/binary"
	"errors"

	"github.com/panjf2000/gnet"
)

var gsFrameDelimiter uint16 = 0xEB90
var gsRequestFrame byte = 0x00
var gsResponseFrame byte = 0x01
var gsFrameFixedLength = 17

type GSFrame struct {
	StartTag    uint16
	SendSeq     uint32 // inc
	RecvSeq     uint32 // inc
	SessionFlag byte   // 0x00 request 0x01 response
	Length      uint32
	Data        []byte
	EndTag      uint16
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

	if len(buf) == 0 {
		return nil, errors.New("incomplete packet")
	}

	c.ResetBuffer()

	return buf, nil
}
