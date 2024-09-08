package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type VersionedDatabase interface {
	DB() *sql.DB
	getCurrentVersion() (int, error)
	setVersion(version int) error
}

func MigrateUp(db VersionedDatabase, version int, query string) error {
	if version < 1 {
		return errors.New("version less than 1")
	}

	currentVersion, err := db.getCurrentVersion()
	if err != nil {
		return err
	}

	if currentVersion >= version {
		return errors.New(fmt.Sprintf("current version (%d) greater than version (%d)", currentVersion, version))
	}

	_, err = db.DB().Exec(query)
	if err != nil {
		return err
	}

	err = db.setVersion(version)
	if err != nil {
		return err
	}

	return nil
}

func MigrateDown(db VersionedDatabase, version int, query string) error {
	if version < 1 {
		return errors.New("version less than 1")
	}

	currentVersion, err := db.getCurrentVersion()
	if err != nil {
		return err
	}

	if currentVersion < version {
		return errors.New(fmt.Sprintf("current version (%d) greater than version (%d)", currentVersion, version))
	}

	_, err = db.DB().Exec(query)
	if err != nil {
		return err
	}

	err = db.setVersion(version - 1)
	if err != nil {
		return err
	}

	return nil
}
