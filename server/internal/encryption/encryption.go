package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"env-vault/server/internal/domain"
)

type AESEncryptionService struct {
	masterKey []byte
}

func NewAESEncryptionService(masterKey string) (*AESEncryptionService, error) {
	if masterKey == "" {
		return nil, errors.New("master key cannot be empty")
	}
	hash := sha256.Sum256([]byte(masterKey))
	return &AESEncryptionService{masterKey: hash[:]}, nil
}

func (s *AESEncryptionService) Encrypt(plainText string) (*domain.Envelope, error) {
	// 1. Generate DEK (32 bytes for AES-256)
	dek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dek); err != nil {
		return nil, err
	}

	// 2. Encrypt DEK with Master Key
	masterBlock, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return nil, err
	}
	masterGCM, err := cipher.NewGCM(masterBlock)
	if err != nil {
		return nil, err
	}
	masterNonce := make([]byte, masterGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, masterNonce); err != nil {
		return nil, err
	}
	encryptedDEK := masterGCM.Seal(masterNonce, masterNonce, dek, nil)

	// 3. Encrypt plainText with DEK
	dekBlock, err := aes.NewCipher(dek)
	if err != nil {
		return nil, err
	}
	dekGCM, err := cipher.NewGCM(dekBlock)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, dekGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	cipherTextAndTag := dekGCM.Seal(nil, nonce, []byte(plainText), nil)

	// In GCM, Seal appends the auth tag to the end of the ciphertext.
	// AES-GCM tag size is 16 bytes.
	tagSize := dekGCM.Overhead()
	if len(cipherTextAndTag) < tagSize {
		return nil, errors.New("cipherText too short")
	}

	cipherTextOnly := cipherTextAndTag[:len(cipherTextAndTag)-tagSize]
	authTag := cipherTextAndTag[len(cipherTextAndTag)-tagSize:]

	return &domain.Envelope{
		EncryptedValue: base64.StdEncoding.EncodeToString(cipherTextOnly),
		EncryptedDEK:   base64.StdEncoding.EncodeToString(encryptedDEK),
		IV:             base64.StdEncoding.EncodeToString(nonce),
		AuthTag:        base64.StdEncoding.EncodeToString(authTag),
	}, nil
}

func (s *AESEncryptionService) Decrypt(env *domain.Envelope) (string, error) {
	// 1. Decrypt DEK with Master Key
	encryptedDEKBytes, err := base64.StdEncoding.DecodeString(env.EncryptedDEK)
	if err != nil {
		return "", err
	}
	masterBlock, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return "", err
	}
	masterGCM, err := cipher.NewGCM(masterBlock)
	if err != nil {
		return "", err
	}
	nonceSize := masterGCM.NonceSize()
	if len(encryptedDEKBytes) < nonceSize {
		return "", errors.New("encrypted DEK too short")
	}
	masterNonce, masterCipher := encryptedDEKBytes[:nonceSize], encryptedDEKBytes[nonceSize:]
	dek, err := masterGCM.Open(nil, masterNonce, masterCipher, nil)
	if err != nil {
		return "", err
	}

	// 2. Decrypt plainText with DEK
	dekBlock, err := aes.NewCipher(dek)
	if err != nil {
		return "", err
	}
	dekGCM, err := cipher.NewGCM(dekBlock)
	if err != nil {
		return "", err
	}

	nonce, err := base64.StdEncoding.DecodeString(env.IV)
	if err != nil {
		return "", err
	}
	cipherTextOnly, err := base64.StdEncoding.DecodeString(env.EncryptedValue)
	if err != nil {
		return "", err
	}
	authTag, err := base64.StdEncoding.DecodeString(env.AuthTag)
	if err != nil {
		return "", err
	}

	// Reconstruct cipherTextAndTag
	cipherTextAndTag := append(cipherTextOnly, authTag...)

	plainText, err := dekGCM.Open(nil, nonce, cipherTextAndTag, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
