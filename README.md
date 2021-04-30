# mdb-bench

`mdb-bench` includes all benchmarks and tests for the MicroDB research.

## Setup

To run the benchmark you will need to following installed,

- Docker
- Make
- Go `>= 1.16.3`
- [vegeta](https://github.com/flaviostutz/vegeta) (optional, for plotting graphs)

## Shipping

Shipping contains a sample service-oriented application based on `go-kit`'s example.

## Docker

Contains all dockerfile, docker-compose files for spinning up a test environment.

For example, a mysql-cluster with the application,

```bash
make run-backend

make run-mysql-cluster

make run-app-mysql-cluster
```

## Benchmark

Using the cli tool under `/bench`, you can benchmark with various load, frequency, and duration.

```bash
go run bench/main.go \
    --url http://localhost:8080 \
    --freq 5000 \
    --dur 5m \
    --name example-test-run \
    --spike \
    --attack \
    --save
```

This above example will start a benchmark against `http://localhost:8080` with load of 5000 rps for a duration of 5 minutes. It will be a load test containing triggers for sudden database spikes. At the end of the test run, it will save the results under `./results/example-test-run`.
