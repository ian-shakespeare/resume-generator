package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/tests"
	"testing"
	"time"
)

func TestCreateWorkExperience(t *testing.T) {
	t.Run("minimumRequiredFields", func(t *testing.T) {
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

		_, err = database.CreateWorkExperience(db, &resume, "employer", "title", time.Now(), true, nil, nil)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})

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

		location := "location"
		finished := time.Now()
		_, err = database.CreateWorkExperience(db, &resume, "employer", "title", time.Now(), true, &location, &finished)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})
}

func TestGetWorkExperience(t *testing.T) {
	t.Run("minimumRequiredFields", func(t *testing.T) {
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

		created, err := database.CreateWorkExperience(db, &resume, "employer", "title", time.Now(), true, nil, nil)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		workExperience := database.GetWorkExperience(db, created.Id)

		if workExperience == nil {
			t.Fatalf("expected %s, received %s", "workExperience", "nil")
		}
	})

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

		location := "location"
		finished := time.Now()
		created, err := database.CreateWorkExperience(db, &resume, "employer", "title", time.Now(), true, &location, &finished)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		workExperience := database.GetWorkExperience(db, created.Id)

		if workExperience == nil {
			t.Fatalf("expected %s, received %s", "workExperience", "nil")
		}
		if *workExperience.Location != location {
			t.Fatalf("expected %s, received %s", location, *workExperience.Location)
		}
	})
}

func TestCreateWorkResposibility(t *testing.T) {
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

		workExperience, err := database.CreateWorkExperience(db, &resume, "employer", "title", time.Now(), true, nil, nil)
		if err != nil {
			t.Fatalf("expected %s, received %s", "workExperience", "nil")
		}

		responsibility, err := database.CreateWorkResponsibility(db, &workExperience, "responsibility")

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(workExperience.Responsibilities) < 1 {
			t.Fatalf("expected %d, received %d", 1, len(workExperience.Responsibilities))
		} else if workExperience.Responsibilities[0].Id != responsibility.Id {
			t.Fatalf("expected %s, received %s", responsibility.Id, workExperience.Responsibilities[0].Id)
		}
	})
}
