// Package voyage provides the Voyage aggregate.
package voyage

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/hojulian/mdb-bench/shipping/location"
)

// Number uniquely identifies a particular Voyage.
type Number string

// Voyage is a uniquely identifiable series of carrier movements.
type Voyage struct {
	Number     Number `gorm:"primaryKey"`
	ScheduleID int
	Schedule   Schedule `gorm:"foreignKey:ScheduleID"`
}

// New creates a voyage with a voyage number and a provided schedule.
func New(n Number, s Schedule) *Voyage {
	return &Voyage{Number: n, Schedule: s}
}

// Schedule describes a voyage schedule.
type Schedule struct {
	gorm.Model
	CarrierMovements []CarrierMovement `gorm:"foreignKey:ScheduleRefer"`
}

// CarrierMovement is a vessel voyage from one location to another.
type CarrierMovement struct {
	gorm.Model
	DepartureLocation location.UNLocode
	ArrivalLocation   location.UNLocode
	DepartureTime     time.Time `gorm:"default:null"`
	ArrivalTime       time.Time `gorm:"default:null"`
	ScheduleRefer     uint
}

// ErrUnknown is used when a voyage could not be found.
var ErrUnknown = errors.New("unknown voyage")

// Repository provides access a voyage store.
type Repository interface {
	Find(Number) (*Voyage, error)
}
