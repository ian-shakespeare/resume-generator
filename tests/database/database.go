package database_test

import (
	"os"
	"resumegenerator/internal/database"
	"testing"
)

const TEST_DB_NAME string = "testing.db"

func setup(t *testing.T) database.VersionedDatabase {
	file, err := os.Create(TEST_DB_NAME)
	if err != nil {
		t.Fatalf("setup %s", err.Error())
	}
	file.Close()

	db, err := database.NewSQLite(TEST_DB_NAME)
	if err != nil {
		t.Fatalf("setup %s", err.Error())
	}

	return db
}

func tearDown(t *testing.T, db database.VersionedDatabase) {
	db.DB().Close()
	err := os.Remove(TEST_DB_NAME)
	if err != nil {
		t.Fatalf("teardown %s", err.Error())
	}
}
