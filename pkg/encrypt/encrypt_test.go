package encrypt

import (
    "os"
    "strings"
    "testing"
)

func TestByte16ToBytes(t *testing.T) {
    input := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
    result := byte16ToBytes(input)
    
    if len(result) != 16 {
        t.Errorf("Expected length 16, got %d", len(result))
    }
    
    for i, v := range input {
        if result[i] != v {
            t.Errorf("At index %d, expected %d, got %d", i, v, result[i])
        }
    }
}

func TestEncPasswordBcrypt(t *testing.T) {
    password := "test123"
    hash := EncPassword(password)
    
    if hash == "" {
        t.Error("Expected non-empty hash")
    }
    
    if len(hash) < 50 {
        t.Errorf("Bcrypt hash should be at least 50 chars, got %d", len(hash))
    }
    
    if strings.HasPrefix(hash, "$2") {
        t.Log("Successfully using bcrypt hash")
    }
}

func TestVerifyPassword(t *testing.T) {
    tests := []struct {
        name     string
        password string
    }{
        {"simple", "test123"},
        {"complex", "P@ssw0rd!123"},
        {"chinese", "密码123"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            hash := EncPassword(tt.password)
            if !VerifyPassword(hash, tt.password) {
                t.Errorf("Password verification failed for %s", tt.password)
            }
            
            if VerifyPassword(hash, "wrongpassword") {
                t.Error("Verification should fail for wrong password")
            }
        })
    }
}

func TestVerifyPasswordBackwardCompatibility(t *testing.T) {
    oldSeed := passwordEncryptSeed
    password := "test123"
    
    oldMD5Hash := Md5Sum([]byte(strings.TrimSpace(password + oldSeed)))
    
    if !VerifyPassword(oldMD5Hash, password) {
        t.Error("Backward compatibility with MD5 failed")
    }
    
    if VerifyPassword(oldMD5Hash, "wrongpassword") {
        t.Error("MD5 verification should fail for wrong password")
    }
}

func TestGetEnvOrDefault(t *testing.T) {
    key := "TEST_ENV_KEY"
    defaultValue := "default"
    customValue := "custom"
    
    result := getEnvOrDefault(key, defaultValue)
    if result != defaultValue {
        t.Errorf("Expected %s, got %s", defaultValue, result)
    }
    
    os.Setenv(key, customValue)
    defer os.Unsetenv(key)
    
    result = getEnvOrDefault(key, defaultValue)
    if result != customValue {
        t.Errorf("Expected %s, got %s", customValue, result)
    }
}

func TestEncPhoneDecPhone(t *testing.T) {
    phones := []string{
        "13800138000",
        "18612345678",
        "15900000000",
    }
    
    for _, phone := range phones {
        encrypted, err := EncPhone(phone)
        if err != nil {
            t.Errorf("EncPhone failed for %s: %v", phone, err)
            continue
        }
        
        decrypted, err := DecPhone(encrypted)
        if err != nil {
            t.Errorf("DecPhone failed for %s: %v", encrypted, err)
            continue
        }
        
        if decrypted != phone {
            t.Errorf("Expected %s, got %s", phone, decrypted)
        }
    }
}

func TestMd5Sum(t *testing.T) {
    data := []byte("test")
    hash := Md5Sum(data)
    
    if len(hash) != 32 {
        t.Errorf("MD5 hash should be 32 chars, got %d", len(hash))
    }
    
    expectedHash := "098f6bcd4621d373cade4e832627b4f6"
    if hash != expectedHash {
        t.Errorf("Expected %s, got %s", expectedHash, hash)
    }
}

func BenchmarkEncPassword(b *testing.B) {
    password := "test123"
    for i := 0; i < b.N; i++ {
        EncPassword(password)
    }
}

func BenchmarkVerifyPassword(b *testing.B) {
    password := "test123"
    hash := EncPassword(password)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        VerifyPassword(hash, password)
    }
}
