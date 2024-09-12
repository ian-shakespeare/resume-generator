package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/pkg/resume"
	"resumegenerator/test"
	"testing"
)

func TestCreateProject(t *testing.T) {
	t.Run("allFields", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.FromJSON([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		p, err := resume.ProjectFromJSON([]byte(test.PROJECT))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateProject(db, &r, &p)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})
}

func TestGetProject(t *testing.T) {
	t.Run("allFields", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.FromJSON([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		p, err := resume.ProjectFromJSON([]byte(test.PROJECT))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateProject(db, &r, &p)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		project := database.GetProject(db, p.Id)

		if project == nil {
			t.Fatalf("expected %s, received %s", "project", "nil")
		}
	})
}
