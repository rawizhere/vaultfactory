package crypto

import (
	"testing"
)

func TestCryptoService_EncryptDecrypt(t *testing.T) {
	crypto := NewCryptoService()

	testData := []byte("test data for encryption")
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	encrypted, err := crypto.Encrypt(testData, key)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	if len(encrypted) == 0 {
		t.Fatal("Encrypted data is empty")
	}

	decrypted, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if string(decrypted) != string(testData) {
		t.Fatalf("Decrypted data doesn't match original. Expected: %s, Got: %s", string(testData), string(decrypted))
	}
}

func TestCryptoService_EncryptDecrypt_DifferentKeys(t *testing.T) {
	crypto := NewCryptoService()

	testData := []byte("test data")
	key1, _ := crypto.GenerateKey()
	key2, _ := crypto.GenerateKey()

	encrypted, err := crypto.Encrypt(testData, key1)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	_, err = crypto.Decrypt(encrypted, key2)
	if err == nil {
		t.Fatal("Expected error when decrypting with wrong key")
	}
}

func TestCryptoService_GenerateKey(t *testing.T) {
	crypto := NewCryptoService()

	key1, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	key2, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	if len(key1) != 32 {
		t.Fatalf("Expected key length 32, got %d", len(key1))
	}

	if string(key1) == string(key2) {
		t.Fatal("Generated keys should be different")
	}
}

func TestCryptoService_HashPassword(t *testing.T) {
	crypto := NewCryptoService()

	password := "testpassword123"
	hash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if len(hash) == 0 {
		t.Fatal("Hash is empty")
	}

	if hash == password {
		t.Fatal("Hash should not equal original password")
	}
}

func TestCryptoService_VerifyPassword(t *testing.T) {
	crypto := NewCryptoService()

	password := "testpassword123"
	hash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	valid := crypto.VerifyPassword(password, hash)
	if !valid {
		t.Fatal("Password verification failed for correct password")
	}

	invalid := crypto.VerifyPassword("wrongpassword", hash)
	if invalid {
		t.Fatal("Password verification should fail for wrong password")
	}
}

func TestCryptoService_VerifyPassword_EmptyPassword(t *testing.T) {
	crypto := NewCryptoService()

	hash, err := crypto.HashPassword("testpassword")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	valid := crypto.VerifyPassword("", hash)
	if valid {
		t.Fatal("Empty password should not be valid")
	}
}

func TestCryptoService_Encrypt_EmptyData(t *testing.T) {
	crypto := NewCryptoService()

	key, _ := crypto.GenerateKey()
	encrypted, err := crypto.Encrypt([]byte{}, key)
	if err != nil {
		t.Fatalf("Failed to encrypt empty data: %v", err)
	}

	decrypted, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Failed to decrypt empty data: %v", err)
	}

	if len(decrypted) != 0 {
		t.Fatal("Decrypted empty data should be empty")
	}
}

func TestCryptoService_Decrypt_InvalidData(t *testing.T) {
	crypto := NewCryptoService()

	key, _ := crypto.GenerateKey()
	_, err := crypto.Decrypt([]byte("invalid data"), key)
	if err == nil {
		t.Fatal("Expected error when decrypting invalid data")
	}
}

func TestCryptoService_Decrypt_TooShortData(t *testing.T) {
	crypto := NewCryptoService()

	key, _ := crypto.GenerateKey()
	_, err := crypto.Decrypt([]byte("short"), key)
	if err == nil {
		t.Fatal("Expected error when decrypting too short data")
	}
}

