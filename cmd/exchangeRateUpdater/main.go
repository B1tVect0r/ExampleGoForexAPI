package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/b1tvect0r/exchangerates/pkg/db"
	"github.com/b1tvect0r/exchangerates/pkg/ratefetcher"
	"github.com/jackc/pgx/v4/pgxpool"
)

func updateRates(rf ratefetcher.RateFetcher, ctx context.Context) error {
	segments, err := rf.MakeExchangeSegments(ctx)
	if err != nil {
		return fmt.Errorf("failed to make exchange segments: %w", err)
	}

	exchangeRates := make([]db.SetExchangeRateParams, 0, (len(segments)*len(segments))/2)

	for _, segment := range segments {
		ratesForSegment, err := rf.FetchRates(ctx, segment)
		if err != nil {
			return err
		}

		log.Printf("Retrieved rates: %v", ratesForSegment)
		exchangeRates = append(exchangeRates, ratesForSegment...)
	}

	for _, r := range exchangeRates {
		log.Printf("Writing rate %v to database", r)
		if err = rf.StoreRate(ctx, r); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	ctx := context.Background()

	dbConnString := os.Getenv("PSQL_CONNECTION_STRING")
	pool, err := pgxpool.Connect(ctx, dbConnString)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to establish connection pool to database at %s", dbConnString))
	}
	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database at %s", dbConnString)
	}

	q := db.New(pool)

	rf, err := ratefetcher.Default(q, os.Getenv("FIXER_API_KEY"))
	if err != nil {
		log.Fatalf("Failed to initialize rate fetcher: %s", err.Error())
	}

	for {
		if err = updateRates(rf, ctx); err != nil {
			log.Printf("FAILED TO UPDATE RATES: %s", err.Error())
		}
		log.Printf("Done updating rates; sleeping for an hour.")
		time.Sleep(1 * time.Hour)
	}
}
