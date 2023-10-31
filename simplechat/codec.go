package simplechat

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

type Encoder interface {
	Marshal(v interface{}) ([]byte, error)
}
type Decoder interface {
	UnMashal([]byte) (interface{}, error)
}

const (
	ppFrameDelimiter  uint16 = 0xEB90
	ppRequestFrame    byte   = 0x00
	ppResponseFrame   byte   = 0x01
	ppFrameHeaderSize int    = 23
)

const (
	ppFRAME_DECODE_STATE_INIT = 0 + iota
	ppFRAME_DECODE_STATE_90
	ppFRAME_DECODE_STATE_EB
	ppFRAME_DECODE_STATE_START
	ppFRAME_DECODE_STATE_FRAME
	ppFRAME_DECODE_STATE_END
	ppFRAME_DECODE_STATE_FRAME_BROKEN
)

type PacketHead struct {
	StartTag uint16
	SendSeq  uint64 // inc
	RecvSeq  uint64 // inc
	Type     byte   // 0x00 request 0x01 response
	Length   uint32
}

func defaultPacketHead() PacketHead {
	return PacketHead{
		StartTag: ppFrameDelimiter,
		SendSeq:  0,
		RecvSeq:  0,
		Type:     0,
		Length:   0,
	}

}

type PacketReactor interface {
	HandlerPacket(conn *connection, header PacketHead, packet []byte) error
}

func NewPacketHead(sendSeq uint64, recvSeq uint64, ty byte) *PacketHead {
	return &PacketHead{
		StartTag: ppFrameDelimiter,
		SendSeq:  sendSeq,
		RecvSeq:  recvSeq,
		Type:     ty,
		Length:   0,
	}
}

type packetEncoder struct {
	head PacketHead
}

func encodePacket(head PacketHead) *packetEncoder {
	return &packetEncoder{head: head}
}

func (pe *packetEncoder) Marshal(v interface{}) ([]byte, error) {

	buff := &bytes.Buffer{}
	data, e := json.Marshal(v)
	if e != nil {
		return nil, e
	}

	pe.head.Length = uint32(len(data))
	e = binary.Write(buff, binary.LittleEndian, pe.head)
	if e != nil {
		return nil, e
	}
	n, e := buff.Write(data)
	if n != int(pe.head.Length) {
		return nil, errors.New("not match")
	}
	e = binary.Write(buff, binary.LittleEndian, pe.head.StartTag)
	if e != nil {
		return nil, e
	}
	return buff.Bytes(), e
}

type PacketParser struct {
	header      PacketHead
	headerBytes [23]byte
	state       int
	endTag      uint16
	handler     PacketReactor
}

func NewPacketParser(hdlr PacketReactor) *PacketParser {
	return &PacketParser{
		state:   ppFRAME_DECODE_STATE_INIT,
		handler: hdlr,
	}
}

func (pp *PacketParser) Parse(conn *connection) error {
	// @TODO lock
	c := conn.conn
	inited := true
	for c.InboundBuffered() > 0 {
		if inited {
			headBuff, e := c.Next(ppFrameHeaderSize)
			if e != nil {
				fmt.Println("decode ", e)
				return errors.New("invalid header length")
			}
			copy(pp.headerBytes[:], headBuff)
			inited = false
		}
		headRd := bytes.NewReader(pp.headerBytes[:])
		binary.Read(headRd, binary.LittleEndian, &pp.header)
		fmt.Printf("magic %04x\n", pp.header.StartTag)
		fmt.Printf("length %d\n", pp.header.Length)
		if pp.header.StartTag == ppFrameDelimiter {
			data, e := c.Next(int(pp.header.Length + 2))
			if e != nil {
				fmt.Println("rd data ", e.Error())
				return errors.New("invalid data length")
			}
			fmt.Println(string(data))
			pp.endTag = binary.LittleEndian.Uint16(data[pp.header.Length : pp.header.Length+2])
			inited = true
			// @TODO
			// routuine
			packetData := make([]byte, pp.header.Length)
			copy(packetData, data[:pp.header.Length])
			go pp.handler.HandlerPacket(conn, pp.header, packetData)
		} else {
			data, e := c.Next(1)
			if e != nil {
				return errors.New("invalid data")
			}
			copy(pp.headerBytes[0:22], pp.headerBytes[1:23])
			copy(pp.headerBytes[22:23], data)
			headRd.Reset(pp.headerBytes[:])
		}
	}

	return nil
}
