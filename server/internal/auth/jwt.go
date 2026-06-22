package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	privateKey      *rsa.PrivateKey
	publicKey       *rsa.PublicKey
	accessTokenTTL  time.Duration
}

type TokenClaims struct {
	UserID       string `json:"uid"`
	Email        string `json:"email"`
	SessionID    string `json:"sid"`
	IsSuperAdmin bool   `json:"sa,omitempty"`
	IsReauth     bool   `json:"reauth,omitempty"`
	jwt.RegisteredClaims
}

func NewJWTService(keyDir string, accessTTL time.Duration) (*JWTService, error) {
	privKey, pubKey, err := loadOrGenerateKeys(keyDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWT keys: %w", err)
	}

	return &JWTService{
		privateKey:     privKey,
		publicKey:      pubKey,
		accessTokenTTL: accessTTL,
	}, nil
}

func (s *JWTService) GenerateAccessToken(userID, email, sessionID string, isSuperAdmin bool) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID:       userID,
		Email:        email,
		SessionID:    sessionID,
		IsSuperAdmin: isSuperAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			Issuer:    "env-vault",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *JWTService) GenerateReauthToken(userID, email, sessionID string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID:       userID,
		Email:        email,
		SessionID:    sessionID,
		IsReauth:     true,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)), // Short-lived
			Issuer:    "env-vault",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *JWTService) ValidateAccessToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func loadOrGenerateKeys(keyDir string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return nil, nil, err
	}

	privPath := filepath.Join(keyDir, "jwt_private.pem")
	pubPath := filepath.Join(keyDir, "jwt_public.pem")

	if _, err := os.Stat(privPath); os.IsNotExist(err) {
		return generateAndSaveKeys(privPath, pubPath)
	}

	return loadKeys(privPath, pubPath)
}

func generateAndSaveKeys(privPath, pubPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
	if err := os.WriteFile(privPath, privPEM, 0600); err != nil {
		return nil, nil, fmt.Errorf("failed to write private key: %w", err)
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	if err := os.WriteFile(pubPath, pubPEM, 0644); err != nil {
		return nil, nil, fmt.Errorf("failed to write public key: %w", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}

func loadKeys(privPath, pubPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privData, err := os.ReadFile(privPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read private key: %w", err)
	}
	privBlock, _ := pem.Decode(privData)
	if privBlock == nil {
		return nil, nil, errors.New("failed to decode private key PEM")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	pubData, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read public key: %w", err)
	}
	pubBlock, _ := pem.Decode(pubData)
	if pubBlock == nil {
		return nil, nil, errors.New("failed to decode public key PEM")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	publicKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("public key is not RSA")
	}

	return privateKey, publicKey, nil
}
