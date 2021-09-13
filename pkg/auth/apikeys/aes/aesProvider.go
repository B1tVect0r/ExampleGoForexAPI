package aes

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/b1tvect0r/exchangerates/pkg/auth/apikeys"
	"github.com/b1tvect0r/exchangerates/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

type aesAPIKeyProvider struct {
	key []byte
}

func (akp *aesAPIKeyProvider) encrypt(plaintext []byte) ([]byte, error) {
	c, err := aes.NewCipher(akp.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM operation mode: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to create nonce: %w", err)
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (akp *aesAPIKeyProvider) decrypt(ciphertext []byte) ([]byte, error) {
	c, err := aes.NewCipher(akp.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM operation mode: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext cannot be shorter than nonce size (at least %d, was %d)", nonceSize, len(ciphertext))
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

func (akp *aesAPIKeyProvider) Create(ctx context.Context, q *db.Queries, kp apikeys.APIKeyParameters) (string, error) {
	hashedSecret, err := bcrypt.GenerateFromPassword(kp.Secret, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash key secret: %w", err)
	}

	kpBytes, err := json.Marshal(&kp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal key params to JSON: %w", err)
	}

	// Encrypt. This makes the secret opaque to the outside world.
	enc, err := akp.encrypt(kpBytes)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt: %w", err)
	}

	_, err = q.SetProjectSecret(ctx, db.SetProjectSecretParams{ProjectID: kp.ProjectID, HashedSecret: hashedSecret})
	if err != nil {
		return "", fmt.Errorf("failed to save hashed secret to database: %w", err)
	}

	// Base-64 encode the encoded bytes
	return base64.URLEncoding.Strict().EncodeToString(enc), nil
}

func (akp *aesAPIKeyProvider) Verify(ctx context.Context, q *db.Queries, apiKey string) error {
	kpBytes, err := base64.URLEncoding.Strict().DecodeString(apiKey)
	if err != nil {
		return fmt.Errorf("failed to decode api key: %w", err)
	}

	dec, err := akp.decrypt(kpBytes)
	if err != nil {
		return fmt.Errorf("failed to decrypt api key: %w", err)
	}

	kp := apikeys.APIKeyParameters{}
	if err = json.Unmarshal(dec, &kp); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	hashedSecret, err := q.GetProjectSecret(ctx, kp.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to retrieve project secret: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword(hashedSecret, kp.Secret); err != nil {
		return fmt.Errorf("API key secret mismatch: %w", err)
	}

	return nil
}

// New returns an AES-driven API key provider that will sign/decrypt secrets using the provided key.
func New(signingKey string) (apikeys.APIKeyProvider, error) {
	return &aesAPIKeyProvider{[]byte(signingKey)}, nil
}
