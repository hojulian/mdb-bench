package main

import (
	"flag"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"

	"github.com/hojulian/mdb-bench/bench/shipping"
)

var (
	url  = flag.String("url", "http://localhost:8080", "target url")
	freq = flag.Int("freq", 100, "requests per second")
	dur  = flag.Duration("dur", 1*time.Minute, "duration in minutes")
	// out  = flag.String("out", "./results", "results output directory")
)

var users = []shipping.User{
	shipping.RegularCustomer(*url),
	shipping.RegularCustomer(*url),
	shipping.RegularCustomer(*url),
	shipping.RegularAuditor(*url),
	shipping.RegularBooker(*url),
}

func main() {
	flag.Parse()

	var targets []vegeta.Target
	for _, u := range users {
		targets = append(targets, u.Interactions()...)
	}

	targeter := vegeta.NewStaticTargeter(targets...)
	attacker := vegeta.NewAttacker()
	rate := vegeta.Rate{Freq: *freq, Per: time.Second}

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, *dur, "Normal load") {
		metrics.Add(res)
	}
	metrics.Close()

	rpt := vegeta.NewTextReporter(&metrics)
	rpt.Report(os.Stdout)
}
