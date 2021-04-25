package database

import (
	"fmt"
	"os"
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
	DatabaseTypeInMem        DatabaseType = "inmem"
	DatabaseTypeMySQL        DatabaseType = "mysql"
	DatabaseTypeMySQLCluster DatabaseType = "mysql-cluster"
	DatabaseTypeMicroDB      DatabaseType = "microdb"

	defaultGormConfig *gorm.Config = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
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

	case DatabaseTypeMySQLCluster:
		nodes, err := strconv.ParseInt(params["MYSQL_NODES"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number of nodes: %w", err)
		}

		dsns := make([]string, 0, int(nodes))
		for i := 0; i < int(nodes); i++ {
			dsn := mySQLConnectionCfg(
				os.Getenv(fmt.Sprintf("MYSQL_HOST_%d", i)),
				os.Getenv(fmt.Sprintf("MYSQL_PORT_%d", i)),
				params["MYSQL_USER"],
				params["MYSQL_PASSWORD"],
				params["MYSQL_DATABASE"],
			)
			dsns = append(dsns, dsn)
		}

		db, err := mySQLDBCluster(dsns...)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database cluster: %w", err)
		}
		return mysqldb.NewCargoRepository(db), nil

	case DatabaseTypeMicroDB:
		_, err := microDB(
			params["NATS_HOST"],
			params["NATS_PORT"],
			params["NATS_CLIENT_ID"],
			params["NATS_CLUSTER_ID"],
			"cargos",
			"handling_histories",
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
	withDefaults, _ := strconv.ParseBool(params["DB_DEFAULTS"])

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
		return mysqldb.NewLocationRepository(db, withDefaults), nil

	case DatabaseTypeMySQLCluster:
		nodes, err := strconv.ParseInt(params["MYSQL_NODES"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number of nodes: %w", err)
		}

		dsns := make([]string, 0, int(nodes))
		for i := 0; i < int(nodes); i++ {
			dsn := mySQLConnectionCfg(
				os.Getenv(fmt.Sprintf("MYSQL_HOST_%d", i)),
				os.Getenv(fmt.Sprintf("MYSQL_PORT_%d", i)),
				params["MYSQL_USER"],
				params["MYSQL_PASSWORD"],
				params["MYSQL_DATABASE"],
			)
			dsns = append(dsns, dsn)
		}

		db, err := mySQLDBCluster(dsns...)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database cluster: %w", err)
		}
		return mysqldb.NewLocationRepository(db, withDefaults), nil

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
	withDefaults, _ := strconv.ParseBool(params["DB_DEFAULTS"])

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
		return mysqldb.NewVoyageRepository(db, withDefaults), nil

	case DatabaseTypeMySQLCluster:
		nodes, err := strconv.ParseInt(params["MYSQL_NODES"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number of nodes: %w", err)
		}

		dsns := make([]string, 0, int(nodes))
		for i := 0; i < int(nodes); i++ {
			dsn := mySQLConnectionCfg(
				os.Getenv(fmt.Sprintf("MYSQL_HOST_%d", i)),
				os.Getenv(fmt.Sprintf("MYSQL_PORT_%d", i)),
				params["MYSQL_USER"],
				params["MYSQL_PASSWORD"],
				params["MYSQL_DATABASE"],
			)
			dsns = append(dsns, dsn)
		}

		db, err := mySQLDBCluster(dsns...)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database cluster: %w", err)
		}
		return mysqldb.NewVoyageRepository(db, withDefaults), nil

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

	case DatabaseTypeMySQLCluster:
		nodes, err := strconv.ParseInt(params["MYSQL_NODES"], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number of nodes: %w", err)
		}

		dsns := make([]string, 0, int(nodes))
		for i := 0; i < int(nodes); i++ {
			dsn := mySQLConnectionCfg(
				os.Getenv(fmt.Sprintf("MYSQL_HOST_%d", i)),
				os.Getenv(fmt.Sprintf("MYSQL_PORT_%d", i)),
				params["MYSQL_USER"],
				params["MYSQL_PASSWORD"],
				params["MYSQL_DATABASE"],
			)
			dsns = append(dsns, dsn)
		}

		db, err := mySQLDBCluster(dsns...)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database cluster: %w", err)
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
	db, err := gorm.Open(mysqld.Open(dsn), defaultGormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db = configureGormDB(db)
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
	dsn := mCfg.FormatDSN()

	return dsn
}

func configureGormDB(db *gorm.DB) *gorm.DB {
	db = db.Set("gorm:table_options", "DEFAULT CHARSET=utf8")
	db = db.Clauses(clause.OnConflict{DoNothing: true})
	return db
}

func mySQLDBCluster(dsn ...string) (*gorm.DB, error) {
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

func microDB(natsHost, natsPort, natsClientID, natsClusterID string, tables ...string) (*client.Client, error) {
	if err := microdb.LoadDataOrigins("dataorigin.yaml"); err != nil {
		return nil, fmt.Errorf("failed to load data origins: %w", err)
	}

	c, err := client.Connect(natsHost, natsPort, natsClientID, natsClusterID, tables...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to microdb: %w", err)
	}
	return c, nil
}
