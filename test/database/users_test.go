package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/pkg/resume"
	"resumegenerator/test"
	"testing"
)

func TestCreateUser(t *testing.T) {
	conn := test.SetupDB(t)
	defer test.TearDownDB(t, conn)

	err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	user, err := database.CreateUser(conn)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	rows, err := conn.Database().Query("SELECT 1 FROM users WHERE user_id = ?", user.Id)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	if !rows.Next() {
		t.Fatalf("expected %s, received %s", "true", "false")
	}

	if user.Id == "" {
		t.Fatalf("expected %s, received %s", "id", "\"\"")
	}
}

func TestGetUser(t *testing.T) {
	t.Run("nonexistant", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.GetUser(conn, "SOME_ID")
		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}

		if user.Id != "" {
			t.Fatalf("expected %s, received %s", "\"\"", user.Id)
		}
	})

	t.Run("existant", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		created, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.GetUser(conn, created.Id)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if created.Id != user.Id {
			t.Fatalf("expected %s, received %s", created.Id, user.Id)
		}
	})
}

func TestAddResume(t *testing.T) {
	t.Run("minimumRequiredFields", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		rows, err := conn.Database().Query("SELECT 1 FROM resumes WHERE resume_id = ?", r.Id)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if !rows.Next() {
			t.Fatalf("expected %s, received %s", "true", "false")
		}

		if len(u.Resumes) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(u.Resumes))
		}
	})

	t.Run("allFields", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.Example())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		rows, err := conn.Database().Query("SELECT 1 FROM resumes WHERE resume_id = ?", r.Id)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if !rows.Next() {
			t.Fatalf("expected %s, received %s", "true", "false")
		}

		if len(u.Resumes) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(u.Resumes))
		}
	})
}
