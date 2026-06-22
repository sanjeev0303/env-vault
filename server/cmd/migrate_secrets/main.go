package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"

	"env-vault/server/internal/config"
	"env-vault/server/internal/database"
	"env-vault/server/internal/encryption"
)

func main() {
	log.Println("Starting secret migration to Envelope Encryption...")

	cfg := config.LoadConfig()
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer db.Close()

	// Initialize new envelope encryptor
	newEncryptor, err := encryption.NewAESEncryptionService(cfg.MasterKey)
	if err != nil {
		log.Fatalf("Failed to init new encryptor: %v", err)
	}

	// Legacy decryptor (same key hash logic, just old AES-GCM directly on value)
	hash := sha256.Sum256([]byte(cfg.MasterKey))
	legacyKey := hash[:]

	ctx := context.Background()
	
	// We also need to migrate secret_versions. Let's do both.
	
	log.Println("Migrating secrets table...")
	migrateSecretsTable(ctx, db, legacyKey, newEncryptor)

	log.Println("Migrating secret_versions table...")
	migrateSecretVersionsTable(ctx, db, legacyKey, newEncryptor)

	log.Println("Migration complete!")
}

func migrateSecretsTable(ctx context.Context, db *sql.DB, legacyKey []byte, newEncryptor *encryption.AESEncryptionService) {
	rows, err := db.QueryContext(ctx, "SELECT id, value FROM secrets WHERE encrypted_dek = ''")
	if err != nil {
		log.Fatalf("Failed to query secrets: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, oldCipher string
		if err := rows.Scan(&id, &oldCipher); err != nil {
			log.Fatalf("Scan error: %v", err)
		}

		plainText, err := legacyDecrypt(legacyKey, oldCipher)
		if err != nil {
			log.Printf("Failed to decrypt secret %s (might already be migrated or corrupted): %v\n", id, err)
			continue
		}

		env, err := newEncryptor.Encrypt(plainText)
		if err != nil {
			log.Fatalf("Failed to encrypt with envelope for %s: %v", id, err)
		}

		_, err = db.ExecContext(ctx, "UPDATE secrets SET value = $1, encrypted_dek = $2, iv = $3, auth_tag = $4 WHERE id = $5",
			env.EncryptedValue, env.EncryptedDEK, env.IV, env.AuthTag, id)
		if err != nil {
			log.Fatalf("Failed to update secret %s: %v", id, err)
		}
		count++
	}
	log.Printf("Migrated %d secrets.\n", count)
}

func migrateSecretVersionsTable(ctx context.Context, db *sql.DB, legacyKey []byte, newEncryptor *encryption.AESEncryptionService) {
	rows, err := db.QueryContext(ctx, "SELECT id, encrypted_value FROM secret_versions WHERE encrypted_dek = ''")
	if err != nil {
		log.Fatalf("Failed to query secret_versions: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, oldCipher string
		if err := rows.Scan(&id, &oldCipher); err != nil {
			log.Fatalf("Scan error: %v", err)
		}

		plainText, err := legacyDecrypt(legacyKey, oldCipher)
		if err != nil {
			log.Printf("Failed to decrypt secret version %s: %v\n", id, err)
			continue
		}

		env, err := newEncryptor.Encrypt(plainText)
		if err != nil {
			log.Fatalf("Failed to encrypt version with envelope for %s: %v", id, err)
		}

		_, err = db.ExecContext(ctx, "UPDATE secret_versions SET encrypted_value = $1, encrypted_dek = $2, iv = $3, auth_tag = $4 WHERE id = $5",
			env.EncryptedValue, env.EncryptedDEK, env.IV, env.AuthTag, id)
		if err != nil {
			log.Fatalf("Failed to update secret_version %s: %v", id, err)
		}
		count++
	}
	log.Printf("Migrated %d secret versions.\n", count)
}

func legacyDecrypt(key []byte, cipherText string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, encrypted := decoded[:nonceSize], decoded[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
