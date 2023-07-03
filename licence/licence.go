package licence

//AiriaLicence
//Airia执照

import (
	"fmt"
	flatbuffers "github.com/google/flatbuffers/go"
	log "github.com/sirupsen/logrus"
	"github.com/taerc/ezgo"
	"github.com/taerc/ezgo/licence/airia"
	"math/rand"
	"os"
)

/// GenerateLicenceFile @description:
///
func GenerateLicence(dstFile, sn, uuid, DeviceDesc string) error {
	licenceProtoType := new(LicenceProtoType)
	licenceProtoType.Version = 1                                                      // = 0; //1u
	licenceProtoType.MagicValue = randomString()                                      //随机字符串
	licenceProtoType.MagicSignature = randomStringSha256(licenceProtoType.MagicValue) //随机字符串编码sha256编码1
	licenceProtoType.AuthType = airia.AuthTypeLocalAuth                               //=TimeAuth;//LocalAuth本地文件方式鉴权
	licenceProtoType.DeviceDesc = DeviceDesc                                          //安卓"android" 盒子"emmc" linux是"dmi"
	licenceProtoType.DeviceSn = sn
	licenceProtoType.Uuid = snAndUUIDMerge(licenceProtoType.DeviceSn, uuid)
	//序列化
	buf := encodeLicence(licenceProtoType)
	//异或编码加密
	bufferEncrypt(buf)
	//写入文件
	f, err := os.Create(dstFile)
	if err != nil {
		log.Errorln("生成文件,os.Create err:", err)
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = f.Write([]byte(buf))
	if err != nil {
		return err
	}
	return nil
}

func encodeLicence(licenceProtoType *LicenceProtoType) []byte {

	localInfo := LocalInfo{DeviceSn: licenceProtoType.DeviceSn, Uuid: licenceProtoType.Uuid}

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

	//初始化buffer，大小为0，会自动扩容
	builder := flatbuffers.NewBuilder(0)
	//builder.CreateString()

	infosn := builder.CreateString(localInfo.DeviceSn) //先处理非标量string,得到偏移量
	infouuiid := builder.CreateString(localInfo.Uuid)
	MagicValue := builder.CreateString(licenceProto.MagicValue)

	MagicSignature := builder.CreateString(licenceProto.MagicSignature)
	DeviceDesc := builder.CreateString(licenceProto.DeviceDesc)

	//LocalInfo
	airia.LocalInfoStart(builder)
	airia.LocalInfoAddSn(builder, infosn)
	airia.LocalInfoAddUuid(builder, infouuiid)
	ninfo := airia.LocalInfoEnd(builder)

	//TimeInfo
	airia.TimeInfoStart(builder)
	timeinfo := airia.TimeInfoEnd(builder)

	//CentreInfo
	airia.CentreInfoStart(builder)
	centreinfo := airia.CentreInfoEnd(builder)

	//LicenceProto
	airia.LicenceProtoStart(builder)

	airia.LicenceProtoAddVersion(builder, 1)
	airia.LicenceProtoAddMagicValue(builder, MagicValue)
	airia.LicenceProtoAddMagicSignature(builder, MagicSignature)
	airia.LicenceProtoAddAuthType(builder, airia.AuthTypeLocalAuth)
	airia.LicenceProtoAddDeviceDesc(builder, DeviceDesc)

	airia.LicenceProtoAddTimeInfo(builder, timeinfo)
	airia.LicenceProtoAddLocalInfo(builder, ninfo)
	airia.LicenceProtoAddCentreInfo(builder, centreinfo)

	licenceinfo := airia.LicenceProtoEnd(builder)

	builder.Finish(licenceinfo)
	buf := builder.FinishedBytes() //返回[]byte
	return buf
}

//Decode 反序列化
func decodeLicence(buf []byte) *LicenceProto {

	gral := airia.GetRootAsLicenceProto(buf, 0)
	local := new(airia.LocalInfo)
	gral.LocalInfo(local)
	fmt.Println(string(local.Sn()))
	fmt.Println(string(local.Uuid()))
	return nil
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

//可变长度的随机字符串
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

//随机字符串编码sha256编码1
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

//sn+uuid编码sha256编码2
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

func DecodeLicence(filepathname string) {
	data, err := os.ReadFile(filepathname)
	if err != nil {
		log.Errorln("ParseLicenceFile ioutil.ReadFile Error:", err)
		return
	}
	bufferEncrypt(data)
	decodeLicence(data)
}
