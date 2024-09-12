package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/pkg/resume"
	"resumegenerator/test"
	"testing"
)

func TestCreateWorkExperience(t *testing.T) {
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

		r, err := resume.FromJSON([]byte(test.FULL_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(
			db,
			&user,
			&r,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		w, err := resume.WorkExperienceFromJSON([]byte(test.MIN_WORK_EXPERIENCE))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateWorkExperience(db, &r, &w)
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

		r, err := resume.FromJSON([]byte(test.FULL_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		w, err := resume.WorkExperienceFromJSON([]byte(test.FULL_WORK_EXPERIENCE))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)

		err = database.CreateWorkExperience(db, &r, &w)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})
}

func TestGetWorkExperience(t *testing.T) {
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

		r, err := resume.FromJSON([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		w, err := resume.WorkExperienceFromJSON([]byte(test.MIN_WORK_EXPERIENCE))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)

		err = database.CreateWorkExperience(db, &r, &w)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		workExperience := database.GetWorkExperience(db, w.Id)

		if workExperience == nil {
			t.Fatalf("expected %s, received %s", "workExperience", "nil")
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

		r, err := resume.FromJSON([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		w, err := resume.WorkExperienceFromJSON([]byte(test.MIN_WORK_EXPERIENCE))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)

		err = database.CreateWorkExperience(db, &r, &w)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		workExperience := database.GetWorkExperience(db, w.Id)

		if workExperience == nil {
			t.Fatalf("expected %s, received %s", "workExperience", "nil")
		}
		if workExperience.Location != w.Location {
			t.Fatalf("expected %s, received %s", w.Location, workExperience.Location)
		}
	})
}
