package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDatabase struct {
	db *sql.DB
}

func NewSQLite(path string) (*SQLiteDatabase, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	return &SQLiteDatabase{db}, nil
}

func (s *SQLiteDatabase) DB() *sql.DB {
	return s.db
}

func (s *SQLiteDatabase) getCurrentVersion() (int, error) {
	currentVersionRow := s.db.QueryRow("PRAGMA user_version")
	if currentVersionRow == nil {
		return 0, errors.New("missing user_version row")
	}

	currentVersionBuf := make([]byte, 1)
	if err := currentVersionRow.Scan(&currentVersionBuf); err != nil {
		return 0, err
	}

	currentVersion, err := strconv.Atoi(string(currentVersionBuf))
	if err != nil {
		return 0, err
	}

	return currentVersion, nil
}

func (s *SQLiteDatabase) setVersion(version int) error {
	pragmaQuery := fmt.Sprintf("PRAGMA user_version = %d", version)
	_, err := s.db.Exec(pragmaQuery)
	return err
}
