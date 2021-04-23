package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/hojulian/mdb-bench/shipping/cargo"
	"github.com/hojulian/mdb-bench/shipping/location"
	"github.com/hojulian/mdb-bench/shipping/voyage"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:test@/test?charset=utf8&parseTime=True"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	db = db.Set("gorm:table_options", "DEFAULT CHARSET=utf8")
	migrator := db.Migrator()

	tables := []interface{}{
		voyage.Voyage{},
		voyage.CarrierMovement{},
		voyage.Schedule{},
		location.Location{},
		cargo.HandlingEvent{},
		cargo.HandlingActivity{},
		cargo.RouteSpecification{},
		cargo.Itinerary{},
		cargo.Delivery{},
		cargo.Leg{},
	}

	for _, t := range tables {
		if err := createIfNotExist(migrator, t); err != nil {
			log.Fatalf("failed to create table: %s", err)
		}
		log.Printf("created table")
	}
}

func createIfNotExist(migrator gorm.Migrator, table interface{}) error {
	if migrator.HasTable(table) {
		return nil
	}

	return migrator.CreateTable(table)
}
