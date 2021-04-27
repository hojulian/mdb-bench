package shipping

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type User interface {
	Interactions() []vegeta.Target
}

func RegularCustomer(url string) User {
	c := &Customer{
		url: url,
	}

	c.TrackCargo()
	c.TrackCargo()

	return c
}

func RegularAuditor(url string) User {
	a := &Auditor{
		url: url,
	}

	a.GetAllLocations()
	a.CheckCargo()
	a.CheckCargo()

	return a
}

func RegularBooker(url string) User {
	b := &Booker{
		url: url,
	}

	b.BookCargo()
	b.ViewCargo()
	b.BookCargo()
	b.ViewCargo()
	b.BookCargo()
	b.ViewCargo()

	return b
}

func HighLoadBooker(url string) User {
	b := &Booker{
		url: url,
	}

	for i := 0; i < 500; i++ {
		b.BookCargo()
		b.ViewCargo()
	}

	return b
}
