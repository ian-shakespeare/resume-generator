package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/tests"
	"testing"
)

func TestCreateProject(t *testing.T) {
	t.Run("allFields", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&user,
			"name",
			"email",
			"phone number",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		_, err = database.CreateProject(db, &resume, "name", "description", "role")

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})
}

func TestGetProject(t *testing.T) {
	t.Run("allFields", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&user,
			"name",
			"email",
			"phone number",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		created, err := database.CreateProject(db, &resume, "name", "description", "role")

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		project := database.GetProject(db, created.Id)

		if project == nil {
			t.Fatalf("expected %s, received %s", "project", "nil")
		}
	})
}

func TestCreateProjectResposibility(t *testing.T) {
	t.Run("allFields", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&user,
			"name",
			"email",
			"phone number",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		project, err := database.CreateProject(db, &resume, "name", "description", "role")
		if err != nil {
			t.Fatalf("expected %s, received %s", "project", "nil")
		}

		responsibility, err := database.CreateProjectResponsibility(db, &project, "responsibility")

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(project.Responsibilities) < 1 {
			t.Fatalf("expected %d, received %d", 1, len(project.Responsibilities))
		} else if project.Responsibilities[0].Id != responsibility.Id {
			t.Fatalf("expected %s, received %s", responsibility.Id, project.Responsibilities[0].Id)
		}
	})
}
