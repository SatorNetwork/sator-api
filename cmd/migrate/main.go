package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // init pg driver

	"github.com/dmitrymomot/go-env"
	migrate "github.com/rubenv/sql-migrate"
)

var (
	dbConnString    = env.MustString("DATABASE_URL")
	migrationsDir   = env.GetString("DATABASE_MIGRATIONS_DIR", "./migrations")
	migrationsTable = env.GetString("DATABASE_MIGRATIONS_TABLE", "migrations")
)

func main() {
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("init db connection error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("db pinng error: %v", err)
	}

	m := migrate.MigrationSet{
		TableName: migrationsTable,
	}
	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}
	n, err := m.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("could not exec migrations: %v", err)
	}

	fmt.Printf("Applied %d migrations!\n", n)
}
