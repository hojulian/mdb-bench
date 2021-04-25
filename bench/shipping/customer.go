package shipping

import (
	"fmt"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Customer struct {
	url     string
	targets []vegeta.Target
}

func (c *Customer) TrackCargo() {
	id := randTrackingID()
	u := fmt.Sprintf("%s/tracking/v1/cargos/%s", c.url, string(id))
	t := vegeta.Target{
		Method: "GET",
		URL:    u,
	}

	c.targets = append(c.targets, t)
}

func (c *Customer) Interactions() []vegeta.Target {
	return c.targets
}
