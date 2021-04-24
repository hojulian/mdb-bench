package microdb

import (
	"context"
	"fmt"

	"github.com/hojulian/microdb/client"
	"github.com/hojulian/microdb/microdb"

	"github.com/hojulian/mdb-bench/shipping/cargo"
)

func LoadDataOrigins(name string) error {
	err := microdb.AddDataOriginFromCfg(name)
	if err != nil {
		return fmt.Errorf("failed to create data origins from config: %w", err)
	}
	return nil
}

type cargoRepository struct {
	c *client.Client
}

type locationRepository struct {
	c *client.Client
}

type voyageRepository struct {
	c *client.Client
}

type handlingEventRepository struct {
	c *client.Client
}

func NewCargoRepository(c *client.Client) cargo.Repository {
	return &cargoRepository{
		c: c,
	}
}

func (r *cargoRepository) Store(c *cargo.Cargo) error {
	ctx := context.Background()
	q := "INSERT INTO `cargos` (`tracking_id`, `origin`, `route_specification_id`, `itinerary_id`, `delivery_id`) VALUES (?, ?, ?, ?, ?);"

	_, err := r.c.Execute(ctx, q, c.TrackingID, c.Origin)
	if err != nil {
		return fmt.Errorf("failed to create cargo: %w", err)
	}

	return nil
}

func (r *cargoRepository) Find(id cargo.TrackingID) (*cargo.Cargo, error) {
	ctx := context.Background()
	q := "SELECT `tracking_id`, `origin`, `route_specification_id`, `itinerary_id`, `delivery_id` FROM `cargos` WHERE `tracking_id` = ? LIMIT 1;"

	res, err := r.c.Query(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cargo: %w", err)
	}

	var cg cargo.Cargo
	for res.Next() {
		err = res.Scan(&cg.TrackingID, &cg.Origin, &cg.RouteSpecificationID, &cg.ItineraryID, &cg.DeliveryID)
		if err != nil {
			return nil, err
		}
	}

	return &cg, nil
}

func (r *cargoRepository) FindAll() []*cargo.Cargo {
	ctx := context.Background()
	q := "SELECT `tracking_id`, `origin`, `route_specification_id`, `itinerary_id`, `delivery_id` FROM `cargos`;"

	res, err := r.c.Query(ctx, q)
	if err != nil {
		panic(fmt.Errorf("failed to retrieve all cargos: %w", err))
	}

	var cgs []*cargo.Cargo
	for res.Next() {
		var cg cargo.Cargo

		err = res.Scan(&cg.TrackingID, &cg.Origin, &cg.RouteSpecificationID, &cg.ItineraryID, &cg.DeliveryID)
		if err != nil {
			panic(fmt.Errorf("failed to retrieve cargo: %w", err))
		}

		cgs = append(cgs, &cg)
	}

	return cgs
}
