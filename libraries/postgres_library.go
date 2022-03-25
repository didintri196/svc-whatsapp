package libraries

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresLibrary struct {
	MigrationDirectory string
	MigrationDialect   string
	DBHost             string
	DBUser             string
	DBPassword         string
	DBPort             string
	DBName             string
	DBSSLMode          string
	LogMode            int
}

func (lib PostgresLibrary) Migrate(db *sql.DB) (err error) {
	migrations := &migrate.FileMigrationSource{
		Dir: lib.MigrationDirectory,
	}

	n, err := migrate.Exec(db, lib.MigrationDialect, migrations, migrate.Up)
	if err != nil {
		return err
	}

	if n > 0 {
		fmt.Printf("Migration: Applied %d migrations \n", n)
	} else {
		fmt.Printf("Migration: No schema change \n")
	}

	return err
}

func (lib PostgresLibrary) ConnectAndValidate() (db *gorm.DB, sql *sql.DB, err error) {
	// init config
	postgresDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		lib.DBHost,
		lib.DBUser,
		lib.DBPassword,
		lib.DBName,
		lib.DBPort,
		lib.DBSSLMode)
	// log mode
	logConfig := logger.Config{
		SlowThreshold: time.Second,
		Colorful:      true,
	}
	switch lib.LogMode {
	case 1:
		logConfig.LogLevel = logger.Info
		logConfig.IgnoreRecordNotFoundError = false
	case 2:
		logConfig.LogLevel = logger.Warn
		logConfig.IgnoreRecordNotFoundError = true
	case 3:
		logConfig.LogLevel = logger.Error
		logConfig.IgnoreRecordNotFoundError = true
	default:
		logConfig = logger.Config{}
	}

	// open connection
	gormConfig := gorm.Config{Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logConfig)}
	db, err = gorm.Open(postgres.Open(postgresDsn), &gormConfig)
	if err != nil {
		return db, sql, err
	}

	// connection configuration
	sql, err = db.DB()
	if err != nil {
		return db, sql, err
	}

	// ping connection
	err = sql.Ping()
	if err != nil {
		return db, sql, err
	}
	return db, sql, err
}
