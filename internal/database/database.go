package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type VersionedDatabase interface {
	DB() *sql.DB
	GetCurrentVersion() (int, error)
	SetVersion(version int) error
}

func MigrateUp(db VersionedDatabase, version int, query string) error {
	if version < 1 {
		return errors.New("version less than 1")
	}

	currentVersion, err := db.GetCurrentVersion()
	if err != nil {
		return err
	}

	if currentVersion >= version {
		return fmt.Errorf("current version (%d) greater than version (%d)", currentVersion, version)
	}

	_, err = db.DB().Exec(query)
	if err != nil {
		return err
	}

	err = db.SetVersion(version)
	if err != nil {
		return err
	}

	return nil
}

func MigrateDown(db VersionedDatabase, version int, query string) error {
	if version < 0 {
		return errors.New("version less than 0")
	}

	currentVersion, err := db.GetCurrentVersion()
	if err != nil {
		return err
	}

	if currentVersion < version {
		return fmt.Errorf("current version (%d) greater than version (%d)", currentVersion, version)
	}

	_, err = db.DB().Exec(query)
	if err != nil {
		return err
	}

	err = db.SetVersion(version - 1)
	if err != nil {
		return err
	}

	return nil
}

func ApplyMigrations(db VersionedDatabase, up []string, down []string, version ...int) error {
	if len(up) != len(down) {
		return errors.New("forward and backward migration count mismatch")
	}

	currentVersion, err := db.GetCurrentVersion()
	if err != nil {
		return err
	}

	targetVersion := len(up)
	if len(version) >= 1 {
		targetVersion = version[0]

		if targetVersion < 0 {
			targetVersion = 0
		} else if targetVersion > len(up) {
			targetVersion = len(up)
		}
	}

	if currentVersion < targetVersion {
		for i := currentVersion; i < targetVersion; i += 1 {
			if err = MigrateUp(db, i+1, up[i]); err != nil {
				return err
			}
		}

		if err = db.SetVersion(targetVersion); err != nil {
			return err
		}
	} else if currentVersion > targetVersion {
		for i := len(down); i > targetVersion; i -= 1 {
			if err = MigrateDown(db, i-1, down[i-1]); err != nil {
				return err
			}
		}

		if err = db.SetVersion(targetVersion); err != nil {
			return err
		}
	}

	return nil
}
