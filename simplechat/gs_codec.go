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

type GSFrameDecode struct {
	header     GSFrameCodec
	data       string
	ringBuffer *ring.Buffer
	peekBuff   []byte
	mutex      *sync.Mutex
}

func NewGSFrameDecoder() *GSFrameDecode {
	return &GSFrameDecode{
		ringBuffer: ring.New(ring.DefaultBufferSize),
		peekBuff:   make([]byte, 1024),
		mutex:      &sync.Mutex{},
	}
}

func (gs *GSFrameDecode) Write(p []byte) (int, error) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.ringBuffer.Write(p)
}

func (gs *GSFrameDecode) Decode() {

}

func (gs *GSFrameDecode) scan() {

	bufSize := len(gs.peekBuff)

	startTag := GSFRAME_DECODE_STATE_INIT
	index := -1
	for i := 0; i < bufSize; i++ {
		if gs.peekBuff[i] == 0x90 {
			startTag = GSFRAME_DECODE_STATE_90
			index = i
			continue
		}
		if startTag == GSFRAME_DECODE_STATE_90 && gs.peekBuff[i] == 0xeb && i == index+1 {
			startTag = GSFRAME_DECODE_STATE_EB
		} else {
			startTag = GSFRAME_DECODE_STATE_INIT
			index = -1
		}

	}

}

const (
	GSFRAME_DECODE_STATE_INIT = 0 + iota
	GSFRAME_DECODE_STATE_PRE
	GSFRAME_DECODE_STATE_90
	GSFRAME_DECODE_STATE_EB
)

func (gs *GSFrameDecode) load() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.peekBuff = gs.ringBuffer.Bytes()
}
