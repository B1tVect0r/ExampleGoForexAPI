package server

import (
	"github.com/b1tvect0r/exchangerates/pkg/auth/apikeys/aes"
)

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
