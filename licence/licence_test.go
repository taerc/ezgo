package licence

import "testing"

func TestParseLicenceFile(t *testing.T) {
	//filePath := "/Users/rotaercw/Desktop/0201010123010013.proto"
	//ParseLicenceFile(filePath)
}

func TestGenerateLicence(t *testing.T) {
	dstPath := "/Users/rotaercw/wkspace/logs/exp.proto"
	//GenerateLicence(dstPath, "1234567890", "uuuuuuuuuuuuu", "emmc")
	DecodeLicence(dstPath)
}
