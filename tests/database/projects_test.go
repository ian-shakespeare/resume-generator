package database_test

import (
	"resumegenerator/internal/database"
	"testing"
)

func TestCreateProject(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

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

	_, err = database.CreateProject(db, &resume, "name", "role")

	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}
}

func TestGetProject(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

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

	created, err := database.CreateProject(db, &resume, "name", "role")

	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	project := database.GetProject(db, created.Id)

	if project == nil {
		t.Fatalf("expected %s, received %s", "project", "nil")
	}
}

func TestProjectResposibility(t *testing.T) {
	db := setup(t)
	defer tearDown(t, db)

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

	project, err := database.CreateProject(db, &resume, "name", "role")
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
}
