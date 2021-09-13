package apikeys

import (
	"context"
	"math/rand"

	"github.com/b1tvect0r/exchangerates/pkg/db"
)

// APIKeyParameters wraps a project id and secret for json marshaling/unmarshaling.
type APIKeyParameters struct {
	ProjectID string `json:"ProjectID"`
	Secret    []byte `json:"Secret"`
}

// APIKeyProvider creates and verifies API keys.
// Todo: Replace the reference to concrete db.Queries with an appropriate interface for testability
type APIKeyProvider interface {
	Create(ctx context.Context, q *db.Queries, kp APIKeyParameters) (string, error)
	Verify(ctx context.Context, q *db.Queries, apiKey string) error
}

var secretRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const secretLength = 32

// MakeSecret creates a random 32-byte array consisting of standard ASCII alphanum characters.
func MakeSecret() []byte {
	r := make([]rune, secretLength)
	for i := range r {
		r[i] = secretRunes[rand.Intn(len(secretRunes))]
	}
	return []byte(string(r))
}
