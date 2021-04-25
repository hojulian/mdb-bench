package shipping

import (
	"fmt"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Auditor struct {
	url     string
	targets []vegeta.Target
}

func (a *Auditor) CheckCargo() {
	id := randTrackingID()
	u := fmt.Sprintf("%s/booking/v1/cargos/%s", a.url, id)
	t := vegeta.Target{
		Method: "GET",
		URL:    u,
	}

	a.targets = append(a.targets, t)
}

func (a *Auditor) CheckAllCargos() {
	u := fmt.Sprintf("%s/booking/v1/cargos", a.url)
	t := vegeta.Target{
		Method: "GET",
		URL:    u,
	}

	a.targets = append(a.targets, t)
}

func (a *Auditor) GetAllLocations() {
	u := fmt.Sprintf("%s/booking/v1/locations", a.url)
	t := vegeta.Target{
		Method: "GET",
		URL:    u,
	}

	a.targets = append(a.targets, t)
}

func (a *Auditor) Interactions() []vegeta.Target {
	return a.targets
}
