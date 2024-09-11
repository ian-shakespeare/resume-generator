package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"resumegenerator/internal/auth"
	"resumegenerator/internal/database"
	"resumegenerator/internal/handlers"
	"resumegenerator/test"
	"testing"
)

func TestHandleCreateProject(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := test.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", nil)

		handlers.HandleCreateProject(w, r, a, db)

		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}

		w.StatusCode = 200

		user1, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		token, err := a.GenToken(&user1)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user2, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&user2,
			"name",
			"email",
			"phoneNumber",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r.SetPathValue("resumeId", resume.Id)
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

		handlers.HandleCreateProject(w, r, a, db)

		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := test.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", nil)

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		token, err := a.GenToken(&user)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r.SetPathValue("resumeId", "random")
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

		handlers.HandleCreateProject(w, r, a, db)

		if w.StatusCode != 404 {
			t.Fatalf("expected %d, received %d", 404, w.StatusCode)
		}
	})

	t.Run("invalidArgument", func(t *testing.T) {
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
			"name",
			"email",
			"phoneNumber",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)
		token, err := a.GenToken(&user)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		w := test.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", nil)
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
		r.SetPathValue("resumeId", resume.Id)

		handlers.HandleCreateProject(w, r, a, db)

		if w.StatusCode != 400 {
			t.Fatalf("expected %d, received %d", 400, w.StatusCode)
		}

		w.StatusCode = 200

		ne := handlers.NewProject{
			Name: "",
		}

		body, err := json.Marshal(ne)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r.Body = io.NopCloser(bytes.NewReader(body))

		handlers.HandleCreateProject(w, r, a, db)

		if w.StatusCode != 400 {
			t.Fatalf("expected %d, received %d", 400, w.StatusCode)
		}
	})

	t.Run("successful", func(t *testing.T) {
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
			"name",
			"email",
			"phoneNumber",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)
		token, err := a.GenToken(&user)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		w := test.NewDummyResponseWriter()

		nr := make([]handlers.NewProjectResponsibility, 0)
		np := handlers.NewProject{
			Name:             "name",
			Description:      "description",
			Role:             "role",
			Responsibilities: nr,
		}

		body, err := json.Marshal(np)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := http.NewRequest("POST", "", io.NopCloser(bytes.NewReader(body)))
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
		r.SetPathValue("resumeId", resume.Id)

		handlers.HandleCreateProject(w, r, a, db)

		if w.StatusCode != 201 {
			t.Fatalf("expected %d, received %d", 201, w.StatusCode)
		}

		var project database.Project
		if err = json.Unmarshal(w.Body, &project); err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		inserted := database.GetProject(db, project.Id)
		if inserted == nil {
			t.Fatalf("expected %s, received %s", "project", "nil")
		}
	})
}

func TestHandleGetProject(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := test.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", nil)

		handlers.HandleGetProjects(w, r, a, db)

		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}

		w.StatusCode = 200

		user1, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		token, err := a.GenToken(&user1)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user2, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&user2,
			"name",
			"email",
			"phoneNumber",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r.SetPathValue("resumeId", resume.Id)
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

		handlers.HandleGetProjects(w, r, a, db)

		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)
		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		token, err := a.GenToken(&user)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		w := test.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", nil)
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

		handlers.HandleGetProjects(w, r, a, db)

		if w.StatusCode != 404 {
			t.Fatalf("expected %d, received %d", 404, w.StatusCode)
		}

		w.StatusCode = 200

		r.SetPathValue("resumeId", "BAD")

		handlers.HandleGetProjects(w, r, a, db)

		if w.StatusCode != 404 {
			t.Fatalf("expected %d, received %d", 404, w.StatusCode)
		}
	})

	t.Run("successful", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := test.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", nil)

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		token, err := a.GenToken(&user)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

		resume, err := database.CreateResume(
			db,
			&user,
			"name",
			"email",
			"phoneNumber",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
		r.SetPathValue("resumeId", resume.Id)

		project, err := database.CreateProject(
			db,
			&resume,
			"name",
			"description",
			"role",
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		handlers.HandleGetProjects(w, r, a, db)

		if w.StatusCode != 200 {
			t.Fatalf("expected %d, received %d", 200, w.StatusCode)
		}

		contentType := w.Header().Get("content-type")
		if contentType != "application/json" {
			t.Fatalf("expected %s, received %s", "application/json", contentType)
		}

		var e []database.Project
		if err = json.Unmarshal(w.Body, &e); err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(e) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(e))
		}

		if e[0].Id != project.Id {
			t.Fatalf("expected %s, received %s", project.Id, e[0].Id)
		}
	})
}
