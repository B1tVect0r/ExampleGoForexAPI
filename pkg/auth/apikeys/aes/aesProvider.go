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

func (akp *aesAPIKeyProvider) Create(ctx context.Context, q *db.Queries, kp apikeys.APIKeyParameters) (string, error) {
	c, err := aes.NewCipher(akp.key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM operation mode: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to create nonce: %w", err)
	}

	hashedSecret, err := bcrypt.GenerateFromPassword(kp.Secret, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash key secret: %w", err)
	}

	kpBytes, err := json.Marshal(&kp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal key params to JSON: %w", err)
	}

	_, err = q.SetProjectSecret(ctx, db.SetProjectSecretParams{ProjectID: kp.ProjectID, HashedSecret: hashedSecret})
	if err != nil {
		return "", fmt.Errorf("failed to save hashed secret to database: %w", err)
	}

	// Seal & base-64 encode the key parameters. This makes it opaque to the outside world.
	return base64.URLEncoding.Strict().EncodeToString(gcm.Seal(nonce, nonce, kpBytes, nil)), nil
}

func (akp *aesAPIKeyProvider) Verify(ctx context.Context, q *db.Queries, apiKey string) error {
	kpBytes, err := base64.URLEncoding.Strict().DecodeString(apiKey)
	if err != nil {
		return fmt.Errorf("failed to decode api key: %w", err)
	}

	var kp *apikeys.APIKeyParameters
	if err = json.Unmarshal(kpBytes, kp); err != nil {
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

func New(signingKey string) (apikeys.APIKeyProvider, error) {
	return &aesAPIKeyProvider{[]byte(signingKey)}, nil
}
