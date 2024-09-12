package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/pkg/resume"
	"resumegenerator/test"
	"testing"
)

func TestCreateEducation(t *testing.T) {
	t.Run("minimumRequiredFields", func(t *testing.T) {
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

		r, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		e, err := resume.NewEducation([]byte(test.MIN_EDUCATION))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateEducation(db, &r, &e)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})

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

		r, err := resume.New([]byte(test.FULL_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		e, err := resume.NewEducation([]byte(test.FULL_EDUCATION))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateEducation(db, &r, &e)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})
}

func TestGetEducation(t *testing.T) {
	t.Run("minimumRequiredFields", func(t *testing.T) {
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

		r, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		e, err := resume.NewEducation([]byte(test.MIN_EDUCATION))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateEducation(db, &r, &e)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		education := database.GetEducation(db, e.Id)

		if education == nil {
			t.Fatalf("expected %s, received %s", "education", "nil")
		}
		if education.Location != "" {
			t.Fatalf("expected %s, received %s", "", education.Location)
		}
	})

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

		r, err := resume.New([]byte(test.FULL_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		e, err := resume.NewEducation([]byte(test.FULL_EDUCATION))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateEducation(db, &r, &e)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		education := database.GetEducation(db, e.Id)

		if education == nil {
			t.Fatalf("expected %s, received %s", "education", "nil")
		}
		if education.Location != e.Location {
			t.Fatalf("expected %s, received %s", e.Location, education.Location)
		}
	})
}
