// 加密处理
package tool

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/google/uuid"
)

// MD5 使用MD5对数据签名 (长度32)
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha256 使用Sha256对数据签名 (长度64)
func Sha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// Base64Encode
func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// Base64Decode
func Base64Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

// UUID4 生成唯一ID
func UUID() string {
	return uuid.New().String()
}

// UUID1 生成唯一ID
func UUID1() string {
	return uuid.Must(uuid.NewUUID()).String()
}
