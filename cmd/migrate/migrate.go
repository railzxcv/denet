package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		url   string
		path  string
		force bool
		steps int
	)
	flag.StringVar(&url, "url", "", "Database connection URL")
	flag.StringVar(&path, "path", "file://migrations", "Path to migrations directory")
	flag.BoolVar(&force, "force", false, "Force migration to run even if it is dirty (default false)")
	flag.IntVar(&steps, "steps", 0, "Number of migration steps to apply (0 for all)")
	flag.Parse()

	if url == "" {
		log.Fatal("Database URL must be specified with --url")
	}

	m, err := migrate.New(path, url)
	if err != nil {
		log.Fatal("migrate error: ", err)
	}
	defer m.Close()
	if force {
		if steps == 0 {
			log.Fatal("Cannot use --force without specifying --steps")
		}
		err = m.Force(steps)
	} else {
		err = m.Steps(steps)
	}
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("migrate error: ", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal(err)
	}

	fmt.Printf("Migration successful! Current version: %d (dirty: %t)\n", version, dirty)
}
