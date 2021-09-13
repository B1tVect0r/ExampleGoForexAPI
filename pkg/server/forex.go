package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

const (
	fromCurParam = "fromCur"
	toCursParam  = "toCur"
)

func (s *Server) ExchangeRatesForCurrency(c *gin.Context) {
	from := c.Param(fromCurParam)
	toCurs := c.QueryArray(toCursParam)

	log.Printf("Hello from ExchangeRatesForCurrency with from='%s', toCur='%v' (len %d)", from, toCurs, len(toCurs))
}
