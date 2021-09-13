package apikeys

type APIKeyProvider interface {
	Create(projectID string) (string, error)
	Verify(apiKey string) error
}
