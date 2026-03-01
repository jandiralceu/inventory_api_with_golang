package pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to generate RSA test keys in memory without saving to files
func generateMemRSAKeys(t *testing.T) (string, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	return string(privPEM), string(pubPEM)
}

func TestNewJWTManager_Success(t *testing.T) {
	privPEM, pubPEM := generateMemRSAKeys(t)

	manager, err := NewJWTManager(privPEM, pubPEM)

	require.NoError(t, err)
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.privateKey)
	assert.NotNil(t, manager.publicKey)
}

func TestNewJWTManager_InvalidPrivateKey(t *testing.T) {
	_, pubPEM := generateMemRSAKeys(t)

	manager, err := NewJWTManager("invalid_private_key_pem_format", pubPEM)

	require.Error(t, err)
	assert.Nil(t, manager)
	assert.Contains(t, err.Error(), "private key")
}

func TestNewJWTManager_InvalidPublicKey(t *testing.T) {
	privPEM, _ := generateMemRSAKeys(t)

	manager, err := NewJWTManager(privPEM, "invalid_public_key_pem_format")

	require.Error(t, err)
	assert.Nil(t, manager)
	assert.Contains(t, err.Error(), "public key")
}

func TestJWTManager_GenerateAndValidateToken_Success(t *testing.T) {
	privPEM, pubPEM := generateMemRSAKeys(t)
	manager, err := NewJWTManager(privPEM, pubPEM)
	require.NoError(t, err)

	userID := uuid.New()
	expiration := time.Minute * 15

	role := "admin"
	token, err := manager.GenerateToken(userID, role, expiration)

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := manager.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
}

func TestJWTManager_ValidateToken_Expired(t *testing.T) {
	privPEM, pubPEM := generateMemRSAKeys(t)
	manager, err := NewJWTManager(privPEM, pubPEM)
	require.NoError(t, err)

	userID := uuid.New()
	expiration := -time.Minute * 15 // Set an expiration in the past

	token, err := manager.GenerateToken(userID, "admin", expiration)
	require.NoError(t, err)

	claims, err := manager.ValidateToken(token)
	require.Error(t, err)
	// Should fail due to "token is expired" from golang-jwt validator
	assert.Nil(t, claims)
}

func TestJWTManager_ValidateToken_InvalidSignature(t *testing.T) {
	privPEM1, pubPEM1 := generateMemRSAKeys(t)
	manager1, err := NewJWTManager(privPEM1, pubPEM1)
	require.NoError(t, err)

	privPEM2, pubPEM2 := generateMemRSAKeys(t)
	manager2, err := NewJWTManager(privPEM2, pubPEM2)
	require.NoError(t, err)

	userID := uuid.New()
	tokenFrom1, err := manager1.GenerateToken(userID, "admin", time.Hour)
	require.NoError(t, err)

	// Validating token signed by Key 1 using Key 2, should fail signature validation
	claims, err := manager2.ValidateToken(tokenFrom1)
	require.Error(t, err)
	assert.Nil(t, claims)
}

func BenchmarkGenerateToken(b *testing.B) {
	privPEM, pubPEM := generateMemRSAKeysBenchmark(b)
	manager, _ := NewJWTManager(privPEM, pubPEM)
	userID := uuid.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.GenerateToken(userID, "admin", time.Hour)
	}
}

func BenchmarkValidateToken(b *testing.B) {
	privPEM, pubPEM := generateMemRSAKeysBenchmark(b)
	manager, _ := NewJWTManager(privPEM, pubPEM)
	userID := uuid.New()
	token, _ := manager.GenerateToken(userID, "admin", time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = manager.ValidateToken(token)
	}
}

// generateMemRSAKeysBenchmark is a helper for benchmarks
func generateMemRSAKeysBenchmark(b *testing.B) (string, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		b.Fatal(err)
	}

	privBytes, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	pubBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})

	return string(privPEM), string(pubPEM)
}
