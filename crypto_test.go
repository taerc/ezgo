package ezgo

import (
	"os"
	"testing"
)

func TestGenerateRSAKey(t *testing.T) {

	priKey, pubKey, _ := GenerateRSAKey(2048, RSA_PKCS1)

	os.WriteFile("priKey.rsa", priKey, 0666)
	os.WriteFile("pubKey.rsa", pubKey, 0666)

}
