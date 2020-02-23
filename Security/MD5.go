package Security

import (
	"crypto/md5"
)

// MD5加密
func MD5Encrypt(origData []byte) []byte {
	h := md5.New()
	h.Write(origData)
	cipherStr := h.Sum(nil)
	return cipherStr
}
