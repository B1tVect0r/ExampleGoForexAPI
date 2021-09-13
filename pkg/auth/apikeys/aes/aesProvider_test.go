package aes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testKey = "16bytekeyfortest"

func TestMustHaveAppropriateKey(t *testing.T) {
	ap := &aesAPIKeyProvider{[]byte("too short")}
	_, err := ap.encrypt([]byte("plaintext"))
	assert.Error(t, err)

	_, err = ap.decrypt([]byte("ciphertext"))
	assert.Error(t, err)
}

func TestCanEncryptAndDecrypt(t *testing.T) {
	ap := &aesAPIKeyProvider{[]byte(testKey)}
	plaintext := "some secret"

	ciphertext, err := ap.encrypt([]byte(plaintext))
	assert.NoError(t, err)

	ptbytes, err := ap.decrypt(ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, string(ptbytes))
}
