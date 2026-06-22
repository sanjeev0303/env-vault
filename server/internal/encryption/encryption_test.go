package encryption

import (
	"testing"
)

func TestEncryptionService(t *testing.T) {
	svc, err := NewAESEncryptionService("super-secret-key")
	if err != nil {
		t.Fatalf("failed to init: %v", err)
	}

	plainText := "my-secret-password"
	env, err := svc.Encrypt(plainText)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if env.EncryptedValue == plainText {
		t.Fatal("encrypted value is plain text")
	}
	if env.EncryptedDEK == "" || env.IV == "" || env.AuthTag == "" {
		t.Fatal("envelope fields missing")
	}

	decrypted, err := svc.Decrypt(env)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if decrypted != plainText {
		t.Fatalf("expected %q, got %q", plainText, decrypted)
	}

	// Test invalid decryption
	env.EncryptedValue = "invalid-base64-string!@#$"
	_, err = svc.Decrypt(env)
	if err == nil {
		t.Fatal("expected error with invalid base64")
	}
}
