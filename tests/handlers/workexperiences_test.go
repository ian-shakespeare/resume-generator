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
	"resumegenerator/tests"
	"testing"
	"time"
)

func TestHandleCreateWorkExperience(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := tests.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", nil)

		handlers.HandleCreateWorkExperience(w, r, a, db)

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

		handlers.HandleCreateWorkExperience(w, r, a, db)

		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := tests.NewDummyResponseWriter()

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

		handlers.HandleCreateWorkExperience(w, r, a, db)

		if w.StatusCode != 404 {
			t.Fatalf("expected %d, received %d", 404, w.StatusCode)
		}
	})

	t.Run("invalidArgument", func(t *testing.T) {
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

		w := tests.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", nil)
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
		r.SetPathValue("resumeId", resume.Id)

		handlers.HandleCreateWorkExperience(w, r, a, db)

		if w.StatusCode != 400 {
			t.Fatalf("expected %d, received %d", 400, w.StatusCode)
		}

		w.StatusCode = 200

		ne := handlers.NewWorkExperience{
			Employer: "",
		}

		body, err := json.Marshal(ne)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r.Body = io.NopCloser(bytes.NewReader(body))

		handlers.HandleCreateWorkExperience(w, r, a, db)

		if w.StatusCode != 400 {
			t.Fatalf("expected %d, received %d", 400, w.StatusCode)
		}
	})

	t.Run("successful", func(t *testing.T) {
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

		w := tests.NewDummyResponseWriter()

		nr := make([]handlers.NewWorkResponsibility, 0)
		ne := handlers.NewWorkExperience{
			Employer:         "degree",
			Title:            "fieldOfStudy",
			Began:            "1970-01-01T00:00:00.000Z",
			Current:          true,
			Responsibilities: nr,
		}

		body, err := json.Marshal(ne)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := http.NewRequest("POST", "", io.NopCloser(bytes.NewReader(body)))
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
		r.SetPathValue("resumeId", resume.Id)

		handlers.HandleCreateWorkExperience(w, r, a, db)

		if w.StatusCode != 201 {
			t.Fatalf("expected %d, received %d", 201, w.StatusCode)
		}

		var workExperience database.WorkExperience
		if err = json.Unmarshal(w.Body, &workExperience); err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		inserted := database.GetWorkExperience(db, workExperience.Id)
		if inserted == nil {
			t.Fatalf("expected %s, received %s", "workExperience", "nil")
		}
	})
}

func TestHandleGetWorkExperience(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := tests.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", nil)

		handlers.HandleGetWorkExperiences(w, r, a, db)

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

		handlers.HandleGetWorkExperiences(w, r, a, db)

		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

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

		w := tests.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", nil)
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

		handlers.HandleGetWorkExperiences(w, r, a, db)

		if w.StatusCode != 404 {
			t.Fatalf("expected %d, received %d", 404, w.StatusCode)
		}

		w.StatusCode = 200

		r.SetPathValue("resumeId", "BAD")

		handlers.HandleGetWorkExperiences(w, r, a, db)

		if w.StatusCode != 404 {
			t.Fatalf("expected %d, received %d", 404, w.StatusCode)
		}
	})

	t.Run("successful", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := tests.NewDummyResponseWriter()

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

		workExperience, err := database.CreateWorkExperience(
			db,
			&resume,
			"employer",
			"title",
			time.Now(),
			true,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		handlers.HandleGetWorkExperiences(w, r, a, db)

		if w.StatusCode != 200 {
			t.Fatalf("expected %d, received %d", 200, w.StatusCode)
		}

		contentType := w.Header().Get("content-type")
		if contentType != "application/json" {
			t.Fatalf("expected %s, received %s", "application/json", contentType)
		}

		var e []database.WorkExperience
		if err = json.Unmarshal(w.Body, &e); err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(e) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(e))
		}

		if e[0].Id != workExperience.Id {
			t.Fatalf("expected %s, received %s", workExperience.Id, e[0].Id)
		}
	})
}
