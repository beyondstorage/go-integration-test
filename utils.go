package tests

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func calculateHashFromReader(r io.Reader) string {
	rMD5 := md5.New()
	if _, err := io.Copy(rMD5, r); err != nil {
		return ""
	}
	return hex.EncodeToString(rMD5.Sum(nil))
}
