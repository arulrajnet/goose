package goose

import (
	"database/sql"
	"fmt"
)

// Down rolls back a single migration from the current version.
func Down(db *sql.DB, dir string) error {
	migrations, err := CollectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}
	currentVersion, err := GetDBVersion(db)
	if err != nil {
		return err
	}
	current, err := migrations.Current(currentVersion)
	if err != nil {
		return fmt.Errorf("no migration %v", currentVersion)
	}
	return current.Down(db)
}

// DownTo rolls back migrations to a specific version.
func DownTo(db *sql.DB, dir string, version int64, opts ...OptionsFunc) error {
	option := &options{}
	for _, f := range opts {
		f(option)
	}
	migrations, err := CollectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}
	if option.noVersioning {
		return downNoVersioning(db, migrations, version)
	}

	for {
		currentVersion, err := GetDBVersion(db)
		if err != nil {
			return err
		}

		current, err := migrations.Current(currentVersion)
		if err != nil {
			log.Printf("goose: no migrations to run. current version: %d\n", currentVersion)
			return nil
		}

		if current.Version <= version {
			log.Printf("goose: no migrations to run. current version: %d\n", currentVersion)
			return nil
		}

		if err = current.Down(db); err != nil {
			return err
		}
	}
}

func downNoVersioning(db *sql.DB, migrations Migrations, targetVersion int64) error {
	// TODO(mf): we're not tracking the seed migrations in the database,
	// which means subsequent "down" operations will always start from the
	// highest seed file.
	// Also, should target version always be 0 and error otherwise?
	for _, current := range migrations {
		current.noVersioning = true
		if err := current.Down(db); err != nil {
			return err
		}
		if current.Version <= targetVersion {
			log.Printf("goose: no migrations to run. current version: %d\n", current.Version)
			return nil
		}
	}
	return nil
}
