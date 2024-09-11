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

const TEST_SIGNING_KEY = "TESTKEY"

func TestHandleCreateResume(t *testing.T) {
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
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		handlers.HandleCreateResume(w, r, a, db)
		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}

		w.StatusCode = 200

		r.Header.Add("authorization", "Bearer")
		handlers.HandleCreateResume(w, r, a, db)
		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}

		w.StatusCode = 200

		r.Header.Set("authorization", "Bearer BAD_TOKEN")
		handlers.HandleCreateResume(w, r, a, db)
		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}
	})

	t.Run("invalidArgument", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := test.NewDummyResponseWriter()

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		token, err := a.GenToken(&user)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		body := []byte("BAD_BODY")

		r, err := http.NewRequest("POST", "", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
		r.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

		handlers.HandleCreateResume(w, r, a, db)
		if w.StatusCode != 400 {
			t.Fatalf("expected %d, received %d", 400, w.StatusCode)
		}

		w.StatusCode = 200

		resume := handlers.NewResume{
			Name: "John Doe",
		}
		body, err = json.Marshal(resume)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		handlers.HandleCreateResume(w, r, a, db)
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

		a := auth.New(TEST_SIGNING_KEY)

		w := test.NewDummyResponseWriter()

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		token, err := a.GenToken(&user)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume := handlers.NewResume{
			Name:        "John Doe",
			Email:       "jdoe@email.com",
			PhoneNumber: "+1 (000) 000-0000",
			Prelude:     "This is a resume",
		}
		body, err := json.Marshal(resume)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := http.NewRequest("POST", "", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
		r.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

		handlers.HandleCreateResume(w, r, a, db)

		if w.StatusCode != 201 {
			t.Fatalf("expected %d, received %d", 201, w.StatusCode)
		}

		contentType := w.Headers.Get("content-type")
		if contentType != "application/json" {
			t.Fatalf("expected %s, received %s", "application/json", contentType)
		}

		var response database.Resume
		err = json.Unmarshal(w.Body, &response)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		stored := database.GetResume(db, response.Id)

		if stored == nil {
			t.Fatalf("expected %s, received %s", "resume", "nil")
		}
	})
}

func TestHandleGetResume(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := test.NewDummyResponseWriter()

		r, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		handlers.HandleGetResume(w, r, a, db)

		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}

		w.StatusCode = 200

		thisUser, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		token, err := a.GenToken(&thisUser)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

		otherUser, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&otherUser,
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

		handlers.HandleGetResume(w, r, a, db)

		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expect %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := test.NewDummyResponseWriter()

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		token, err := a.GenToken(&user)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := http.NewRequest("GET", "", nil)
		r.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))
		r.SetPathValue("resumeId", "NON_EXISTANT")

		handlers.HandleGetResume(w, r, a, db)

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

		token, err := a.GenToken(&user)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := http.NewRequest("GET", "", nil)
		r.SetPathValue("resumeId", resume.Id)
		r.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

		handlers.HandleGetResume(w, r, a, db)

		if w.StatusCode != 200 {
			t.Fatalf("expected %d, received %d", 200, w.StatusCode)
		}

		contentType := w.Headers.Get("content-type")
		if contentType != "application/json" {
			t.Fatalf("expected %s, received %s", "application/json", contentType)
		}

		var response handlers.FullResume
		err = json.Unmarshal(w.Body, &response)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if response.Resume.Id != resume.Id {
			t.Fatalf("expected %s, received %s", resume.Id, response.Resume.Id)
		}
	})
}
