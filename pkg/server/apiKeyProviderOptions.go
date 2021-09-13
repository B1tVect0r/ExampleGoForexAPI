package server

import (
	"github.com/b1tvect0r/exchangerates/pkg/auth/apikeys/aes"
)

// WithAESAPIKeyProvider configures the server to use the AES-backed key provider.
func WithAESAPIKeyProvider(signingKey string) func(*Server) error {
	return func(s *Server) error {
		akp, err := aes.New(signingKey)
		if err != nil {
			return err
		}

		s.kp = akp

		return nil
	}
}
