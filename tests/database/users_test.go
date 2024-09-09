package database_test

import (
	"resumegenerator/internal/database"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations(), 1)

	// Act
	_, err := database.CreateUser(db)

	// Assert
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}
}

func TestGetUser(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())

	created, err := database.CreateUser(db)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	// Act
	user := database.GetUser(db, created.Id)

	// Assert
	if user == nil {
		t.Fatalf("expected %s, received %s", "user", "nil")
	}
}
