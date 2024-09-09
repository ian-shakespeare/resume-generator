package database_test

import (
	"resumegenerator/internal/database"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrateUp(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	query := `
CREATE TABLE migrate_test_up (
  id INT PRIMARY KEY
)`

	// Act
	err := database.MigrateUp(db, 1, query)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	// Assert
	userVersion, err := db.GetCurrentVersion()
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	if userVersion != 1 {
		t.Fatalf("expected %s, received %d", "1", userVersion)
	}

	_, err = db.DB().Exec("INSERT INTO migrate_test_up (id) VALUES (1)")
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}
}

func TestMigrateUpToZero(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Act
	err := database.MigrateUp(db, 0, "")

	// Assert
	if err == nil {
		t.Fatalf("expected %s, received %s", "error", "nil")
	}
}

func TestMigrateUpExisting(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	query1 := `
CREATE TABLE migrate_test_up_existing1 (
  id INT PRIMARY KEY
)`

	err := database.MigrateUp(db, 1, query1)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	query2 := `
CREATE TABLE migrate_test_up_existing2 (
  id INT PRIMARY KEY
)`

	// Act
	err = database.MigrateUp(db, 1, query2)

	// Assert
	if err == nil {
		t.Fatalf("expected %s, received %s", "error", "nil")
	}
}

func TestMigrateDown(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	upQuery := `
CREATE TABLE migrate_test_down (
  id INT PRIMARY KEY
)`

	err := database.MigrateUp(db, 1, upQuery)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	// Act
	err = database.MigrateDown(db, 1, "DROP TABLE migrate_test_down")
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	// Assert
	userVersion, err := db.GetCurrentVersion()
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	if userVersion != 0 {
		t.Fatalf("expected %s, received %d", "0", userVersion)
	}

	_, err = db.DB().Exec("INSERT INTO migrate_test_down (id) VALUES (1)")
	if err == nil {
		t.Fatalf("expected %s, received %s", "error", "nil")
	}
}

func TestMigrateDownNegative(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Act
	err := database.MigrateDown(db, -1, "")

	// Assert
	if err == nil {
		t.Fatalf("expected %s, received %s", "error", "nil")
	}
}

func TestMigrateDownNonExisting(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Act
	err := database.MigrateDown(db, 1, "")

	// Assert
	if err == nil {
		t.Fatalf("expected %s, received %s", "error", "nil")
	}
}

func TestApplyMigrations(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	up := []string{
		"CREATE TABLE apply_migrations1 (id TEXT)",
		"CREATE TABLE apply_migrations2 (id TEXT)",
		"CREATE TABLE apply_migrations3 (id TEXT)",
	}

	down := []string{
		"DROP TABLE apply_migrations1",
		"DROP TABLE apply_migrations2",
		"DROP TABLE apply_migrations3",
	}

	// Act
	err := database.ApplyMigrations(db, up, down)

	// Assert
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	userVersion, err := db.GetCurrentVersion()
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	if userVersion != 3 {
		t.Fatalf("expected %d, received %d", 3, userVersion)
	}

	_, err = db.DB().Exec("INSERT INTO apply_migrations3 (id) VALUES ('SOME_ID')")
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}
}

func TestApplySpecificMigration(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	up := []string{
		"CREATE TABLE apply_migrations1 (id TEXT)",
		"CREATE TABLE apply_migrations2 (id TEXT)",
		"CREATE TABLE apply_migrations3 (id TEXT)",
	}

	down := []string{
		"DROP TABLE apply_migrations1",
		"DROP TABLE apply_migrations2",
		"DROP TABLE apply_migrations3",
	}

	// Act
	err := database.ApplyMigrations(db, up, down, 2)

	// Assert
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	userVersion, err := db.GetCurrentVersion()
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	if userVersion != 2 {
		t.Fatalf("expected %d, received %d", 2, userVersion)
	}

	_, err = db.DB().Exec("INSERT INTO apply_migrations2 (id) VALUES ('SOME_ID')")
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	_, err = db.DB().Exec("INSERT INTO apply_migrations3 (id) VALUES ('SOME_ID')")
	if err == nil {
		t.Fatalf("expected %s, received %s", "error", "nil")
	}
}

func TestApplyMigrationsDown(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	up := []string{
		"CREATE TABLE apply_migrations1 (id TEXT)",
		"CREATE TABLE apply_migrations2 (id TEXT)",
		"CREATE TABLE apply_migrations3 (id TEXT)",
	}

	down := []string{
		"DROP TABLE apply_migrations1",
		"DROP TABLE apply_migrations2",
		"DROP TABLE apply_migrations3",
	}

	err := database.ApplyMigrations(db, up, down)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	// Act
	err = database.ApplyMigrations(db, up, down, 1)

	// Assert
	userVersion, err := db.GetCurrentVersion()
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	if userVersion != 1 {
		t.Fatalf("expected %d, received %d", 1, userVersion)
	}

	_, err = db.DB().Exec("INSERT INTO apply_migrations1 (id) VALUES ('SOME_ID')")
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	_, err = db.DB().Exec("INSERT INTO apply_migrations2 (id) VALUES ('SOME_ID')")
	if err == nil {
		t.Fatalf("expected %s, received %s", "error", "nil")
	}
}

func TestApplyMigrationsDownAll(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	up := []string{
		"CREATE TABLE apply_migrations1 (id TEXT)",
		"CREATE TABLE apply_migrations2 (id TEXT)",
		"CREATE TABLE apply_migrations3 (id TEXT)",
	}

	down := []string{
		"DROP TABLE apply_migrations1",
		"DROP TABLE apply_migrations2",
		"DROP TABLE apply_migrations3",
	}

	err := database.ApplyMigrations(db, up, down)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	// Act
	err = database.ApplyMigrations(db, up, down, 0)

	// Assert
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	userVersion, err := db.GetCurrentVersion()
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	if userVersion != 0 {
		t.Fatalf("expected %d, received %d", 0, userVersion)
	}

	_, err = db.DB().Exec("INSERT INTO apply_migrations1 (id) VALUES ('SOME_ID')")
	if err == nil {
		t.Fatalf("expected %s, received %s", "error", "nil")
	}
}
