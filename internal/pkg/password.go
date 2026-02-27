package platform

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordHasher abstracts password hashing and verification.
// This allows injecting different implementations (e.g., identity hasher for tests)
// and facilitates future algorithm rotations.
type PasswordHasher interface {
	// Hash generates a cryptographically secure hash of the provided password.
	Hash(password string) (string, error)
	// Verify checks if the provided password matches the encoded hash.
	Verify(password, hash string) (bool, error)
}

// Argon2PasswordHasher implements PasswordHasher using Argon2id,
// the winner of the Password Hashing Competition.
type Argon2PasswordHasher struct{}

// NewArgon2PasswordHasher returns a PasswordHasher backed by Argon2id.
func NewArgon2PasswordHasher() PasswordHasher {
	return &Argon2PasswordHasher{}
}

// Hash invokes the internal Argon2id hashing logic.
func (h *Argon2PasswordHasher) Hash(password string) (string, error) {
	return hashPassword(password)
}

// Verify invokes the internal Argon2id verification logic.
func (h *Argon2PasswordHasher) Verify(password, hash string) (bool, error) {
	return verifyPassword(password, hash)
}

// Argon2id parameters following OWASP recommendations for balanced security/performance.
// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
var (
	argonTime    uint32 = 1
	argonMemory  uint32 = 64 * 1024 // 64 MB
	argonThreads uint8  = 4
	argonKeyLen  uint32 = 32
	argonSaltLen        = 16
)

// ErrInvalidHash is returned when the hash string doesn't follow the Argon2 standard format.
var ErrInvalidHash = errors.New("invalid password hash format")

// hashPassword hashes a plain-text password using Argon2id with a unique 16-byte salt.
// Returns a string in the format: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
func hashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemory, argonTime, argonThreads,
		encodedSalt, encodedHash,
	), nil
}

// verifyPassword checks if a plain-text password matches an Argon2id hash.
// It extracts the salt from the encoded hash and recomputes the key for comparison.
func verifyPassword(password, encodedHash string) (bool, error) {
	salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// Use constant-time comparison to prevent timing attacks where attackers
	// could deduce password length or prefix by measuring response times.
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

// decodeHash extracts the salt and hash from an encoded Argon2id string.
func decodeHash(encodedHash string) ([]byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode hash: %w", err)
	}

	return salt, hash, nil
}
