package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/test"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrateUp(t *testing.T) {
	t.Run("toFutureVersion", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		query := `
CREATE TABLE migrate_test_up (
  id INT PRIMARY KEY
)`

		err := database.MigrateUp(conn, 1, query)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		userVersion, err := conn.DatabaseVersion()
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if userVersion != 1 {
			t.Fatalf("expected %s, received %d", "1", userVersion)
		}

		_, err = conn.Database().Exec("INSERT INTO migrate_test_up (id) VALUES (1)")
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})

	t.Run("toInitialVersion", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.MigrateUp(conn, 0, "")

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})

	t.Run("toPastVersion", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		query1 := `
CREATE TABLE migrate_test_up_existing1 (
  id INT PRIMARY KEY
)`

		err := database.MigrateUp(conn, 1, query1)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		query2 := `
CREATE TABLE migrate_test_up_existing2 (
  id INT PRIMARY KEY
)`

		err = database.MigrateUp(conn, 1, query2)

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})
}

func TestMigrateDown(t *testing.T) {
	t.Run("toPastVersion", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		upQuery := `
CREATE TABLE migrate_test_down (
  id INT PRIMARY KEY
)`

		err := database.MigrateUp(conn, 1, upQuery)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.MigrateDown(conn, 1, "DROP TABLE migrate_test_down")
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		userVersion, err := conn.DatabaseVersion()
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if userVersion != 0 {
			t.Fatalf("expected %s, received %d", "0", userVersion)
		}

		_, err = conn.Database().Exec("INSERT INTO migrate_test_down (id) VALUES (1)")
		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})

	t.Run("toNegativeVersion", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.MigrateDown(conn, -1, "")

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})

	t.Run("toNonExistantVersion", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.MigrateDown(conn, 1, "")

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})
}

func TestApplyMigrations(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

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

		err := database.ApplyMigrations(conn, up, down)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		userVersion, err := conn.DatabaseVersion()
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if userVersion != 3 {
			t.Fatalf("expected %d, received %d", 3, userVersion)
		}

		_, err = conn.Database().Exec("INSERT INTO apply_migrations3 (id) VALUES ('SOME_ID')")
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})

	t.Run("specific", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

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

		err := database.ApplyMigrations(conn, up, down, 2)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		userVersion, err := conn.DatabaseVersion()
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if userVersion != 2 {
			t.Fatalf("expected %d, received %d", 2, userVersion)
		}

		_, err = conn.Database().Exec("INSERT INTO apply_migrations2 (id) VALUES ('SOME_ID')")
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		_, err = conn.Database().Exec("INSERT INTO apply_migrations3 (id) VALUES ('SOME_ID')")
		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})

	t.Run("down", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

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

		err := database.ApplyMigrations(conn, up, down)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.ApplyMigrations(conn, up, down, 1)

		userVersion, err := conn.DatabaseVersion()
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if userVersion != 1 {
			t.Fatalf("expected %d, received %d", 1, userVersion)
		}

		_, err = conn.Database().Exec("INSERT INTO apply_migrations1 (id) VALUES ('SOME_ID')")
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		_, err = conn.Database().Exec("INSERT INTO apply_migrations2 (id) VALUES ('SOME_ID')")
		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})

	t.Run("downAll", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

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

		err := database.ApplyMigrations(conn, up, down)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		// Act
		err = database.ApplyMigrations(conn, up, down, 0)

		// Assert
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		userVersion, err := conn.DatabaseVersion()
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if userVersion != 0 {
			t.Fatalf("expected %d, received %d", 0, userVersion)
		}

		_, err = conn.Database().Exec("INSERT INTO apply_migrations1 (id) VALUES ('SOME_ID')")
		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})
}
