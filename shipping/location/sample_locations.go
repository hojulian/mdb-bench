package location

// Sample UN locodes.
var (
	SESTO UNLocode = "SESTO"
	AUMEL UNLocode = "AUMEL"
	CNHKG UNLocode = "CNHKG"
	USNYC UNLocode = "USNYC"
	USCHI UNLocode = "USCHI"
	JNTKO UNLocode = "JNTKO"
	DEHAM UNLocode = "DEHAM"
	NLRTM UNLocode = "NLRTM"
	FIHEL UNLocode = "FIHEL"
)

// Sample locations.
var (
	Stockholm = &Location{UNLocode: SESTO, Name: "Stockholm"}
	Melbourne = &Location{UNLocode: AUMEL, Name: "Melbourne"}
	Hongkong  = &Location{UNLocode: CNHKG, Name: "Hongkong"}
	NewYork   = &Location{UNLocode: USNYC, Name: "New York"}
	Chicago   = &Location{UNLocode: USCHI, Name: "Chicago"}
	Tokyo     = &Location{UNLocode: JNTKO, Name: "Tokyo"}
	Hamburg   = &Location{UNLocode: DEHAM, Name: "Hamburg"}
	Rotterdam = &Location{UNLocode: NLRTM, Name: "Rotterdam"}
	Helsinki  = &Location{UNLocode: FIHEL, Name: "Helsinki"}
)
