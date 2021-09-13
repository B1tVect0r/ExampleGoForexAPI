package ratefetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/b1tvect0r/exchangerates/pkg/db"
)

type ExchangeSegment struct {
	FromCurrency string
	ToCurrencies []string
}

type RateFetcher interface {
	MakeExchangeSegments(ctx context.Context) ([]ExchangeSegment, error)
	FetchRates(ctx context.Context, es ExchangeSegment) ([]db.SetExchangeRateParams, error)
	StoreRate(ctx context.Context, er db.SetExchangeRateParams) error
}

const fixerAPIRatesEndpoint = "http://data.fixer.io/api/latest"

type fixerRates struct {
	Base  string             `json:"Base"`
	Rates map[string]float32 `json:"Rates"`
}

type defaultRateFetcher struct {
	*db.Queries
	fixerAPIKey string
}

func Default(q *db.Queries, fixerAPIKey string) (RateFetcher, error) {
	return &defaultRateFetcher{q, fixerAPIKey}, nil
}

func (drf *defaultRateFetcher) MakeExchangeSegments(ctx context.Context) ([]ExchangeSegment, error) {
	currencies, err := drf.GetCurrencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currencies: %w", err)
	}

	log.Printf("Retrieved currencies: %v", currencies)

	segments := make([]ExchangeSegment, len(currencies)-1)
	// We assume that the exchange rate from B=>A is the inverse of the rate from A=>B,
	// so we don't need to fetch every rate and can infer some instead.
	for i := range segments {
		segments[i] = ExchangeSegment{currencies[i], currencies[i+1:]}
	}

	log.Printf("Created segments %v", segments)

	return segments, nil
}

func (drf *defaultRateFetcher) FetchRates(ctx context.Context, ep ExchangeSegment) ([]db.SetExchangeRateParams, error) {
	log.Printf("Fetching rates for segment %v", ep)
	rates, err := http.Get(fmt.Sprintf(
		"%s?access_key=%s&base=%s&symbols=%s",
		fixerAPIRatesEndpoint,
		drf.fixerAPIKey,
		ep.FromCurrency,
		strings.Join(ep.ToCurrencies, ","),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates for %s: %w", ep.FromCurrency, err)
	}

	defer func() {
		_ = rates.Body.Close()
	}()

	bytes, err := ioutil.ReadAll(rates.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	fr := fixerRates{}
	if err = json.Unmarshal(bytes, &fr); err != nil {
		return nil, fmt.Errorf("failed to parse into json: %w", err)
	}

	out := make([]db.SetExchangeRateParams, 0, len(ep.ToCurrencies))
	for _, oc := range ep.ToCurrencies {
		if r, ok := fr.Rates[oc]; ok {
			out = append(out, db.SetExchangeRateParams{
				FromCurrency: ep.FromCurrency,
				ToCurrency:   oc,
				Rate:         r,
			})
		}
	}

	return out, nil
}

func (drf *defaultRateFetcher) StoreRate(ctx context.Context, erp db.SetExchangeRateParams) error {
	if err := drf.SetExchangeRate(ctx, erp); err != nil {
		return fmt.Errorf("failed to store exchange rate: %w", err)
	}

	return nil
}
