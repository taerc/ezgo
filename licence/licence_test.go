package licence

import (
	"testing"

	"github.com/go-playground/assert"
	_ "github.com/taerc/ezgo"
)

func TestParseLicenceFile(t *testing.T) {
	//filePath := "/Users/rotaercw/Desktop/0201010123010013.proto"
	//ParseLicenceFile(filePath)
}

func TestGenerateLicence(t *testing.T) {
	dstPath := "/Users/rotaercw/wkspace/logs/exp.proto"

	// generate
	sn := "1234567890"
	uuid := "qwertyuiop"
	desc := "emmc"
	GenerateLicence(dstPath, sn, uuid, desc)

	// true
	r := LicenceCheck(dstPath, sn, uuid)
	assert.Equal(t, r, true)

	// false
	sn = "1234567891"
	r = LicenceCheck(dstPath, sn, uuid)
	assert.Equal(t, r, false)

}
