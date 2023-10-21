package licence

import (
	"fmt"
	"math/rand"
	"os"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/taerc/ezgo/licence/proto"
	ezgo "github.com/taerc/ezgo/pkg"
)

// GenerateLicence
func GenerateLicence(dstFile, sn, uuid, deviceDesc string) error {
	builder := flatbuffers.NewBuilder(0)
	// time
	proto.TimeInfoStart(builder)
	timeinfo := proto.TimeInfoEnd(builder)
	// local
	protoSn := builder.CreateString(sn)
	protoUUID := builder.CreateString(snAndUUIDMerge(sn, uuid))
	proto.LocalInfoStart(builder)
	// sn
	proto.LocalInfoAddSn(builder, protoSn)
	// uuid
	proto.LocalInfoAddUuid(builder, protoUUID)
	localinfo := proto.LocalInfoEnd(builder)

	// centre
	proto.CentreInfoStart(builder)
	centreinfo := proto.CentreInfoEnd(builder)
	//LicenceProto
	magicValue := randomString()
	protoMagicSign := builder.CreateString(randomStringSha256(magicValue))
	protoMagicValue := builder.CreateString(magicValue)
	protoDesc := builder.CreateString(deviceDesc)

	proto.LicenceProtoStart(builder)
	//version
	proto.LicenceProtoAddVersion(builder, 1)
	// magic
	proto.LicenceProtoAddMagicValue(builder, protoMagicValue)
	// magic sign
	proto.LicenceProtoAddMagicSignature(builder, protoMagicSign)
	proto.LicenceProtoAddAuthType(builder, proto.AuthTypeLocalAuth)
	// desc
	proto.LicenceProtoAddDeviceDesc(builder, protoDesc)
	// timeinfo
	proto.LicenceProtoAddTimeInfo(builder, timeinfo)
	// localinfo
	proto.LicenceProtoAddLocalInfo(builder, localinfo)
	// centre info
	proto.LicenceProtoAddCentreInfo(builder, centreinfo)
	licenceinfo := proto.LicenceProtoEnd(builder)

	builder.Finish(licenceinfo)
	buf := builder.FinishedBytes()
	// return buf
	//异或编码加密
	bufferEncrypt(buf)
	return mixupFileWrite(dstFile, buf)
}

// DecodeLicence
func DecodeLicence(filepathname string) (*proto.LicenceProto, error) {
	data, err := mixupFileRead(filepathname)
	if err != nil {
		return nil, err
	}
	bufferEncrypt(data)
	return decodeLicence(data), nil
}

func LicenceCheck(licPath string, sn string, uuid string) bool {
	data, err := mixupFileRead(licPath)
	if err != nil {
		return false
	}
	bufferEncrypt(data)
	lic := decodeLicence(data)
	fmt.Println(snAndUUIDMerge(sn, uuid))
	fmt.Println(string(lic.LocalInfo(nil).Uuid()))
	if string(lic.LocalInfo(nil).Uuid()) == snAndUUIDMerge(sn, uuid) {
		return true
	}

	return false
}

// Decode 反序列化
func decodeLicence(buf []byte) *proto.LicenceProto {

	gral := proto.GetRootAsLicenceProto(buf, 0)
	local := new(proto.LocalInfo)
	gral.LocalInfo(local)
	fmt.Println(string(local.Sn()))
	fmt.Println(string(local.Uuid()))
	return gral
}

var _mask = [8]byte{0xe6, 0x78, 0x4e, 0xfb, 0xe6, 0x78, 0x4e, 0xfb}

func bufferEncrypt(msg []byte) {

	if len(msg) < 8 {
		return
	}
	n := len(msg) - 7
	for i := 0; i < n; i += 8 {
		msg[i+0] ^= _mask[0]
		msg[i+1] ^= _mask[1]
		msg[i+2] ^= _mask[2]
		msg[i+3] ^= _mask[3]
		msg[i+4] ^= _mask[4]
		msg[i+5] ^= _mask[5]
		msg[i+6] ^= _mask[6]
		msg[i+7] ^= _mask[7]
	}

}

// 可变长度的随机字符串
func randomString() string {
	//len(b)==90
	b := []byte{
		0x22, 0x60, 0x21, 0x40, 0x23, 0x24, 0x25, 0x5e, 0x26, 0x2a, 0x28, 0x29,
		0x2d, 0x3d, 0x2c, 0x2e, 0x3e, 0x3c, 0x3f, 0x3b, 0x3a, 0x27, 0x5b, 0x5d,
		0x7b, 0x7d, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a,
		0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56,
		0x57, 0x58, 0x59, 0x5a, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68,
		0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74,
		0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35,
		0x36, 0x37, 0x38, 0x39, 0x2b, 0x2f}

	randomInt := 32 + rand.Intn(48-32)
	bytes := make([]byte, randomInt)
	for i := 0; i < randomInt; i++ {
		bytes[i] = b[rand.Intn(len(b))]
	}
	return string(bytes)
}

// 随机字符串编码sha256编码1
func randomStringSha256(id string) string {
	//len(b1): 128
	b1 := []byte{
		0x6b, 0x74, 0x50, 0x32, 0x57, 0x73, 0x47, 0x52, 0x42, 0x57, 0x4f, 0x33,
		0x2b, 0x2b, 0x4d, 0x42, 0x33, 0x39, 0x6a, 0x54, 0x45, 0x58, 0x6b, 0x6f,
		0x45, 0x2f, 0x35, 0x4b, 0x6c, 0x74, 0x68, 0x51, 0x72, 0x67, 0x4d, 0x4c,
		0x4f, 0x77, 0x6d, 0x4b, 0x68, 0x76, 0x74, 0x38, 0x31, 0x56, 0x4c, 0x73,
		0x69, 0x77, 0x6b, 0x33, 0x49, 0x69, 0x2f, 0x39, 0x41, 0x64, 0x39, 0x33,
		0x63, 0x74, 0x46, 0x4c, 0x38, 0x69, 0x2b, 0x74, 0x51, 0x66, 0x2b, 0x58,
		0x51, 0x73, 0x7a, 0x47, 0x65, 0x36, 0x58, 0x6d, 0x2f, 0x79, 0x72, 0x2f,
		0x63, 0x62, 0x31, 0x75, 0x6b, 0x4e, 0x46, 0x43, 0x75, 0x46, 0x63, 0x43,
		0x34, 0x39, 0x70, 0x37, 0x2f, 0x50, 0x69, 0x37, 0x76, 0x77, 0x65, 0x72,
		0x39, 0x2b, 0x56, 0x31, 0x7a, 0x70, 0x2f, 0x59, 0x59, 0x58, 0x53, 0x4f,
		0x75, 0x6b, 0x68, 0x53, 0x79, 0x56, 0x31, 0x58}
	//id==随机字符串==string(bytes)
	//id := string(bytes)
	//id+（截取前64位）+id+(剩余的)+id
	str1 := []byte(string(b1))[:64]
	str2 := []byte(string(b1))[64:]
	s1 := ezgo.SHA256(id + string(str1) + id + string(str2) + id)
	return string(s1)
}

// sn+uuid编码sha256编码2
func snAndUUIDMerge(sn, uuid string) string {
	//len(b2): 128
	b2 := []byte{
		0x59, 0x6c, 0x76, 0x41, 0x42, 0x52, 0x2f, 0x46, 0x54, 0x4a, 0x48, 0x6a,
		0x4e, 0x72, 0x37, 0x71, 0x4c, 0x4f, 0x46, 0x64, 0x63, 0x5a, 0x33, 0x31,
		0x2f, 0x72, 0x38, 0x52, 0x55, 0x6a, 0x4d, 0x45, 0x49, 0x46, 0x66, 0x59,
		0x4f, 0x4c, 0x73, 0x46, 0x79, 0x43, 0x4c, 0x56, 0x33, 0x2f, 0x46, 0x77,
		0x50, 0x37, 0x79, 0x48, 0x68, 0x52, 0x68, 0x4b, 0x61, 0x47, 0x4a, 0x67,
		0x34, 0x34, 0x2f, 0x73, 0x2b, 0x69, 0x6a, 0x4d, 0x57, 0x73, 0x50, 0x31,
		0x63, 0x62, 0x6c, 0x65, 0x64, 0x44, 0x6a, 0x6b, 0x49, 0x4d, 0x38, 0x38,
		0x76, 0x49, 0x4b, 0x38, 0x4d, 0x55, 0x67, 0x38, 0x76, 0x2f, 0x2f, 0x67,
		0x65, 0x35, 0x35, 0x31, 0x30, 0x77, 0x61, 0x50, 0x78, 0x62, 0x44, 0x4c,
		0x6e, 0x62, 0x70, 0x61, 0x34, 0x32, 0x39, 0x70, 0x31, 0x59, 0x6e, 0x58,
		0x4d, 0x50, 0x4a, 0x77, 0x51, 0x73, 0x6b, 0x45}

	idstr := "AiRiA" + sn + "SHIP" + uuid

	str1 := []byte(string(b2))[:32]
	str2 := []byte(string(b2))[32:]
	return ezgo.SHA256(idstr + string(str1) + idstr + string(str2) + idstr)
}

// mixupFileWrite
// [xxx-pading][data][datalen-32:datalen-16][datalen-16:]
func mixupFileWrite(filename string, data []byte) error {

	var bufLen int = 1024
	if len(data) > 960 {
		bufLen = (len(data) + 127) / 64 * 64
	}
	randData := ezgo.BytesRandom(bufLen)
	idx := int(ezgo.HashBytesCRC32(randData[bufLen-16:]) & 0x1f)
	copy(randData[idx:], data)
	bL := ezgo.UInt32Bytes(uint32(len(data)))

	for i := 0; i < 4; i += 1 {
		randData[bufLen-32+i] = bL[i]
	}
	return os.WriteFile(filename, randData, os.ModePerm)
}

// mixupFileRead
func mixupFileRead(filePath string) ([]byte, error) {

	if data, e := os.ReadFile(filePath); e == nil {
		idx := int(ezgo.HashBytesCRC32(data[len(data)-16:]) & 0x1f)
		n := int(ezgo.BytesUInt32(data[len(data)-32:]))
		return data[idx : idx+n], nil
	} else {
		return nil, e
	}
}
