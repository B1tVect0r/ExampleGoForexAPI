package aes

import (
	"github.com/b1tvect0r/exchangerates/pkg/auth/apikeys"
)

type aesAPIKeyProvider struct {
	key []byte
}

func (akp *aesAPIKeyProvider) Create(projectID string) (string, error) {
	return "", nil
}

func (akp *aesAPIKeyProvider) Verify(apiKey string) error {
	return nil
}

func New(signingKey string) (apikeys.APIKeyProvider, error) {
	return &aesAPIKeyProvider{[]byte(signingKey)}, nil
}
