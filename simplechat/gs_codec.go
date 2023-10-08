package simplechat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"

	"github.com/panjf2000/gnet/v2/pkg/buffer/ring"
)

const (
	gsFrameDelimiter  uint16 = 0xEB90
	gsRequestFrame    byte   = 0x00
	gsResponseFrame   byte   = 0x01
	gsFrameHeaderSize        = 23
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
	header     GSFrameHeader
	data       string
	ringBuffer *ring.Buffer
	peekBuff   []byte
	state      int
	mutex      *sync.Mutex
}

func NewGSFrameDecoder() *GSFrameDecoder {
	return &GSFrameDecoder{
		ringBuffer: ring.New(ring.DefaultBufferSize),
		peekBuff:   make([]byte, 1024),
		mutex:      &sync.Mutex{},
		state:      GSFRAME_DECODE_STATE_INIT,
	}
}

func (gs *GSFrameDecoder) Write(p []byte) (int, error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.ringBuffer.Write(p)
}

func (gs *GSFrameDecoder) Decode() {

	n, _ := gs.load()
	for n > gsFrameHeaderSize {
		gs.scan()
		n, _ = gs.load()
	}
}

func (gs *GSFrameDecoder) scan() {

	edgeLimit := len(gs.peekBuff) - gsFrameHeaderSize
	buffSize := len(gs.peekBuff)
	index := -1
	for i := 0; i < edgeLimit; i++ {
		if gs.state == GSFRAME_DECODE_STATE_INIT && gs.peekBuff[i] == 0x90 {
			gs.state = GSFRAME_DECODE_STATE_90
			index = i
			continue
		}
		if gs.state == GSFRAME_DECODE_STATE_90 && gs.peekBuff[i] == 0xeb && i == index+1 {
			gs.state = GSFRAME_DECODE_STATE_EB
			// try to decode header
			// size checking
			e := gs.fillHeader(gs.peekBuff[i+1 : gsFrameHeaderSize+i-1])
			fmt.Println(e)
			fmt.Println(" data size ", buffSize-gsFrameHeaderSize-i+1)
			if gs.header.Length > uint32(buffSize-gsFrameHeaderSize-i+1) {
				// data merge
				gs.state = GSFRAME_DECODE_STATE_FRAME_BROKEN
			} else {
				if gs.peekBuff[i+21+int(gs.header.Length)] == 0x90 && gs.peekBuff[i+22+int(gs.header.Length)] == 0xeb {
					gs.state = GSFRAME_DECODE_STATE_END
					//TODO
					fmt.Println(gs.peekBuff[i+21 : i+21+int(gs.header.Length)])
					gs.state = GSFRAME_DECODE_STATE_INIT
					i += 22 + int(gs.header.Length)
				} else {

				}

			}
		} else {
			gs.state = GSFRAME_DECODE_STATE_INIT
			index = -1
		}

	}

}

func (gs *GSFrameDecoder) load() (int, error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	// gs.peekBuff = gs.ringBuffer.Bytes()
	n, e := gs.ringBuffer.Read(gs.peekBuff)
	fmt.Printf("load %02x %02x \n", gs.peekBuff[0], gs.peekBuff[1])
	return n, e
}

func (gs *GSFrameDecoder) fillHeader(data []byte) error {
	// binary.Read(r, binary.LittleEndian, &gs.header)
	gs.header.SendSeq = binary.LittleEndian.Uint64(data[:8])
	gs.header.RecvSeq = binary.LittleEndian.Uint64(data[8:16])
	gs.header.Type = data[16]
	gs.header.Length = binary.LittleEndian.Uint32(data[17:21])

	fmt.Println(gs.header.SendSeq)
	fmt.Println(gs.header.RecvSeq)
	fmt.Println(gs.header.Type)
	fmt.Println(gs.header.Length)
	return nil
}
