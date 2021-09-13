package server

import (
	"fmt"

	"github.com/b1tvect0r/exchangerates/pkg/auth/apikeys"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	kp apikeys.APIKeyProvider
}

type Opt func(*Server) error

func New(opts ...Opt) (*Server, error) {
	s := &Server{Engine: gin.Default()}
	s.withDefaultRoutes()

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, fmt.Errorf("failed to apply option to server: %w", err)
		}
	}

	if s.kp == nil {
		return nil, fmt.Errorf("server must have an APIKeyProvider; please provide one using one of the WithXAPIKeyProvider options")
	}

	return s, nil
}
