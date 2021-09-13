package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const authorizationHeaderKey = "Authorization"
const authorizationType = "Bearer"

func VerifyAPIKeyMiddleware(s *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var apiKey string
		authHeader := c.Request.Header.Get(authorizationHeaderKey)
		headerFrags := strings.Split(authHeader, " ")
		if len(headerFrags) == 2 && headerFrags[0] == authorizationType {
			apiKey = headerFrags[1]
		} else {
			log.Printf("Incorrect number of fragments or incorrect auth type: %v", headerFrags)
		}

		if apiKey == "" {
			c.Writer.Header().Add("WWW-Authenticate", authorizationType)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err := s.kp.Verify(c.Request.Context(), s.dal, apiKey); err != nil {
			c.AbortWithError(http.StatusForbidden, fmt.Errorf("failed to verify API key: %w", err))
			return
		}

		c.Next()
	}
}
