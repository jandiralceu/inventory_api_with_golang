package platform

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword_Success(t *testing.T) {
	password := "my_secure_password"

	hash, err := hashPassword(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Check Argon2id format
	assert.True(t, strings.HasPrefix(hash, "$argon2id$v="), "hash should start with argon2id prefix")

	// Ensure the hash has 6 parts split by '$'
	parts := strings.Split(hash, "$")
	assert.Equal(t, 6, len(parts), "hash should have 6 parts")
}

func TestVerifyPassword_Success(t *testing.T) {
	password := "my_secure_password"
	hash, err := hashPassword(password)
	require.NoError(t, err)

	match, err := verifyPassword(password, hash)

	require.NoError(t, err)
	assert.True(t, match, "password should match the hash")
}

func TestVerifyPassword_WrongPassword(t *testing.T) {
	password := "my_secure_password"
	wrongPassword := "wrong_password"
	hash, err := hashPassword(password)
	require.NoError(t, err)

	match, err := verifyPassword(wrongPassword, hash)

	require.NoError(t, err)
	assert.False(t, match, "wrong password should not match the hash")
}

func TestVerifyPassword_InvalidHashFormat(t *testing.T) {
	password := "my_secure_password"
	invalidHash := "$argon2id$v=19$m=65536,t=1,p=4$invalid_salt" // Missing hash part

	match, err := verifyPassword(password, invalidHash)

	require.Error(t, err)
	assert.Equal(t, ErrInvalidHash, err)
	assert.False(t, match)
}

func TestVerifyPassword_InvalidBase64Salt(t *testing.T) {
	password := "my_secure_password"
	// 5 parts correctly separated, but part 4 (salt) is invalid base64 (e.g., contains '?' or length is weird)
	invalidHash := "$argon2id$v=19$m=65536,t=1,p=4$@invalid_salt@$c29tZWhhc2g="

	match, err := verifyPassword(password, invalidHash)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode salt")
	assert.False(t, match)
}

func TestVerifyPassword_InvalidBase64Hash(t *testing.T) {
	password := "my_secure_password"
	invalidHash := "$argon2id$v=19$m=65536,t=1,p=4$c29tZXNhbHQ$@invalid_hash@"

	match, err := verifyPassword(password, invalidHash)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode hash")
	assert.False(t, match)
}
