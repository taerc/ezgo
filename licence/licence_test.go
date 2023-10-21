package licence

import (
	"testing"

	"github.com/go-playground/assert/v2"
	_ "github.com/taerc/ezgo"
)

func TestParseLicenceFile(t *testing.T) {
	//filePath := "/Users/rotaercw/Desktop/0201010123010013.proto"
	//ParseLicenceFile(filePath)
}

func TestGenerateLicence(t *testing.T) {
	dstPath := "/Users/rotaercw/wkspace/logs/exp.proto"

	sn := "1234567890"
	uuid := "qwertyuiop"
	desc := "emmc"
	GenerateLicence(dstPath, sn, uuid, desc)
	p, e := DecodeLicence(dstPath)
	if e != nil {
		t.Log(e)
	}
	assert.Equal(t, p.LocalInfo.Uuid, uuid)
	// assert.Equal(t, p.LocalInfo.DeviceSn, sn)
}
