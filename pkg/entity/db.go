package entity

import (
	"context"
	"github.com/checkr/flagr/pkg/instrumentation/instana"
	"os"
	"sync"

	_ "github.com/jinzhu/gorm/dialects/mysql"    // mysql driver
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // sqlite driver

	retry "github.com/avast/retry-go"
	"github.com/checkr/flagr/pkg/config"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var (
	singletonDB    *gorm.DB
	singletonOnce  sync.Once
	instanaAdapter func(context.Context, *gorm.DB) *gorm.DB
)

// AutoMigrateTables stores the entity tables that we can auto migrate in gorm
var AutoMigrateTables = []interface{}{
	Constraint{},
	Distribution{},
	FlagSnapshot{},
	Flag{},
	Segment{},
	User{},
	Variant{},
	Tag{},
	FlagEntityType{},
}

func connectDB() (db *gorm.DB, err error) {
	err = retry.Do(
		func() error {
			db, err = gorm.Open(config.Config.DBDriver, config.Config.DBConnectionStr)
			return err
		},
		retry.Attempts(config.Config.DBConnectionRetryAttempts),
		retry.Delay(config.Config.DBConnectionRetryDelay),
	)
	return db, err
}

// GetDB gets the db singleton
func GetDB(ctx context.Context) *gorm.DB {
	singletonOnce.Do(func() {
		db, err := connectDB()
		if err != nil {
			if config.Config.DBConnectionDebug {
				logrus.WithField("err", err).Fatal("failed to connect to db")
			} else {
				logrus.Fatal("failed to connect to db")
			}
		}
		if err := registerInstanaCallback(db); err != nil {
			if config.Config.DBConnectionDebug {
				logrus.WithField("err", err).Fatal("failed to inject Instana to db")
			} else {
				logrus.Fatal("failed to inject Instana to db")
			}
		}
		db.SetLogger(logrus.StandardLogger())
		db.Debug().AutoMigrate(AutoMigrateTables...)
		singletonDB = db
	})

	if config.Config.InstanaEnabled {
		return instanaAdapter(ctx, singletonDB)
	}
	return singletonDB
}

// NewSQLiteDB creates a new sqlite db
// useful for backup exports and unit tests
func NewSQLiteDB(filePath string) *gorm.DB {
	os.Remove(filePath)
	db, err := gorm.Open("sqlite3", filePath)
	if err != nil {
		logrus.WithField("err", err).Errorf("failed to connect to db:%s", filePath)
		panic(err)
	}
	db.SetLogger(logrus.StandardLogger())
	db.AutoMigrate(AutoMigrateTables...)
	return db
}

// NewTestDB creates a new test db
func NewTestDB() *gorm.DB {
	return NewSQLiteDB(":memory:")
}

// PopulateTestDB seeds the test db
func PopulateTestDB(flag Flag) *gorm.DB {
	testDB := NewTestDB()
	testDB.Create(&flag)
	return testDB
}

// registerInstanaCallback registers necessary callbacks for Instana tracing of database operations
func registerInstanaCallback(db *gorm.DB) error {
	if config.Config.InstanaEnabled {
		ia, err := instana.NewAdapter(db, config.Config.DBDriver, config.Config.DBConnectionStr)
		if err != nil {
			return err
		}
		instanaAdapter = ia
	}
	return nil
}
