package licence

import (
	"fmt"
	"math/rand"
	"os"

	flatbuffers "github.com/google/flatbuffers/go"
	log "github.com/sirupsen/logrus"
	"github.com/taerc/ezgo/licence/proto"
	ezgo "github.com/taerc/ezgo/pkg"
)

// type TimeInfo struct {
// 	MagicValue     string
// 	BuildY         int // = -1;
// 	BuildM         int // = -1;
// 	BuildD         int // = -1;
// 	Expired        int // = -1;
// 	LastY          int //= -1;
// 	LastM          int //= -1;
// 	LastD          int //= -1;
// 	LastH          int //= 0;
// 	MagicSignature string
// }

// type LocalInfo struct {
// 	DeviceSn string
// 	Uuid     string //hashid(sn+uuid编码sha256编码2）
// }

// type CentreInfo struct {
// 	Url        string
// 	SignLocal  string
// 	SignCentre string
// }

// type LicenceProto struct {
// 	Version        int            // = 0; //1u
// 	MagicValue     string         //随机字符串
// 	MagicSignature string         //随机字符串编码sha256编码1
// 	AuthType       proto.AuthType //=TimeAuth;//LocalAuth本地文件方式鉴权
// 	DeviceDesc     string         //安卓"android" 盒子"emmc"
// 	TimeInfo       TimeInfo       //空
// 	LocalInfo      LocalInfo      //sn和uuid //sn和 hashid(sn+uuid编码sha256编码2）
// 	CentreInfo     CentreInfo     //空
// }

// type LicenceProtoType struct {
// 	Version        int            // = 0; //1u
// 	MagicValue     string         //随机字符串
// 	MagicSignature string         //随机字符串编码sha256编码1
// 	AuthType       proto.AuthType //=TimeAuth;//LocalAuth本地文件方式鉴权
// 	DeviceDesc     string         //安卓"android" 盒子"emmc"
// 	//LocalInfo      LocalInfo  //sn和uuid
// 	DeviceSn string
// 	Uuid     string
// }

// GenerateLicence
func GenerateLicence(dstFile, sn, uuid, DeviceDesc string) error {
	//初始化buffer，大小为0，会自动扩容
	builder := flatbuffers.NewBuilder(0)
	//builder.CreateString()
	lic := new(proto.LicenceProto)
	// licenceProtoType := new(LicenceProtoType)
	licenceProtoType.Version = 1                                                      // = 0; //1u
	licenceProtoType.MagicValue = randomString()                                      //随机字符串
	licenceProtoType.MagicSignature = randomStringSha256(licenceProtoType.MagicValue) //随机字符串编码sha256编码1
	licenceProtoType.AuthType = proto.AuthTypeLocalAuth                               //=TimeAuth;//LocalAuth本地文件方式鉴权
	licenceProtoType.DeviceDesc = DeviceDesc                                          //安卓"android" 盒子"emmc" linux是"dmi"
	licenceProtoType.DeviceSn = sn
	licenceProtoType.Uuid = snAndUUIDMerge(licenceProtoType.DeviceSn, uuid)
	//序列化
	// buf := encodeLicence(licenceProtoType)
	// localInfo := LocalInfo{DeviceSn: licenceProtoType.DeviceSn, Uuid: licenceProtoType.Uuid}

	timeInfo := TimeInfo{}
	centreInfo := CentreInfo{}

	licenceProto := new(LicenceProto)
	licenceProto.Version = licenceProtoType.Version               // = 0; //1u 类型
	licenceProto.MagicValue = licenceProtoType.MagicValue         //随机字符串
	licenceProto.MagicSignature = licenceProtoType.MagicSignature //随机字符串编码sha256编码1
	licenceProto.AuthType = licenceProtoType.AuthType             //=TimeAuth;
	licenceProto.DeviceDesc = licenceProtoType.DeviceDesc         //安卓"android" 盒子"emmc"
	licenceProto.TimeInfo = timeInfo                              //空

	licenceProto.CentreInfo = centreInfo //空

	licenceProto.LocalInfo = localInfo //sn和uuid //sn和 hashid(sn+uuid编码sha256编码2）

	infosn := builder.CreateString(localInfo.DeviceSn) //先处理非标量string,得到偏移量
	infouuiid := builder.CreateString(localInfo.Uuid)
	MagicValue := builder.CreateString(licenceProto.MagicValue)

	MagicSignature := builder.CreateString(licenceProto.MagicSignature)
	DeviceDesc := builder.CreateString(licenceProto.DeviceDesc)

	//LocalInfo
	proto.LocalInfoStart(builder)
	proto.LocalInfoAddSn(builder, infosn)
	proto.LocalInfoAddUuid(builder, infouuiid)
	ninfo := proto.LocalInfoEnd(builder)

	//TimeInfo
	proto.TimeInfoStart(builder)
	timeinfo := proto.TimeInfoEnd(builder)

	//CentreInfo
	proto.CentreInfoStart(builder)
	centreinfo := proto.CentreInfoEnd(builder)

	//LicenceProto
	proto.LicenceProtoStart(builder)

	proto.LicenceProtoAddVersion(builder, 1)
	proto.LicenceProtoAddMagicValue(builder, MagicValue)
	proto.LicenceProtoAddMagicSignature(builder, MagicSignature)
	proto.LicenceProtoAddAuthType(builder, proto.AuthTypeLocalAuth)
	proto.LicenceProtoAddDeviceDesc(builder, DeviceDesc)

	proto.LicenceProtoAddTimeInfo(builder, timeinfo)
	proto.LicenceProtoAddLocalInfo(builder, ninfo)
	proto.LicenceProtoAddCentreInfo(builder, centreinfo)

	licenceinfo := proto.LicenceProtoEnd(builder)

	builder.Finish(licenceinfo)
	buf := builder.FinishedBytes() //返回[]byte
	return buf
	//异或编码加密
	bufferEncrypt(buf)
	return mixupFileWrite(dstFile, buf)
}

// DecodeLicence
func DecodeLicence(filepathname string) (*proto.LicenceProto, error) {
	data, err := mixupFileRead(filepathname)
	if err != nil {
		log.Errorln("ParseLicenceFile ioutil.ReadFile Error:", err)
		return nil, err
	}
	bufferEncrypt(data)
	return decodeLicence(data), nil
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
