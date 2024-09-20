package transfer

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}
