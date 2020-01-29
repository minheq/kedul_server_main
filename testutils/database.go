package testutils

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// SetupDB connects to database and runs migrations
func SetupDB() (db *sql.DB, cleanup func() error) {
	db, _ = sql.Open("postgres", "postgres://postgres@127.0.0.1:5432/kedul?sslmode=disable")
	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	migrations, _ := migrate.NewWithDatabaseInstance("file://../migrations", "kedul", driver)

	migrations.Down()
	migrations.Up()

	return db, db.Close
}
