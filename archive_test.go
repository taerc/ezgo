package ezgo

import "testing"

func TestZipCompress(t *testing.T) {

	ZipArchive([]string{"/Users/rotaercw/Downloads/hs.jpeg", "/Users/rotaercw/Downloads/test.cpp", "/Users/rotaercw/Downloads/test-image.jpeg"}, "test.zip")
}
