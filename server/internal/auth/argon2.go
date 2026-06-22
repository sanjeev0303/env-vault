package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func DefaultArgon2Config() Argon2Config {
	return Argon2Config{
		Memory:      64 * 1024, // 64MB
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func HashPassword(password string, cfg Argon2Config) (string, error) {
	salt := make([]byte, cfg.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, cfg.Iterations, cfg.Memory, cfg.Parallelism, cfg.KeyLength)

	// Format: $argon2id$v=19$m=MEMORY,t=ITERATIONS,p=PARALLELISM$SALT$HASH
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%x$%x",
		cfg.Memory, cfg.Iterations, cfg.Parallelism, salt, hash), nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	var memory uint32
	var iterations uint32
	var parallelism uint8
	var saltHex, hashHex string

	_, err := fmt.Sscanf(encodedHash, "$argon2id$v=19$m=%d,t=%d,p=%d$%s",
		&memory, &iterations, &parallelism, &saltHex)
	if err != nil {
		return false, fmt.Errorf("failed to parse hash: %w", err)
	}

	// Split salt$hash
	parts := splitLast(saltHex, "$")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid hash format")
	}
	saltHex = parts[0]
	hashHex = parts[1]

	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	expectedHash, err := hex.DecodeString(hashHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	keyLength := uint32(len(expectedHash))
	computedHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	return subtle.ConstantTimeCompare(computedHash, expectedHash) == 1, nil
}

// GenerateRefreshToken creates a cryptographically secure random token.
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// HashToken creates a SHA-256 hash of a token for storage.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func splitLast(s, sep string) []string {
	for i := len(s) - 1; i >= 0; i-- {
		if string(s[i]) == sep {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}
