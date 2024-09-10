package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/tests"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Run("allFields", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations(), 1)

		_, err := database.CreateUser(db)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})
}

func TestGetUser(t *testing.T) {
	t.Run("allFields", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())

		created, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user := database.GetUser(db, created.Id)

		if user == nil {
			t.Fatalf("expected %s, received %s", "user", "nil")
		}
	})
}
