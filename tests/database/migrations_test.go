package database_test

import (
	"database/sql"
	"fmt"
	"os"
	"resumegenerator/internal/database"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const TEST_DB_NAME string = "testing.db"

func setup(t *testing.T) {
	file, err := os.Create(TEST_DB_NAME)
	if err != nil {
		t.Fatalf("setup %s", err.Error())
	}
	file.Close()
}

func tearDown(t *testing.T) {
	err := os.Remove(TEST_DB_NAME)
	if err != nil {
		t.Fatalf("teardown %s", err.Error())
	}
}

func formatExpected(expected string, received string) string {
	return fmt.Sprintf("expected %s, received %s", expected, received)
}

func TestMigrateUp(t *testing.T) {
	setup(t)
	defer tearDown(t)

	// Arrange
	db, err := sql.Open("sqlite3", TEST_DB_NAME)
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}
	query := `
CREATE TABLE migrate_test_up (
  id INT PRIMARY KEY
)`

	// Act
	err = database.MigrateUp(db, 1, query)
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}

	// Assert
	userVersionRow := db.QueryRow("PRAGMA user_version")
	if userVersionRow == nil {
		t.Fatal(formatExpected("*sql.row", "nil"))
	}
	userVersionBytes := make([]byte, 1)
	err = userVersionRow.Scan(&userVersionBytes)
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}

	userVersion, err := strconv.Atoi(string(userVersionBytes))
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}

	if userVersion != 1 {
		t.Fatal(formatExpected("1", string(userVersionBytes)))
	}

	_, err = db.Exec("INSERT INTO migrate_test_up (id) VALUES (1)")
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}
}

func TestMigrateUpToZero(t *testing.T) {
	setup(t)
	defer tearDown(t)

	// Arrange
	db, err := sql.Open("sqlite3", TEST_DB_NAME)
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}

	// Act
	err = database.MigrateUp(db, 0, "")

	// Assert
	if err == nil {
		t.Fatal(formatExpected("error", "nil"))
	}
}

func TestMigrateDown(t *testing.T) {
	setup(t)
	defer tearDown(t)

	// Arrange
	db, err := sql.Open("sqlite3", TEST_DB_NAME)
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}

	upQuery := `
CREATE TABLE migrate_test_down (
  id INT PRIMARY KEY
)`

	err = database.MigrateUp(db, 1, upQuery)
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}

	// Act
	err = database.MigrateDown(db, 1, "DROP TABLE migrate_test_down")
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}

	// Assert
	userVersionRow := db.QueryRow("PRAGMA user_version")
	if userVersionRow == nil {
		t.Fatal(formatExpected("*sql.Row", "nil"))
	}
	userVersionBytes := make([]byte, 1)
	err = userVersionRow.Scan(&userVersionBytes)
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}

	userVersion, err := strconv.Atoi(string(userVersionBytes))
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}

	if userVersion != 0 {
		t.Fatal(formatExpected("0", string(userVersionBytes)))
	}

	_, err = db.Exec("INSERT INTO migrate_test_down (id) VALUES (1)")
	if err == nil {
		t.Fatal(formatExpected("error", "nil"))
	}
}

func TestMigrateDownPastZero(t *testing.T) {
	setup(t)
	defer tearDown(t)

	// Arrange
	db, err := sql.Open("sqlite3", TEST_DB_NAME)
	if err != nil {
		t.Fatal(formatExpected("nil", err.Error()))
	}

	// Act
	err = database.MigrateDown(db, 0, "")

	// Assert
	if err == nil {
		t.Fatal(formatExpected("error", "nil"))
	}
}
