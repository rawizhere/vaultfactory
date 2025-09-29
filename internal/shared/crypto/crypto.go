package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

type CryptoService struct {
	argonParams *ArgonParams
}

type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func NewCryptoService() *CryptoService {
	return &CryptoService{
		argonParams: &ArgonParams{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
}

func (c *CryptoService) Encrypt(data []byte, key []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := aead.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (c *CryptoService) Decrypt(encryptedData []byte, key []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

func (c *CryptoService) GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

func (c *CryptoService) HashPassword(password string) (string, error) {
	salt, err := c.generateRandomBytes(c.argonParams.SaltLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, c.argonParams.Iterations, c.argonParams.Memory, c.argonParams.Parallelism, c.argonParams.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, c.argonParams.Memory, c.argonParams.Iterations, c.argonParams.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func (c *CryptoService) VerifyPassword(password, encodedHash string) bool {
	salt, hash, err := c.decodeHash(encodedHash)
	if err != nil {
		return false
	}

	otherHash := argon2.IDKey([]byte(password), salt, c.argonParams.Iterations, c.argonParams.Memory, c.argonParams.Parallelism, c.argonParams.KeyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true
	}
	return false
}

func (c *CryptoService) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *CryptoService) decodeHash(encodedHash string) (salt, hash []byte, err error) {
	vals := make([]string, 6)
	parts := 0
	for _, r := range encodedHash {
		if r == '$' {
			parts++
			if parts > 5 {
				break
			}
			continue
		}
		vals[parts] += string(r)
	}

	if parts != 5 {
		return nil, nil, fmt.Errorf("invalid hash format")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, fmt.Errorf("incompatible argon2 version")
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, err
	}

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, err
	}

	return salt, hash, nil
}
