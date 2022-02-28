package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Sum(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
