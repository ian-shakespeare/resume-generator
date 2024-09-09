package database_test

import (
	"resumegenerator/internal/database"
	"testing"
)

func TestCreateResume(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations(), 1)

	user, err := database.CreateUser(db)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	// Act
	_, err = database.CreateResume(
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

	// Assert
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}
}

func TestCreateResumeNoNil(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())

	user, err := database.CreateUser(db)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	location := "SOME_PLACE"
	linkedIn := "SOME_LINKEDIN"
	github := "SOME_GITHUB"
	facebook := "SOME_FACEBOOK"
	instagram := "SOME_INSTAGRAM"
	twitter := "SOME_TWITTER"
	portfolio := "SOME_PORTFOLIO"

	// Act
	_, err = database.CreateResume(
		db,
		&user,
		"John Doe",
		"jdoe@email.com",
		"+1 (000) 000-0000",
		"prelude",
		&location,
		&linkedIn,
		&github,
		&facebook,
		&instagram,
		&twitter,
		&portfolio,
	)

	// Assert
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}
}

func TestGetResume(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())

	user, err := database.CreateUser(db)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	created, err := database.CreateResume(
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
	resume := database.GetResume(db, created.Id)

	// Assert
	if resume == nil {
		t.Fatalf("expected %s, received %s", "resume", "nil")
	}
	if resume.Location != nil {
		t.Fatalf("expected %s, received %p", "nil", resume.Location)
	}
}

func TestGetResumeNoNil(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

	// Arrange
	database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())

	user, err := database.CreateUser(db)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	location := "SOME_PLACE"
	linkedIn := "SOME_LINKEDIN"
	github := "SOME_GITHUB"
	facebook := "SOME_FACEBOOK"
	instagram := "SOME_INSTAGRAM"
	twitter := "SOME_TWITTER"
	portfolio := "SOME_PORTFOLIO"

	created, err := database.CreateResume(
		db,
		&user,
		"John Doe",
		"jdoe@email.com",
		"+1 (000) 000-0000",
		"prelude",
		&location,
		&linkedIn,
		&github,
		&facebook,
		&instagram,
		&twitter,
		&portfolio,
	)

	// Act
	resume := database.GetResume(db, created.Id)

	// Assert
	if resume == nil {
		t.Fatalf("expected %s, received %s", "resume", "nil")
	}
	if *resume.Location != location {
		t.Fatalf("expected %s, received %s", location, *resume.Location)
	}
}
