package database

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/hojulian/microdb/client"
	mysqld "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/hojulian/mdb-bench/shipping/cargo"
	"github.com/hojulian/mdb-bench/shipping/database/inmem"
	mysqldb "github.com/hojulian/mdb-bench/shipping/database/mysql"
	"github.com/hojulian/mdb-bench/shipping/location"
	"github.com/hojulian/mdb-bench/shipping/voyage"
)

type DatabaseType string

var (
	DatabaseTypeInMem   DatabaseType = "inmem"
	DatabaseTypeMySQL   DatabaseType = "mysql"
	DatabaseTypeMicroDB DatabaseType = "microdb"
)

func NewCargoRepository(databaseType DatabaseType, params map[string]string) (cargo.Repository, error) {
	switch databaseType {
	case DatabaseTypeInMem:
		return inmem.NewCargoRepository(), nil
	case DatabaseTypeMySQL:
		db, err := mySQLDB(
			params["MYSQL_HOST"],
			params["MYSQL_PORT"],
			params["MYSQL_USER"],
			params["MYSQL_PASSWORD"],
			params["MYSQL_DATABASE"],
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return mysqldb.NewCargoRepository(db), nil
	case DatabaseTypeMicroDB:
		_, err := microDB(
			params["NATS_HOST"],
			params["NATS_PORT"],
			params["NATS_CLIENT_ID"],
			params["NATS_CLUSTER_ID"],
			"handling_events",
			"handling_activities",
			"route_specifications",
			"itineraries",
			"deliveries",
			"legs",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connec to database: %w", err)
		}
	}

	return nil, fmt.Errorf("unsupported database type")
}

func NewLocationRepository(databaseType DatabaseType, params map[string]string) (location.Repository, error) {
	switch databaseType {
	case DatabaseTypeInMem:
		return inmem.NewLocationRepository(), nil
	case DatabaseTypeMySQL:
		db, err := mySQLDB(
			params["MYSQL_HOST"],
			params["MYSQL_PORT"],
			params["MYSQL_USER"],
			params["MYSQL_PASSWORD"],
			params["MYSQL_DATABASE"],
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return mysqldb.NewLocationRepository(db), nil
	case DatabaseTypeMicroDB:
		_, err := microDB(
			params["NATS_HOST"],
			params["NATS_PORT"],
			params["NATS_CLIENT_ID"],
			params["NATS_CLUSTER_ID"],
			"locations",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connec to database: %w", err)
		}
	}

	return nil, fmt.Errorf("unsupported database type")
}

func NewVoyageRepository(databaseType DatabaseType, params map[string]string) (voyage.Repository, error) {
	switch databaseType {
	case DatabaseTypeInMem:
		return inmem.NewVoyageRepository(), nil
	case DatabaseTypeMySQL:
		db, err := mySQLDB(
			params["MYSQL_HOST"],
			params["MYSQL_PORT"],
			params["MYSQL_USER"],
			params["MYSQL_PASSWORD"],
			params["MYSQL_DATABASE"],
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return mysqldb.NewVoyageRepository(db), nil
	case DatabaseTypeMicroDB:
		_, err := microDB(
			params["NATS_HOST"],
			params["NATS_PORT"],
			params["NATS_CLIENT_ID"],
			params["NATS_CLUSTER_ID"],
			"voyages",
			"carrier_movements",
			"schedules",
			"locations",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connec to database: %w", err)
		}
	}

	return nil, fmt.Errorf("unsupported database type")
}

func NewHandlingEventRepository(databaseType DatabaseType, params map[string]string) (cargo.HandlingEventRepository, error) {
	switch databaseType {
	case DatabaseTypeInMem:
		return inmem.NewHandlingEventRepository(), nil
	case DatabaseTypeMySQL:
		db, err := mySQLDB(
			params["MYSQL_HOST"],
			params["MYSQL_PORT"],
			params["MYSQL_USER"],
			params["MYSQL_PASSWORD"],
			params["MYSQL_DATABASE"],
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return mysqldb.NewHandlingEventRepository(db), nil
	case DatabaseTypeMicroDB:
		_, err := microDB(
			params["NATS_HOST"],
			params["NATS_PORT"],
			params["NATS_CLIENT_ID"],
			params["NATS_CLUSTER_ID"],
			"voyages",
			"carrier_movements",
			"schedules",
			"locations",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connec to database: %w", err)
		}
	}

	return nil, fmt.Errorf("unsupported database type")
}

func mySQLDB(host, port, user, password, database string) (*gorm.DB, error) {
	dsn := mySQLConnectionCfg(host, port, user, password, database)
	db, err := gorm.Open(mysqld.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db = db.Set("gorm:table_options", "DEFAULT CHARSET=utf8")
	db = db.Clauses(clause.OnConflict{DoNothing: true})
	return db, nil
}

func mySQLConnectionCfg(host, port, user, password, database string) string {
	mCfg := mysql.NewConfig()
	mCfg.Net = "tcp"
	mCfg.Addr = fmt.Sprintf("%s:%s", host, port)
	mCfg.User = user
	mCfg.Passwd = password
	mCfg.DBName = database
	mCfg.ParseTime = true

	return mCfg.FormatDSN()
}

func microDB(natsHost, natsPort, natsClientID, natsClusterID string, tables ...string) (*client.Client, error) {
	c, err := client.Connect(natsHost, natsPort, natsClientID, natsClusterID, tables...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to microdb: %w", err)
	}
	return c, nil
}
