package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"

	"github.com/hojulian/mdb-bench/bench/shipping"
)

var (
	url       = flag.String("url", "http://localhost:8080", "target url")
	freq      = flag.Int("freq", 100, "requests per second")
	spike     = flag.Bool("spike", false, "create database spike")
	spikeFreq = flag.Int("spike-freq", 3, "spikes per minute")
	dur       = flag.Duration("dur", 1*time.Minute, "duration in minutes")
	name      = flag.String("name", "example-run", "name of the test")
	attack    = flag.Bool("attack", false, "run attack")
	save      = flag.Bool("save", false, "save results")
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
	rand.Shuffle(len(targets), func(i, j int) { targets[i], targets[j] = targets[j], targets[i] })

	// Create results directory
	if err := os.MkdirAll(fmt.Sprintf("./results/%s", *name), os.ModePerm); err != nil {
		log.Fatalf("failed to create results directory: %s", err)
	}

	// Encode targets
	tf, err := os.Create(fmt.Sprintf("./results/%s/targets.json", *name))
	if err != nil {
		log.Fatalf("failed to create targets file: %s", err)
	}
	defer tf.Close()

	e := vegeta.NewJSONTargetEncoder(tf)
	for _, t := range targets {
		if err := e.Encode(&t); err != nil {
			log.Printf("failed to encode target to file: %s", err)
			return
		}
	}

	s := &strings.Builder{}
	wg := &sync.WaitGroup{}
	if *spike {
		wg.Add(1)
		log.Println("Running with SPIKES")
		go Spikes(wg, s)
	} else {
		log.Println("Running")
	}

	if *attack {
		f, err := os.Create(fmt.Sprintf("./results/%s/results.json", *name))
		if err != nil {
			log.Fatalf("failed to open results directory: %s", err)
		}
		defer f.Close()

		encoder := vegeta.NewJSONEncoder(f)
		targeter := vegeta.NewStaticTargeter(targets...)
		attacker := vegeta.NewAttacker()
		rate := vegeta.Rate{Freq: *freq, Per: time.Second}

		var metrics vegeta.Metrics
		for res := range attacker.Attack(targeter, rate, *dur, *name) {
			metrics.Add(res)
			if *save {
				if err := encoder.Encode(res); err != nil {
					log.Printf("failed to encode result to file: %s", err)
					return
				}
			}
		}
		metrics.Close()

		rpt := vegeta.NewTextReporter(&metrics)
		rpt.Report(os.Stdout)
	}

	log.Println("Waiting for other attackers to finish...")
	wg.Wait()
	log.Printf("===== SPIKES =====\n%s", s.String())
	log.Println("Done.")
}

func Spikes(wg *sync.WaitGroup, out *strings.Builder) {
	targets := shipping.HighLoadBooker(*url).Interactions()
	attacker := vegeta.NewAttacker()
	targeter := vegeta.NewStaticTargeter(targets...)
	rate := vegeta.Rate{Freq: *spikeFreq, Per: time.Minute}

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, *dur, fmt.Sprintf("%s_spike", *name)) {
		metrics.Add(res)
	}
	metrics.Close()

	rpt := vegeta.NewTextReporter(&metrics)
	rpt.Report(out)

	wg.Done()
}
