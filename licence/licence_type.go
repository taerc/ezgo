package licence

import "github.com/taerc/ezgo/licence/lic"

type TimeInfo struct {
	MagicValue     string
	BuildY         int // = -1;
	BuildM         int // = -1;
	BuildD         int // = -1;
	Expired        int // = -1;
	LastY          int //= -1;
	LastM          int //= -1;
	LastD          int //= -1;
	LastH          int //= 0;
	MagicSignature string
}

type LocalInfo struct {
	DeviceSn string
	Uuid     string //hashid(sn+uuid编码sha256编码2）
}

type CentreInfo struct {
	Url        string
	SignLocal  string
	SignCentre string
}

type LicenceProto struct {
	Version        int          // = 0; //1u
	MagicValue     string       //随机字符串
	MagicSignature string       //随机字符串编码sha256编码1
	AuthType       lic.AuthType //=TimeAuth;//LocalAuth本地文件方式鉴权
	DeviceDesc     string       //安卓"android" 盒子"emmc"
	TimeInfo       TimeInfo     //空
	LocalInfo      LocalInfo    //sn和uuid //sn和 hashid(sn+uuid编码sha256编码2）
	CentreInfo     CentreInfo   //空
}

type LicenceProtoType struct {
	Version        int          // = 0; //1u
	MagicValue     string       //随机字符串
	MagicSignature string       //随机字符串编码sha256编码1
	AuthType       lic.AuthType //=TimeAuth;//LocalAuth本地文件方式鉴权
	DeviceDesc     string       //安卓"android" 盒子"emmc"
	//LocalInfo      LocalInfo  //sn和uuid
	DeviceSn string
	Uuid     string
}
