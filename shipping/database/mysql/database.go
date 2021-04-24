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
	res := r.db.Create(c)
	if err := res.Error; err != nil {
		return fmt.Errorf("failed to create cargo: %w", err)
	}

	return nil
}

func (r *cargoRepository) Find(id cargo.TrackingID) (*cargo.Cargo, error) {
	var cg cargo.Cargo

	res := r.db.Find(&cg, id)
	if err := res.Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve cargo: %w", err)
	}
	return &cg, nil
}

func (r *cargoRepository) FindAll() []*cargo.Cargo {
	var cgs []*cargo.Cargo

	res := r.db.Find(&cgs)
	if err := res.Error; err != nil {
		panic(fmt.Errorf("failed to retrieve all cargos: %w", err))
	}

	return cgs
}

func NewLocationRepository(db *gorm.DB) location.Repository {
	if err := createTableIfNotExist(db.Migrator(), &location.Location{}); err != nil {
		panic(fmt.Errorf("failed to create location table: %w", err))
	}

	res := db.Create(&defaultLocations)
	if err := res.Error; err != nil {
		panic(fmt.Errorf("failed to create default locations: %w", err))
	}

	return &locationRepository{
		db: db,
	}
}

func (r *locationRepository) Find(locode location.UNLocode) (*location.Location, error) {
	var loc location.Location

	res := r.db.Find(&loc, locode)
	if err := res.Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve location: %w", err)
	}
	return &loc, nil
}

func (r *locationRepository) FindAll() []*location.Location {
	var locs []*location.Location

	res := r.db.Find(&locs)
	if err := res.Error; err != nil {
		panic(fmt.Errorf("failed to retrieve all locations: %w", err))
	}

	return locs
}

func NewVoyageRepository(db *gorm.DB) voyage.Repository {
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

	res := db.Create(&defaultVoyages)
	if err := res.Error; err != nil {
		panic(fmt.Errorf("failed to create default voyages: %w", err))
	}

	return &voyageRepository{
		db: db,
	}
}

func (r *voyageRepository) Find(number voyage.Number) (*voyage.Voyage, error) {
	var voy voyage.Voyage

	res := r.db.Find(&voy, number)
	if err := res.Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve voyage: %w", err)
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
	var evts []cargo.HandlingEvent

	res := r.db.Find(&evts, id)
	if err := res.Error; err != nil {
		panic(fmt.Errorf("failed to retrieve handling events: %w", err))
	}

	return cargo.HandlingHistory{HandlingEvents: evts}
}
