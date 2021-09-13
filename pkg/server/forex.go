package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/b1tvect0r/exchangerates/pkg/db"
	"github.com/gin-gonic/gin"
)

const (
	fromCurParam = "fromCur"
	toCursParam  = "toCur"
)

func (s *Server) ExchangeRatesForCurrency(c *gin.Context) {
	from := strings.ToUpper(c.Param(fromCurParam))
	toCurs := c.QueryArray(toCursParam)

	for i, c := range toCurs {
		toCurs[i] = strings.ToUpper(c)
	}

	var rates []db.ExchangeRate
	var err error

	log.Printf("Making query for rates going from %s to %v", from, toCurs)

	if len(toCurs) > 0 {
		rates, err = s.dal.GetExchangeRatesForCurrency(c.Request.Context(), db.GetExchangeRatesForCurrencyParams{
			FromCurrency: from,
			ToCurrency:   toCurs,
		})
	} else {
		rates, err = s.dal.GetAllExchangeRatesForCurrency(c.Request.Context(), from)
	}

	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve exchange rates: %w", err))
		return
	}

	log.Printf("Retrieved %d rates going from %s to %v", len(rates), from, toCurs)

	c.JSON(http.StatusOK, rates)
}
