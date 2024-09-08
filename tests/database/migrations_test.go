package database_test

import (
	"os"
	"resumegenerator/internal/database"
	"resumegenerator/tests"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrateUp(t *testing.T) {
	setup(t)
	defer tearDown(t)

	// Arrange
	db, err := database.NewSQLite(TEST_DB_NAME)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}
	defer db.DB().Close()
	query := `
CREATE TABLE migrate_test_up (
  id INT PRIMARY KEY
)`

	// Act
	err = database.MigrateUp(db, 1, query)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}

	// Assert
	userVersionRow := db.DB().QueryRow("PRAGMA user_version")
	if userVersionRow == nil {
		t.Fatal(tests.FormatExpected("*sql.row", "nil"))
	}
	userVersionBytes := make([]byte, 1)
	err = userVersionRow.Scan(&userVersionBytes)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}

	userVersion, err := strconv.Atoi(string(userVersionBytes))
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}

	if userVersion != 1 {
		t.Fatal(tests.FormatExpected("1", string(userVersionBytes)))
	}

	_, err = db.DB().Exec("INSERT INTO migrate_test_up (id) VALUES (1)")
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}
}

func TestMigrateUpToZero(t *testing.T) {
	setup(t)
	defer tearDown(t)

	// Arrange
	db, err := database.NewSQLite(TEST_DB_NAME)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}
	defer db.DB().Close()

	// Act
	err = database.MigrateUp(db, 0, "")

	// Assert
	if err == nil {
		t.Fatal(tests.FormatExpected("error", "nil"))
	}
}

func TestMigrateUpExisting(t *testing.T) {
	setup(t)
	defer tearDown(t)

	// Arrange
	db, err := database.NewSQLite(TEST_DB_NAME)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}
	defer db.DB().Close()

	query1 := `
CREATE TABLE migrate_test_up_existing1 (
  id INT PRIMARY KEY
)`

	err = database.MigrateUp(db, 1, query1)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}

	query2 := `
CREATE TABLE migrate_test_up_existing2 (
  id INT PRIMARY KEY
)`

	// Act
	err = database.MigrateUp(db, 1, query2)

	// Assert
	if err == nil {
		t.Fatal(tests.FormatExpected("error", "nil"))
	}
}

func TestMigrateDown(t *testing.T) {
	setup(t)
	defer tearDown(t)

	// Arrange
	db, err := database.NewSQLite(TEST_DB_NAME)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}
	defer db.DB().Close()

	upQuery := `
CREATE TABLE migrate_test_down (
  id INT PRIMARY KEY
)`

	err = database.MigrateUp(db, 1, upQuery)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}

	// Act
	err = database.MigrateDown(db, 1, "DROP TABLE migrate_test_down")
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}

	// Assert
	userVersionRow := db.DB().QueryRow("PRAGMA user_version")
	if userVersionRow == nil {
		t.Fatal(tests.FormatExpected("*sql.Row", "nil"))
	}
	userVersionBytes := make([]byte, 1)
	err = userVersionRow.Scan(&userVersionBytes)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}

	userVersion, err := strconv.Atoi(string(userVersionBytes))
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}

	if userVersion != 0 {
		t.Fatal(tests.FormatExpected("0", string(userVersionBytes)))
	}

	_, err = db.DB().Exec("INSERT INTO migrate_test_down (id) VALUES (1)")
	if err == nil {
		t.Fatal(tests.FormatExpected("error", "nil"))
	}
}

func TestMigrateDownPastZero(t *testing.T) {
	setup(t)
	defer tearDown(t)

	// Arrange
	db, err := database.NewSQLite(TEST_DB_NAME)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}
	defer db.DB().Close()

	// Act
	err = database.MigrateDown(db, 0, "")

	// Assert
	if err == nil {
		t.Fatal(tests.FormatExpected("error", "nil"))
	}
}

func TestMigrateDownNonExisting(t *testing.T) {
	setup(t)
	defer tearDown(t)

	// Arrange
	db, err := database.NewSQLite(TEST_DB_NAME)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}
	defer db.DB().Close()

	// Act
	err = database.MigrateDown(db, 1, "")

	// Assert
	if err == nil {
		t.Fatal(tests.FormatExpected("error", "nil"))
	}
}

func TestMigrations(t *testing.T) {
	// Arrange
	testDbConn, exists := os.LookupEnv("TEST_DB_CONN")
	if !exists {
		t.Fatal(tests.FormatExpected("true", "false"))
	}

	db, err := database.NewPostgres(testDbConn)
	if err != nil {
		t.Fatal(tests.FormatExpected("nil", err.Error()))
	}
	defer db.DB().Close()

	forwards := database.ForwardMigrations()
	for i := 0; i < len(forwards); i += 1 {
		// Act
		err = database.MigrateUp(db, i+1, forwards[i])

		// Assert
		if err != nil {
			t.Fatal(tests.FormatExpected("nil", err.Error()))
		}
	}

	backwards := database.BackwardMigrations()
	for i := len(backwards) - 1; i >= 0; i -= 1 {
		// Act
		err = database.MigrateDown(db, i+1, backwards[i])

		// Assert
		if err != nil {
			t.Fatal(tests.FormatExpected("nil", err.Error()))
		}
	}
}
