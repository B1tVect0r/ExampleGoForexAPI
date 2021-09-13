package apikeys

import (
	"context"
	"math/rand"

	"github.com/b1tvect0r/exchangerates/pkg/db"
)

type APIKeyParameters struct {
	ProjectID string `json:"ProjectID"`
	Secret    []byte `json:"Secret"`
}

type APIKeyProvider interface {
	Create(ctx context.Context, q *db.Queries, kp APIKeyParameters) (string, error)
	Verify(ctx context.Context, q *db.Queries, apiKey string) error
}

var secretRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const secretLength = 32

func MakeSecret() []byte {
	r := make([]rune, secretLength)
	for i := range r {
		r[i] = secretRunes[rand.Intn(len(secretRunes))]
	}
	return []byte(string(r))
}
