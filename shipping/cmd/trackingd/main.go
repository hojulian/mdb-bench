package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/hojulian/mdb-bench/shipping/cargo"
	"github.com/hojulian/mdb-bench/shipping/database"
	"github.com/hojulian/mdb-bench/shipping/tracking"
)

func main() {
	var (
		addr         = envString("PORT", "8080")
		httpAddr     = flag.String("http.addr", ":"+addr, "HTTP listen address")
		databaseType = flag.String("database", "inmem", "database type")
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// Create repos
	cargos, handlingEvents, err := repos(*databaseType)
	if err != nil {
		logger.Log("error_msg", "failed to create repos", "error", err)
		return
	}

	var ts tracking.Service
	fieldKeys := []string{"method"}
	ts = tracking.NewService(cargos, handlingEvents)
	ts = tracking.NewLoggingService(log.With(logger, "component", "tracking"), ts)
	ts = tracking.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "tracking_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "tracking_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		ts,
	)

	// Start service
	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/tracking/v1/", tracking.MakeHandler(ts, httpLogger))
	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func repos(databaseType string) (cargo.Repository, cargo.HandlingEventRepository, error) {
	t := database.DatabaseType(databaseType)
	params := map[string]string{
		"MYSQL_NODES":     envString("MYSQL_NODES", "2"),
		"MYSQL_HOST":      envString("MYSQL_HOST", "127.0.0.1"),
		"MYSQL_PORT":      envString("MYSQL_PORT", "3306"),
		"MYSQL_USER":      envString("MYSQL_USER", "root"),
		"MYSQL_PASSWORD":  envString("MYSQL_PASSWORD", "test"),
		"MYSQL_DATABASE":  envString("MYSQL_DATABASE", "test"),
		"NATS_HOST":       envString("NATS_HOST", "127.0.0.1"),
		"NATS_PORT":       envString("NATS_PORT", "4222"),
		"NATS_CLIENT_ID":  envString("NATS_CLIENT_ID", "bookingd-client"),
		"NATS_CLUSTER_ID": envString("NATS_CLUSTER_ID", "nats-cluster"),
	}

	cargos, err := database.NewCargoRepository(t, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to creat cargos repo: %w", err)
	}

	handlingEvents, err := database.NewHandlingEventRepository(t, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to creat handling events repo: %w", err)
	}

	return cargos, handlingEvents, nil
}
