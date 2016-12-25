package util

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
)

// GetDir get pwd
func GetDir() string {
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		return ""
	}
	return filepath.Dir(path)
}

// MD5 checksum for str
func MD5(str string) string {
	hexStr := md5.Sum([]byte(str))
	return hex.EncodeToString(hexStr[:])
}
