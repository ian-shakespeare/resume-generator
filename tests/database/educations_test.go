package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/tests"
	"testing"
	"time"
)

func TestCreateEducation(t *testing.T) {
	db := tests.SetupDB(t)
	defer tests.TearDownDB(t, db)

	// Arrange
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

	// Act
	_, err = database.CreateEducation(
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

	// Assert
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}
}

func TestCreateEducationNoNil(t *testing.T) {
	db := tests.SetupDB(t)
	defer tests.TearDownDB(t, db)

	// Arrange
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

	// Act
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

	// Assert
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}
}

func TestGetEducation(t *testing.T) {
	db := tests.SetupDB(t)
	defer tests.TearDownDB(t, db)

	// Arrange
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

	// Act
	education := database.GetEducation(db, created.Id)

	// Assert
	if education == nil {
		t.Fatalf("expected %s, received %s", "education", "nil")
	}
	if education.Location != nil {
		t.Fatalf("expected %s, received %p", "nil", education.Location)
	}
}

func TestGetEducationNoNil(t *testing.T) {
	db := tests.SetupDB(t)
	defer tests.TearDownDB(t, db)

	// Arrange
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

	// Act
	education := database.GetEducation(db, created.Id)

	// Assert
	if education == nil {
		t.Fatalf("expected %s, received %s", "education", "nil")
	}
	if *education.Location != location {
		t.Fatalf("expected %s, received %s", location, *education.Location)
	}
}
