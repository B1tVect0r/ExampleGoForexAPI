package server

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/b1tvect0r/exchangerates/pkg/auth/apikeys"
	"github.com/gin-gonic/gin"
)

var projectIDRunes = []rune("ABCDEFGHJKLMNPQRSTUVWXYZ")

func makeProjectID(length int) string {
	r := make([]rune, length)
	for i := range r {
		r[i] = projectIDRunes[rand.Intn(len(projectIDRunes))]
	}
	return string(r)
}

type CreateProjectOutput struct {
	APIKey string `json:"APIKey"`
}

const projectIDLength = 10

func (s *Server) CreateProject(c *gin.Context) {
	pID := makeProjectID(projectIDLength)
	pSecret := apikeys.MakeSecret()

	key, err := s.kp.Create(c.Request.Context(), s.dal, apikeys.APIKeyParameters{ProjectID: pID, Secret: pSecret})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create API key: %w", err))
		return
	}

	c.JSON(http.StatusOK, &CreateProjectOutput{key})
}
