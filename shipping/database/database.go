package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hojulian/microdb/client"
	mysqld "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"

	"github.com/hojulian/mdb-bench/shipping/cargo"
	"github.com/hojulian/mdb-bench/shipping/database/inmem"
	"github.com/hojulian/mdb-bench/shipping/database/microdb"
	mysqldb "github.com/hojulian/mdb-bench/shipping/database/mysql"
	"github.com/hojulian/mdb-bench/shipping/location"
	"github.com/hojulian/mdb-bench/shipping/voyage"
)

type DatabaseType string

var (
	CargoTables = []string{
		"cargos",
		"handling_histories",
		"handling_events",
		"handling_activities",
		"route_specifications",
		"itineraries",
		"deliveries",
		"legs",
	}
	LocationTables = []string{
		"locations",
	}
	VoyageTables = []string{
		"voyages",
		"carrier_movements",
		"schedules",
		"locations",
	}
	HandlingEventTables = []string{
		"voyages",
		"carrier_movements",
		"schedules",
		"locations",
	}

	DatabaseTypeInMem        DatabaseType = "inmem"
	DatabaseTypeMySQL        DatabaseType = "mysql"
	DatabaseTypeMySQLCluster DatabaseType = "mysql-cluster"
	DatabaseTypeMicroDB      DatabaseType = "microdb"

	defaultGormConfig *gorm.Config = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
)

func NewInMemCargoRepository() cargo.Repository {
	return inmem.NewCargoRepository()
}

func NewMySQLCargoRepository(params map[string]string, cluster bool) (cargo.Repository, error) {
	db, err := mysqlGorm(params, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return mysqldb.NewCargoRepository(db), nil
}

// Required tables:
// "cargos",
// "handling_histories",
// "handling_events",
// "handling_activities",
// "route_specifications",
// "itineraries",
// "deliveries",
// "legs",
func NewMicroDBCargoRepository(c *client.Client) cargo.Repository {
	return microdb.NewCargoRepository(c)
}

func NewInMemLocationRepository() location.Repository {
	return inmem.NewLocationRepository()
}

func NewMySQLLocationRepository(params map[string]string, cluster bool) (location.Repository, error) {
	withDefaults, err := strconv.ParseBool(params["DB_DEFAULTS"])
	if err != nil {
		return nil, fmt.Errorf("invalid value for DB_DEFAULTS: %w, must be bool", err)
	}

	db, err := mysqlGorm(params, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return mysqldb.NewLocationRepository(db, withDefaults), nil
}

// Required tables:
// "locations",
func NewMicroDBLocationRepository(c *client.Client) location.Repository {
	return microdb.NewLocationRepository(c)
}

func NewInMemVoyageRepository() voyage.Repository {
	return inmem.NewVoyageRepository()
}

func NewMySQLVoyageRepository(params map[string]string, cluster bool) (voyage.Repository, error) {
	withDefaults, err := strconv.ParseBool(params["DB_DEFAULTS"])
	if err != nil {
		return nil, fmt.Errorf("invalid value for DB_DEFAULTS: %w, must be bool", err)
	}

	db, err := mysqlGorm(params, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return mysqldb.NewVoyageRepository(db, withDefaults), nil
}

// Required tables:
// "voyages",
// "carrier_movements",
// "schedules",
// "locations",
func NewMicroDBVoyageRepository(c *client.Client) voyage.Repository {
	return microdb.NewVoyageRepository(c)
}

func NewInMemHandlingEventRepository() cargo.HandlingEventRepository {
	return inmem.NewHandlingEventRepository()
}

func NewMySQLHandlingEventRepository(params map[string]string, cluster bool) (cargo.HandlingEventRepository, error) {
	db, err := mysqlGorm(params, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return mysqldb.NewHandlingEventRepository(db), nil
}

// Required tables:
// "voyages",
// "carrier_movements",
// "schedules",
// "locations",
func NewMicroDBHandlingEventRepository(c *client.Client) cargo.HandlingEventRepository {
	return microdb.NewHandlingEventRepository(c)
}

func mysqlGorm(params map[string]string, cluster bool) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if cluster {
		nodes, err := strconv.ParseInt(params["MYSQL_NODES"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number of nodes: %w", err)
		}

		dsns := make([]string, 0, int(nodes))
		for i := 0; i < int(nodes); i++ {
			dsn := mysqlConnectionCfg(
				params[fmt.Sprintf("MYSQL_HOST_%d", i)],
				params[fmt.Sprintf("MYSQL_PORT_%d", i)],
				params["MYSQL_USER"],
				params["MYSQL_PASSWORD"],
				params["MYSQL_DATABASE"],
			)
			dsns = append(dsns, dsn)
		}

		db, err = mysqlDBCluster(dsns...)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database cluster: %w", err)
		}
	} else {
		db, err = mysqlDB(
			params["MYSQL_HOST"],
			params["MYSQL_PORT"],
			params["MYSQL_USER"],
			params["MYSQL_PASSWORD"],
			params["MYSQL_DATABASE"],
		)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
	}

	return db, nil
}

func mysqlDB(host, port, user, password, database string) (*gorm.DB, error) {
	dsn := mysqlConnectionCfg(host, port, user, password, database)
	db, err := gorm.Open(mysqld.Open(dsn), defaultGormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db = configureGormDB(db)
	return db, nil
}

func mysqlConnectionCfg(host, port, user, password, database string) string {
	mCfg := mysql.NewConfig()
	mCfg.Net = "tcp"
	mCfg.Addr = fmt.Sprintf("%s:%s", host, port)
	mCfg.User = user
	mCfg.Passwd = password
	mCfg.DBName = database
	mCfg.ParseTime = true
	dsn := mCfg.FormatDSN()

	return dsn
}

func configureGormDB(db *gorm.DB) *gorm.DB {
	db = db.Set("gorm:table_options", "DEFAULT CHARSET=utf8")
	db = db.Clauses(clause.OnConflict{DoNothing: true})
	return db
}

func mysqlDBCluster(dsn ...string) (*gorm.DB, error) {
	if len(dsn) == 0 {
		return nil, fmt.Errorf("require at least 1 node in cluster")
	}

	dbs := make([]gorm.Dialector, 0, len(dsn))
	for _, d := range dsn {
		db := mysqld.Open(d)
		dbs = append(dbs, db)
	}

	db, err := gorm.Open(dbs[0], defaultGormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{dbs[0]},
			Replicas: dbs[1:],
			Policy:   dbresolver.RandomPolicy{},
		}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)
	db = configureGormDB(db)
	return db, nil
}
