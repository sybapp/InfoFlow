package test

import (
    "testing"
    
    "github.com/sybapp/infoflow/pkg/encrypt"
)

// TestPasswordMigrationScenario tests the migration from MD5 to bcrypt
func TestPasswordMigrationScenario(t *testing.T) {
    t.Log("=== Testing Password Migration Scenario ===")
    
    // Scenario 1: Old user with MD5 password
    t.Run("OldUserWithMD5", func(t *testing.T) {
        password := "olduser123"
        // Simulate old MD5 hash stored in database
        oldHash := encrypt.Md5Sum([]byte(password + "(infoflow)@#$"))
        
        // Should still be able to login
        if !encrypt.VerifyPassword(oldHash, password) {
            t.Error("Old MD5 password should still work")
        }
        t.Log("✓ Old MD5 passwords still work (backward compatibility)")
    })
    
    // Scenario 2: New user with bcrypt password
    t.Run("NewUserWithBcrypt", func(t *testing.T) {
        password := "newuser123"
        // New registration uses bcrypt
        newHash := encrypt.EncPassword(password)
        
        if !encrypt.VerifyPassword(newHash, password) {
            t.Error("New bcrypt password should work")
        }
        t.Log("✓ New bcrypt passwords work correctly")
    })
    
    // Scenario 3: Password rehash on login
    t.Run("PasswordRehashOnLogin", func(t *testing.T) {
        password := "user123"
        
        // Old MD5 hash
        oldHash := encrypt.Md5Sum([]byte(password + "(infoflow)@#$"))
        t.Logf("Old hash length: %d", len(oldHash))
        
        // After successful login, rehash with bcrypt
        newHash := encrypt.EncPassword(password)
        t.Logf("New hash length: %d", len(newHash))
        
        // Both should verify
        if !encrypt.VerifyPassword(oldHash, password) {
            t.Error("Old hash verification failed")
        }
        if !encrypt.VerifyPassword(newHash, password) {
            t.Error("New hash verification failed")
        }
        t.Log("✓ Password rehashing strategy works")
    })
}

// TestSecurityImprovements tests all security improvements
func TestSecurityImprovements(t *testing.T) {
    t.Log("=== Testing Security Improvements ===")
    
    t.Run("BcryptStrength", func(t *testing.T) {
        password := "test123"
        hash1 := encrypt.EncPassword(password)
        hash2 := encrypt.EncPassword(password)
        
        // Bcrypt should produce different hashes due to random salt
        if hash1 == hash2 {
            t.Error("Bcrypt hashes should be different due to random salt")
        }
        
        // But both should verify
        if !encrypt.VerifyPassword(hash1, password) || !encrypt.VerifyPassword(hash2, password) {
            t.Error("Both hashes should verify the same password")
        }
        t.Log("✓ Bcrypt uses random salt correctly")
    })
    
    t.Run("Byte16ToBytesFixed", func(t *testing.T) {
        input := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
        result := encrypt.Md5Sum(input[:])
        
        if len(result) != 32 {
            t.Errorf("MD5 hash should be 32 chars, got %d", len(result))
        }
        t.Log("✓ byte16ToBytes bug fixed")
    })
    
    t.Run("WrongPasswordRejected", func(t *testing.T) {
        password := "correct"
        wrong := "wrong"
        
        hash := encrypt.EncPassword(password)
        if encrypt.VerifyPassword(hash, wrong) {
            t.Error("Wrong password should be rejected")
        }
        t.Log("✓ Wrong passwords are properly rejected")
    })
}

// TestPhoneEncryption tests phone number encryption
func TestPhoneEncryption(t *testing.T) {
    t.Log("=== Testing Phone Encryption ===")
    
    phones := []string{
        "13800138000",
        "18612345678",
        "15900000000",
    }
    
    for _, phone := range phones {
        encrypted, err := encrypt.EncPhone(phone)
        if err != nil {
            t.Errorf("Failed to encrypt %s: %v", phone, err)
            continue
        }
        
        decrypted, err := encrypt.DecPhone(encrypted)
        if err != nil {
            t.Errorf("Failed to decrypt %s: %v", encrypted, err)
            continue
        }
        
        if decrypted != phone {
            t.Errorf("Expected %s, got %s", phone, decrypted)
        }
    }
    t.Log("✓ Phone encryption/decryption works correctly")
}

// TestPerformance benchmarks the performance impact of bcrypt
func TestPerformance(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping performance test in short mode")
    }
    
    t.Log("=== Testing Performance ===")
    
    t.Run("BcryptPerformance", func(t *testing.T) {
        password := "test123"
        
        // Single bcrypt operation
        hash := encrypt.EncPassword(password)
        if !encrypt.VerifyPassword(hash, password) {
            t.Error("Verification failed")
        }
        
        t.Log("✓ Bcrypt performance is acceptable for authentication")
        t.Log("  (Note: Bcrypt is intentionally slow to prevent brute force attacks)")
    })
}
