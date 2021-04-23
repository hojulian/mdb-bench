package voyage

import "github.com/hojulian/mdb-bench/shipping/location"

// A set of sample voyages.
var (
	V100 = New("V100", Schedule{
		CarrierMovements: []CarrierMovement{
			{DepartureLocation: location.CNHKG, ArrivalLocation: location.JNTKO},
			{DepartureLocation: location.JNTKO, ArrivalLocation: location.USNYC},
		},
	})

	V300 = New("V300", Schedule{
		CarrierMovements: []CarrierMovement{
			{DepartureLocation: location.JNTKO, ArrivalLocation: location.NLRTM},
			{DepartureLocation: location.NLRTM, ArrivalLocation: location.DEHAM},
			{DepartureLocation: location.DEHAM, ArrivalLocation: location.AUMEL},
			{DepartureLocation: location.AUMEL, ArrivalLocation: location.JNTKO},
		},
	})

	V400 = New("V400", Schedule{
		CarrierMovements: []CarrierMovement{
			{DepartureLocation: location.DEHAM, ArrivalLocation: location.SESTO},
			{DepartureLocation: location.SESTO, ArrivalLocation: location.FIHEL},
			{DepartureLocation: location.FIHEL, ArrivalLocation: location.DEHAM},
		},
	})
)

// These voyages are hard-coded into the current pathfinder. Make sure
// they exist.
var (
	V0100S = New("0100S", Schedule{CarrierMovements: []CarrierMovement{}})
	V0200T = New("0200T", Schedule{CarrierMovements: []CarrierMovement{}})
	V0300A = New("0300A", Schedule{CarrierMovements: []CarrierMovement{}})
	V0301S = New("0301S", Schedule{CarrierMovements: []CarrierMovement{}})
	V0400S = New("0400S", Schedule{CarrierMovements: []CarrierMovement{}})
)
