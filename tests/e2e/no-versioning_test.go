package e2e

import (
	"database/sql"
	"testing"

	"github.com/matryer/is"
	"github.com/pressly/goose/v3"
)

func TestUpNoVersioning(t *testing.T) {
	const (
		wantSeedOwnerCount = 250
		wantOwnerCount     = 5
	)
	is := is.New(t)
	db, err := newDockerDB(t)
	is.NoErr(err)
	goose.SetDialect(*dialect)

	err = goose.Up(db, migrationsDir)
	is.NoErr(err)
	baseVersion, err := goose.GetDBVersion(db)
	is.NoErr(err)

	// Run (all) up migrations from the seed dir
	{
		err = goose.Up(db, seedDir, goose.WithNoVersioning())
		is.NoErr(err)
		// Confirm no changes to the versioned schema in the DB
		currentVersion, err := goose.GetDBVersion(db)
		is.NoErr(err)
		is.Equal(baseVersion, currentVersion)
		seedOwnerCount, err := countSeedOwners(db)
		is.NoErr(err)
		is.Equal(seedOwnerCount, wantSeedOwnerCount)
	}

	// Run (all) down migrations from the seed dir
	{
		err = goose.Down(db, seedDir, goose.WithNoVersioning())
		is.NoErr(err)
		// Confirm no changes to the versioned schema in the DB
		currentVersion, err := goose.GetDBVersion(db)
		is.NoErr(err)
		is.Equal(baseVersion, currentVersion)
		seedOwnerCount, err := countSeedOwners(db)
		is.NoErr(err)
		is.Equal(seedOwnerCount, 0)
	}

	// The migrations added 4 non-seed owners, they must remain
	// in the database afterwards
	ownerCount, err := countOwners(db)
	is.NoErr(err)
	is.Equal(ownerCount, wantOwnerCount)
}

func countSeedOwners(db *sql.DB) (int, error) {
	q := `SELECT count(*)FROM owners WHERE owner_name LIKE'seed-user-%'`
	var count int
	err := db.QueryRow(q).Scan(&count)
	return count, err
}

func countOwners(db *sql.DB) (int, error) {
	q := `SELECT count(*)FROM owners`
	var count int
	err := db.QueryRow(q).Scan(&count)
	return count, err
}
