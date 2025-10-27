package encrypt

import (
    "crypto/md5"
    "encoding/base64"
    "encoding/hex"
    "os"
    "strings"

    "github.com/zeromicro/go-zero/core/codec"
    "golang.org/x/crypto/bcrypt"
)

var (
    passwordEncryptSeed string
    phoneAesKey         string
)

func init() {
    passwordEncryptSeed = getEnvOrDefault("PASSWORD_ENCRYPT_SEED", "(infoflow)@#$")
    phoneAesKey = getEnvOrDefault("PHONE_AES_KEY", "5A2E746B08D846502F37A6E2D85D583B")
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func EncPassword(password string) string {
    hash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(password)), bcrypt.DefaultCost)
    if err != nil {
        return Md5Sum([]byte(strings.TrimSpace(password + passwordEncryptSeed)))
    }
    return string(hash)
}

func VerifyPassword(hashedPassword, password string) bool {
    password = strings.TrimSpace(password)
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err == nil {
        return true
    }
    return hashedPassword == Md5Sum([]byte(password+passwordEncryptSeed))
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
    return in[:]
}
