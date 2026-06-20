package encryption

import (
	"testing"
)

func TestAESEncryptionService(t *testing.T) {
	masterKey := "my-secret-master-key-12345"
	service, err := NewAESEncryptionService(masterKey)
	if err != nil {
		t.Fatalf("failed to create encryption service: %v", err)
	}

	plainText := "database-password-value-123!"
	cipherText, err := service.Encrypt(plainText)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if cipherText == plainText {
		t.Errorf("ciphertext is identical to plaintext")
	}

	decryptedText, err := service.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if decryptedText != plainText {
		t.Errorf("expected decrypted text %q, got %q", plainText, decryptedText)
	}
}

func TestAESEncryptionService_InvalidMasterKey(t *testing.T) {
	_, err := NewAESEncryptionService("")
	if err == nil {
		t.Error("expected error when creating service with empty master key, got nil")
	}
}

func TestAESEncryptionService_DecryptInvalidCiphertext(t *testing.T) {
	service, _ := NewAESEncryptionService("some-key")
	_, err := service.Decrypt("invalid-base64-string!@#$")
	if err == nil {
		t.Error("expected error decrypting invalid base64, got nil")
	}
}
