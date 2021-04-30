package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/hojulian/mdb-bench/shipping/booking"
	"github.com/hojulian/mdb-bench/shipping/cargo"
	"github.com/hojulian/mdb-bench/shipping/database"
	"github.com/hojulian/mdb-bench/shipping/location"
	"github.com/hojulian/mdb-bench/shipping/routing"
	"github.com/hojulian/mdb-bench/shipping/voyage"
)

var locations = []location.UNLocode{
	location.SESTO,
	location.AUMEL,
	location.CNHKG,
	location.USNYC,
	location.USCHI,
	location.JNTKO,
	location.DEHAM,
	location.NLRTM,
	location.FIHEL,
}

var (
	cargosCount = flag.Int("cargos", 10000, "number of cargos")
)

func main() {
	var (
		rsurl             = envString("ROUTINGSERVICE_URL", "http://localhost:7878")
		routingServiceURL = flag.String("service.routing", rsurl, "routing service URL")
		ctx               = context.Background()
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// Create repos
	cargos, locations, _, handlingEvents, err := repos()
	if err != nil {
		panic(fmt.Errorf("failed to create repos: %w", err))
	}

	var rs routing.Service
	rs = routing.NewProxyingMiddleware(ctx, *routingServiceURL)(rs)

	bs := booking.NewService(cargos, locations, handlingEvents, rs)

	// Seed
	logger.Log("stage", "cargos")
	for i := 0; i < *cargosCount; i++ {
		if err := bookRandomCargo(bs); err != nil {
			logger.Log("seed", "failed", "error", err)
			return
		}
		if i%5000 == 0 {
			logger.Log("cargos", i)
		}
	}
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func repos() (cargo.Repository, location.Repository, voyage.Repository, cargo.HandlingEventRepository, error) {
	params := map[string]string{
		"MYSQL_HOST":     envString("MYSQL_HOST", "127.0.0.1"),
		"MYSQL_PORT":     envString("MYSQL_PORT", "3306"),
		"MYSQL_USER":     envString("MYSQL_USER", "root"),
		"MYSQL_PASSWORD": envString("MYSQL_PASSWORD", "test"),
		"MYSQL_DATABASE": envString("MYSQL_DATABASE", "test"),
	}

	cargos, err := database.NewMySQLCargoRepository(params, false)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create cargos repo: %w", err)
	}

	locations, err := database.NewMySQLLocationRepository(params, false)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create locations repo: %w", err)
	}

	voyages, err := database.NewMySQLVoyageRepository(params, false)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create voyages repo: %w", err)
	}

	handlingEvents, err := database.NewMySQLHandlingEventRepository(params, false)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create handling events repo: %w", err)
	}

	return cargos, locations, voyages, handlingEvents, nil
}

func bookRandomCargo(bs booking.Service) error {
	o := randomLoc()
	d := randomLoc()
	t := time.Now().AddDate(0, 0, rand.Intn(30))
	_, err := bs.BookNewCargo(o, d, t)
	if err != nil {
		return fmt.Errorf("failed to create cargo: %w", err)
	}
	return err
}

func randomLoc() location.UNLocode {
	return locations[rand.Intn(len(locations))]
}
