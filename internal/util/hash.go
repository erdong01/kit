package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func MD5HashString(data string) string {
	m := md5.New()
	m.Write([]byte(data))
	has := m.Sum(nil)
	return hex.EncodeToString(has)
}

// MD5Hash MD5哈希值
func MD5Hash(b []byte) string {
	h := md5.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA1Hash SHA1哈希值
func SHA1Hash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA1HashString SHA1哈希值
func SHA1HashString(s string) string {
	return SHA1Hash([]byte(s))
}
