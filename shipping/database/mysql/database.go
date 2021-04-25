package mysql

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/hojulian/mdb-bench/shipping/cargo"
	"github.com/hojulian/mdb-bench/shipping/location"
	"github.com/hojulian/mdb-bench/shipping/voyage"
)

var (
	defaultLocations = []location.Location{
		*location.Stockholm,
		*location.Melbourne,
		*location.Hongkong,
		*location.Tokyo,
		*location.Rotterdam,
		*location.Hamburg,
		*location.Chicago,
		*location.Helsinki,
		*location.NewYork,
	}

	defaultVoyages = []voyage.Voyage{
		*voyage.V100,
		*voyage.V300,
		*voyage.V400,
		*voyage.V0100S,
		*voyage.V0200T,
		*voyage.V0300A,
		*voyage.V0301S,
		*voyage.V0400S,
	}
)

type cargoRepository struct {
	db *gorm.DB
}

type locationRepository struct {
	db *gorm.DB
}

type voyageRepository struct {
	db *gorm.DB
}

type handlingEventRepository struct {
	db *gorm.DB
}

func createTableIfNotExist(migrator gorm.Migrator, table interface{}) error {
	if migrator.HasTable(table) {
		return nil
	}

	return migrator.CreateTable(table)
}

func NewCargoRepository(db *gorm.DB) cargo.Repository {
	requiredTables := []interface{}{
		cargo.Cargo{},
		cargo.HandlingHistory{},
		cargo.HandlingEvent{},
		cargo.HandlingActivity{},
		cargo.RouteSpecification{},
		cargo.Itinerary{},
		cargo.Delivery{},
		cargo.Leg{},
	}

	for _, t := range requiredTables {
		if err := createTableIfNotExist(db.Migrator(), t); err != nil {
			panic(fmt.Errorf("failed to create cargo tables: %w", err))
		}
	}

	return &cargoRepository{
		db: db,
	}
}

func (r *cargoRepository) Store(c *cargo.Cargo) error {
	db, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve db instance: %w", err)
	}

	_, err = db.Exec(sqlInsertCargo, c.TrackingID, c.Origin, c.RouteSpecificationID, c.ItineraryID, c.DeliveryID)
	if err != nil {
		return fmt.Errorf("failed to insert cargo: %w", err)
	}

	return nil
}

func (r *cargoRepository) Find(id cargo.TrackingID) (*cargo.Cargo, error) {
	var cg cargo.Cargo
	// s := r.db.Session(&gorm.Session{DryRun: true}).Joins("RouteSpecification").Joins("Itinerary").Joins("Delivery").Find(&cg, "tracking_id", id).Statement
	// q := s.SQL.String()

	db, err := r.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve db instance: %w", err)
	}

	rows, err := db.Query(sqlFindCargoByID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := r.db.ScanRows(rows, &cg); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	return &cg, nil
}

func (r *cargoRepository) FindAll() []*cargo.Cargo {
	var cgs []*cargo.Cargo

	db, err := r.db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to retrieve db instance: %w", err))
	}

	rows, err := db.Query(sqlFindAllCargos)
	if err != nil {
		panic(fmt.Errorf("failed to read rows: %w", err))
	}
	if err := rows.Err(); err != nil {
		panic(fmt.Errorf("failed to read rows: %w", err))
	}
	defer rows.Close()

	for rows.Next() {
		var cg cargo.Cargo
		if err := r.db.ScanRows(rows, &cg); err != nil {
			panic(fmt.Errorf("failed to scan row: %w", err))
		}
		cgs = append(cgs, &cg)
	}

	return cgs
}

func NewLocationRepository(db *gorm.DB, withDefaults bool) location.Repository {
	if err := createTableIfNotExist(db.Migrator(), &location.Location{}); err != nil {
		panic(fmt.Errorf("failed to create location table: %w", err))
	}

	if withDefaults {
		res := db.Create(&defaultLocations)
		if err := res.Error; err != nil {
			panic(fmt.Errorf("failed to create default locations: %w", err))
		}
	}

	return &locationRepository{
		db: db,
	}
}

func (r *locationRepository) Find(locode location.UNLocode) (*location.Location, error) {
	var loc location.Location

	db, err := r.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve db instance: %w", err)
	}

	rows, err := db.Query(sqlFindLocationByID, locode)
	if err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := r.db.ScanRows(rows, &loc); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	return &loc, nil
}

func (r *locationRepository) FindAll() []*location.Location {
	var locs []*location.Location

	db, err := r.db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to retrieve db instance: %w", err))
	}

	rows, err := db.Query(sqlFindAllLocations)
	if err != nil {
		panic(fmt.Errorf("failed to read rows: %w", err))
	}
	if err := rows.Err(); err != nil {
		panic(fmt.Errorf("failed to read rows: %w", err))
	}
	defer rows.Close()

	for rows.Next() {
		var loc location.Location
		if err := r.db.ScanRows(rows, &loc); err != nil {
			panic(fmt.Errorf("failed to scan row: %w", err))
		}
		locs = append(locs, &loc)
	}

	return locs
}

func NewVoyageRepository(db *gorm.DB, withDefaults bool) voyage.Repository {
	requiredTables := []interface{}{
		voyage.Voyage{},
		voyage.CarrierMovement{},
		voyage.Schedule{},
		location.Location{},
	}

	for _, t := range requiredTables {
		if err := createTableIfNotExist(db.Migrator(), t); err != nil {
			panic(fmt.Errorf("failed to create voyage tables: %w", err))
		}
	}

	if withDefaults {
		res := db.Create(&defaultVoyages)
		if err := res.Error; err != nil {
			panic(fmt.Errorf("failed to create default voyages: %w", err))
		}
	}

	return &voyageRepository{
		db: db,
	}
}

func (r *voyageRepository) Find(number voyage.Number) (*voyage.Voyage, error) {
	var voy voyage.Voyage

	db, err := r.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve db instance: %w", err)
	}

	rows, err := db.Query(sqlFindVoyageByID, number)
	if err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := r.db.ScanRows(rows, &voy); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	return &voy, nil
}

func NewHandlingEventRepository(db *gorm.DB) cargo.HandlingEventRepository {
	return &handlingEventRepository{
		db: db,
	}
}

func (r *handlingEventRepository) Store(e cargo.HandlingEvent) {
	res := r.db.Create(&e)
	if err := res.Error; err != nil {
		panic(fmt.Errorf("failed to create handling event: %w", err))
	}
}

func (r *handlingEventRepository) QueryHandlingHistory(id cargo.TrackingID) cargo.HandlingHistory {
	// var evt cargo.HandlingEvent
	// res := r.db.Find(&evt, "tracking_id = ?", id)
	// if err := res.Error; err != nil {
	// 	fmt.Println(fmt.Errorf("failed to retrieve handling even: %w", err))
	// }

	var his cargo.HandlingHistory

	return his
}
