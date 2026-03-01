package pkg

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken  = errors.New("invalid or expired token")
	ErrInvalidKey    = errors.New("failed to parse RSA key")
	ErrInvalidClaims = errors.New("invalid token claims")
)

// JWTManager handles JWT token generation and validation using asymmetric RSA keys.
// It uses RS256 for signing tokens.
type JWTManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// UserClaims represents the data stored in the JWT claims.
type UserClaims struct {
	UserID uuid.UUID
	Role   string
}

// NewJWTManager creates a new JWTManager by parsing PEM-encoded RSA keys.
// The private key is required for signing, and the public key for verification.
func NewJWTManager(privateKeyPEM, publicKeyPEM string) (*JWTManager, error) {
	privKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("private key: %w", err)
	}

	pubKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("public key: %w", err)
	}

	return &JWTManager{
		privateKey: privKey,
		publicKey:  pubKey,
	}, nil
}

// GenerateToken creates a signed JWT containing the user ID in the 'sub' claim.
// It also includes standard 'iat' (issued at), 'exp' (expiration), and 'jti' (token ID) claims.
func (j *JWTManager) GenerateToken(userID uuid.UUID, role string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(expiration).Unix(),
		"jti":  uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken verifies the token signature against the public key and expiration.
// Returns the user ID extracted from the 'sub' claim if valid.
func (j *JWTManager) ValidateToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, ErrInvalidClaims
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return nil, ErrInvalidClaims
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return &UserClaims{
		UserID: userID,
		Role:   role,
	}, nil
}

func parsePrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, ErrInvalidKey
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidKey, err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%w: not an RSA key", ErrInvalidKey)
	}

	return rsaKey, nil
}

func parsePublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, ErrInvalidKey
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidKey, err)
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("%w: not an RSA key", ErrInvalidKey)
	}

	return rsaKey, nil
}
