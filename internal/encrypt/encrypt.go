package encrypt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

//
// ------------------ Argon2id password hashing & verification ------------------
//

// EncryptPassword hashes the password using Argon2id
func EncryptPassword(password string) (string, error) {
	var (
		memory      uint32 = 64 * 1024 // 64MB
		iterations  uint32 = 3
		parallelism uint8  = 2
		saltLength  uint32 = 16
		keyLength   uint32 = 32
	)

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		memory, iterations, parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

// VerifyPassword checks whether the password is correct
func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format: expected 6 parts, got %d", len(parts))
	}

	if parts[1] != "argon2id" || parts[2] != "v=19" {
		return false, fmt.Errorf("unsupported hash algorithm or version")
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, fmt.Errorf("invalid parameters format: %v", err)
	}

	saltB64 := parts[4]
	hashB64 := parts[5]

	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(hashB64)
	if err != nil {
		return false, err
	}

	newHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))
	return subtleConstantTimeCompare(hash, newHash), nil
}

func subtleConstantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := range a {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
