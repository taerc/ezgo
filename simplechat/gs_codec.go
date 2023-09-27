package simplechat

type GSFrame struct {
	StartTag    string
	SendSeq     int  // inc
	RecvSeq     int  // inc
	SessionFlag byte // 0x00 request 0x01 response
	Length      int
	Data        string
	EndTag      string
}

func (f *GSFrame) Encode() {

}

func (f *GSFrame) Decode() {

}
