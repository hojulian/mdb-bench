package shipping

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"

	"github.com/hojulian/mdb-bench/shipping/cargo"
	"github.com/hojulian/mdb-bench/shipping/location"
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

type Booker struct {
	url     string
	targets []vegeta.Target
}

type bookCargoRequest struct {
	Origin          location.UNLocode `json:"origin"`
	Destination     location.UNLocode `json:"destination"`
	ArrivalDeadline string            `json:"arrival_deadline"`
}

func (b *Booker) ViewCargo() {
	id := cargo.NextTrackingID()
	u := fmt.Sprintf("%s/booking/v1/cargos/%s", b.url, id)
	t := vegeta.Target{
		Method: "GET",
		URL:    u,
	}

	b.targets = append(b.targets, t)
}

func (b *Booker) ViewAllCargos() {
	u := fmt.Sprintf("%s/booking/v1/cargos", b.url)
	t := vegeta.Target{
		Method: "GET",
		URL:    u,
	}

	b.targets = append(b.targets, t)
}

func (b *Booker) BookCargo() error {
	u := fmt.Sprintf("%s/booking/v1/cargos", b.url)

	rs := &bookCargoRequest{
		Origin:          randomLoc(),
		Destination:     randomLoc(),
		ArrivalDeadline: time.Now().AddDate(0, 0, rand.Intn(30)).Format(time.RFC3339),
	}

	body, err := json.Marshal(rs)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	t := vegeta.Target{
		Method: "POST",
		URL:    u,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: []byte(body),
	}

	b.targets = append(b.targets, t)
	return nil
}

func randomLoc() location.UNLocode {
	return locations[rand.Intn(len(locations))]
}

func (b *Booker) AssignCargoToRoute(id cargo.TrackingID, itinerary cargo.Itinerary) error {
	u := fmt.Sprintf("%s/booking/v1/cargos/%s/assign_to_route", b.url, id)

	body, err := json.Marshal(itinerary)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	t := vegeta.Target{
		Method: "POST",
		URL:    u,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: []byte(body),
	}

	b.targets = append(b.targets, t)
	return nil
}

func (b *Booker) Interactions() []vegeta.Target {
	return b.targets
}
