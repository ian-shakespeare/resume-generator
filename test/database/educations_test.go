package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/test"
	"testing"
	"time"
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

		resume, err := database.CreateResume(
			db,
			&user,
			"John Doe",
			"jdoe@email.com",
			"+1 (000) 000-0000",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		_, err = database.CreateEducation(
			db,
			&resume,
			"degree",
			"fieldOfStudy",
			"institution",
			time.Now(),
			true,
			nil,
			nil,
			nil,
		)

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
		resume, err := database.CreateResume(
			db,
			&user,
			"John Doe",
			"jdoe@email.com",
			"+1 (000) 000-0000",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		location := "SOME_PLACE"
		finished := time.Now()
		gpa := "SOME_GPA"

		_, err = database.CreateEducation(
			db,
			&resume,
			"BS",
			"Computer Science",
			"Utah Tech University",
			time.Now(),
			true,
			&location,
			&finished,
			&gpa,
		)

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

		resume, err := database.CreateResume(
			db,
			&user,
			"John Doe",
			"jdoe@email.com",
			"+1 (000) 000-0000",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		created, err := database.CreateEducation(
			db,
			&resume,
			"BS",
			"Computer Science",
			"Utah Tech University",
			time.Now(),
			true,
			nil,
			nil,
			nil,
		)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		education := database.GetEducation(db, created.Id)

		if education == nil {
			t.Fatalf("expected %s, received %s", "education", "nil")
		}
		if education.Location != nil {
			t.Fatalf("expected %s, received %p", "nil", education.Location)
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

		resume, err := database.CreateResume(
			db,
			&user,
			"John Doe",
			"jdoe@email.com",
			"+1 (000) 000-0000",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		location := "SOME_PLACE"
		finished := time.Now()
		gpa := "SOME_GPA"

		created, err := database.CreateEducation(
			db,
			&resume,
			"BS",
			"Computer Science",
			"Utah Tech University",
			time.Now(),
			true,
			&location,
			&finished,
			&gpa,
		)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		education := database.GetEducation(db, created.Id)

		if education == nil {
			t.Fatalf("expected %s, received %s", "education", "nil")
		}
		if *education.Location != location {
			t.Fatalf("expected %s, received %s", location, *education.Location)
		}
	})
}
