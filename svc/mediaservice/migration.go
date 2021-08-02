package mediaservice

import (
	"database/sql"
	"fmt"
	"log"

	migrate "github.com/rubenv/sql-migrate"
)

// Migrate runs database migrations
func Migrate(db *sql.DB, driver, dir string) error {
	migrations := &migrate.FileMigrationSource{
		Dir: dir,
	}

	n, err := migrate.Exec(db, driver, migrations, migrate.Up)
	if err != nil {
		return err
	}

	log.Println(fmt.Printf("Applied %d migrations!\n", n))

	return nil
}
