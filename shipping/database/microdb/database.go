package microdb

import (
	"context"
	"fmt"

	"github.com/hojulian/microdb/client"
	"github.com/hojulian/microdb/microdb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/hojulian/mdb-bench/shipping/cargo"
	"github.com/hojulian/mdb-bench/shipping/location"
	"github.com/hojulian/mdb-bench/shipping/voyage"
)

func LoadDataOrigins(name string) error {
	err := microdb.AddDataOriginFromCfg(name)
	if err != nil {
		return fmt.Errorf("failed to create data origins from config: %w", err)
	}
	return nil
}

type cargoRepository struct {
	g  *gorm.DB
	do *microdb.DataOrigin
	c  *client.Client
}

type locationRepository struct {
	g *gorm.DB
	c *client.Client
}

type voyageRepository struct {
	g *gorm.DB
	c *client.Client
}

type handlingEventRepository struct {
	g *gorm.DB
	c *client.Client
}

func NewCargoRepository(c *client.Client) cargo.Repository {
	// Sqlite here is not actually used. This is only for getting gorm to play well.
	g, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to initialize gorm: %w", err))
	}

	do, err := microdb.GetDataOrigin("cargos")
	if err != nil {
		panic(fmt.Errorf("failed to get data origin for cargos: %w", err))
	}

	return &cargoRepository{
		c:  c,
		g:  g,
		do: do,
	}
}

func (r *cargoRepository) Store(c *cargo.Cargo) error {
	ctx := context.Background()

	_, err := r.c.Execute(
		ctx,
		sqlInsertCargo,
		string(c.TrackingID),
		string(c.Origin),
		c.RouteSpecificationID,
		c.ItineraryID,
		c.DeliveryID,
	)

	if err != nil {
		return fmt.Errorf("failed to insert cargo: %w", err)
	}

	return nil
}

func (r *cargoRepository) Find(id cargo.TrackingID) (*cargo.Cargo, error) {
	ctx := context.Background()

	rows, err := r.c.Query(ctx, sqlFindCargoByID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	defer rows.Close()

	var cg cargo.Cargo
	if rows.Next() {
		if err := r.g.Model(&cg).ScanRows(rows, &cg); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	return &cg, nil
}

func (r *cargoRepository) FindAll() []*cargo.Cargo {
	ctx := context.Background()

	rows, err := r.c.Query(ctx, sqlFindAllCargos)
	if err != nil {
		panic(fmt.Errorf("failed to retrieve all cargos: %w", err))
	}
	defer rows.Close()

	var cgs []*cargo.Cargo
	for rows.Next() {
		var cg cargo.Cargo

		if err := r.g.Model(&cg).ScanRows(rows, &cg); err != nil {
			panic(fmt.Errorf("failed to scan row: %w", err))
		}

		cgs = append(cgs, &cg)
	}

	return cgs
}

func NewLocationRepository(c *client.Client) location.Repository {
	// Sqlite here is not actually used. This is only for getting gorm to play well.
	g, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to initialize gorm: %w", err))
	}

	return &locationRepository{
		c: c,
		g: g,
	}
}

func (r *locationRepository) Find(locode location.UNLocode) (*location.Location, error) {
	ctx := context.Background()

	rows, err := r.c.Query(ctx, sqlFindLocationByID, locode)
	if err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	defer rows.Close()

	var loc location.Location
	if rows.Next() {
		if err := r.g.Model(&loc).ScanRows(rows, &loc); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	return &loc, nil
}

func (r *locationRepository) FindAll() []*location.Location {
	ctx := context.Background()

	rows, err := r.c.Query(ctx, sqlFindAllLocations)
	if err != nil {
		panic(fmt.Errorf("failed to read rows: %w", err))
	}
	if err := rows.Err(); err != nil {
		panic(fmt.Errorf("failed to read rows: %w", err))
	}
	defer rows.Close()

	var locs []*location.Location
	for rows.Next() {
		var loc location.Location

		if err := r.g.Model(&loc).ScanRows(rows, &loc); err != nil {
			panic(fmt.Errorf("failed to scan row: %w", err))
		}

		locs = append(locs, &loc)
	}

	return locs
}

func NewVoyageRepository(c *client.Client) voyage.Repository {
	// Sqlite here is not actually used. This is only for getting gorm to play well.
	g, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to initialize gorm: %w", err))
	}

	return &voyageRepository{
		c: c,
		g: g,
	}
}

func (r *voyageRepository) Find(number voyage.Number) (*voyage.Voyage, error) {
	ctx := context.Background()

	rows, err := r.c.Query(ctx, sqlFindVoyageByID, number)
	if err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read row: %w", err)
	}
	defer rows.Close()

	var voy voyage.Voyage
	if rows.Next() {
		if err := r.g.Model(&voy).ScanRows(rows, &voy); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
	}

	return &voy, nil
}

func NewHandlingEventRepository(c *client.Client) cargo.HandlingEventRepository {
	// Sqlite here is not actually used. This is only for getting gorm to play well.
	g, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to initialize gorm: %w", err))
	}

	return &handlingEventRepository{
		c: c,
		g: g,
	}
}

func (r *handlingEventRepository) Store(e cargo.HandlingEvent) {
	ctx := context.Background()

	_, err := r.c.Execute(ctx, sqlInsertHandlingEvent, e.TrackingID, e.ActivityID, e.HandlingHistoryRefer)
	if err != nil {
		panic(fmt.Errorf("failed to insert handling event: %w", err))
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
