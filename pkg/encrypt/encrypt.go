package encrypt

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/zeromicro/go-zero/core/codec"
)

const (
	passwordEncryptSeed = "(infoflow)@#$"
	phoneAesKey         = "5A2E746B08D846502F37A6E2D85D583B"
)

func EncPassword(password string) string {
	return Md5Sum([]byte(strings.TrimSpace(password + passwordEncryptSeed)))
}

func EncPhone(mobile string) (string, error) {
	data, err := codec.EcbEncrypt([]byte(phoneAesKey), []byte(mobile))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func DecPhone(mobile string) (string, error) {
	originalData, err := base64.StdEncoding.DecodeString(mobile)
	if err != nil {
		return "", err
	}
	data, err := codec.EcbDecrypt([]byte(phoneAesKey), originalData)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func Md5Sum(data []byte) string {
	return hex.EncodeToString(byte16ToBytes(md5.Sum(data)))
}

func byte16ToBytes(in [16]byte) []byte {
	tmp := make([]byte, 16)
	for _, value := range in {
		tmp = append(tmp, value)
	}
	return tmp[16:]
}
