package testutils

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// SetupDB connects to database and runs migrations
func SetupDB() (db *sql.DB, cleanup func() error) {
	dbURL := os.Getenv("DATABASE_URL")
	databaseName := os.Getenv("DATABASE_NAME")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("error opening database with DATABASE_URL=%s", dbURL)
	}

	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	migrations, err := migrate.NewWithDatabaseInstance("file://../migrations", databaseName, driver)

	if err != nil {
		log.Fatal("error instantiating migrations")
	}

	err = migrations.Down()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("error running migrations down, %v", err)
	}

	err = migrations.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("error running migrations up, %v", err)
	}

	return db, db.Close
}
