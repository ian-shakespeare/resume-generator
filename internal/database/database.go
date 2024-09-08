package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func MigrateUp(db *sql.DB, version int, query string) error {
	if version < 1 {
		return errors.New("version less than 1")
	}

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	pragmaQuery := fmt.Sprintf("PRAGMA user_version = %d", version)
	_, err = db.Exec(pragmaQuery)
	if err != nil {
		return err
	}

	return nil
}

func MigrateDown(db *sql.DB, version int, query string) error {
	if version < 1 {
		return errors.New("version less than 1")
	}

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	pragmaQuery := fmt.Sprintf("PRAGMA user_version = %d", version-1)
	_, err = db.Exec(pragmaQuery)
	if err != nil {
		return err
	}

	return nil
}
