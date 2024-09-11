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

func TestHandleCreateEducation(t *testing.T) {
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

		handlers.HandleCreateEducation(w, r, a, db)

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

		handlers.HandleCreateEducation(w, r, a, db)

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

		handlers.HandleCreateEducation(w, r, a, db)

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

		handlers.HandleCreateEducation(w, r, a, db)

		if w.StatusCode != 400 {
			t.Fatalf("expected %d, received %d", 400, w.StatusCode)
		}

		w.StatusCode = 200

		ne := handlers.NewEducation{
			DegreeType: "",
		}

		body, err := json.Marshal(ne)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r.Body = io.NopCloser(bytes.NewReader(body))

		handlers.HandleCreateEducation(w, r, a, db)

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

		ne := handlers.NewEducation{
			DegreeType:   "degree",
			FieldOfStudy: "fieldOfStudy",
			Institution:  "institution",
			Began:        "1970-01-01T00:00:00.000Z",
			Current:      true,
		}

		body, err := json.Marshal(ne)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := http.NewRequest("POST", "", io.NopCloser(bytes.NewReader(body)))
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
		r.SetPathValue("resumeId", resume.Id)

		handlers.HandleCreateEducation(w, r, a, db)

		if w.StatusCode != 201 {
			t.Fatalf("expected %d, received %d", 201, w.StatusCode)
		}

		var education database.Education
		if err = json.Unmarshal(w.Body, &education); err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		inserted := database.GetEducation(db, education.Id)
		if inserted == nil {
			t.Fatalf("expected %s, received %s", "education", "nil")
		}
	})
}

func TestHandleGetEducation(t *testing.T) {
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

		handlers.HandleGetEducations(w, r, a, db)

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

		handlers.HandleGetEducations(w, r, a, db)

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

		handlers.HandleGetEducations(w, r, a, db)

		if w.StatusCode != 404 {
			t.Fatalf("expected %d, received %d", 404, w.StatusCode)
		}

		w.StatusCode = 200

		r.SetPathValue("resumeId", "BAD")

		handlers.HandleGetEducations(w, r, a, db)

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

		education, err := database.CreateEducation(
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

		handlers.HandleGetEducations(w, r, a, db)

		if w.StatusCode != 200 {
			t.Fatalf("expected %d, received %d", 200, w.StatusCode)
		}

		contentType := w.Header().Get("content-type")
		if contentType != "application/json" {
			t.Fatalf("expected %s, received %s", "application/json", contentType)
		}

		var e []database.Education
		if err = json.Unmarshal(w.Body, &e); err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(e) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(e))
		}

		if e[0].Id != education.Id {
			t.Fatalf("expected %s, received %s", education.Id, e[0].Id)
		}
	})
}
