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
	a.CheckAllCargos()
	a.CheckCargo()
	a.CheckAllCargos()
	a.CheckCargo()

	return a
}

func RegularBooker(url string) User {
	b := &Booker{
		url: url,
	}

	b.ViewAllCargos()
	b.BookCargo()
	b.ViewCargo()
	b.ViewAllCargos()
	b.BookCargo()
	b.ViewCargo()
	b.ViewAllCargos()
	b.BookCargo()
	b.ViewAllCargos()
	b.ViewCargo()

	return b
}

func BookOnlyBooker(url string) User {
	b := &Booker{
		url: url,
	}

	b.BookCargo()
	b.BookCargo()
	b.BookCargo()
	b.BookCargo()

	return b
}
