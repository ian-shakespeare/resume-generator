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
	"resumegenerator/pkg/resume"
	"resumegenerator/test"
	"testing"
)

const TEST_SIGNING_KEY = "TESTKEY"

type newResume struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	Prelude     string  `json:"prelude"`
	Location    *string `json:"location"`
	LinkedIn    *string `json:"linkedIn"`
	Github      *string `json:"github"`
	Facebook    *string `json:"facebook"`
	Instagram   *string `json:"instagram"`
	Twitter     *string `json:"twitter"`
	Portfolio   *string `json:"portfolio"`
}

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

		w.Reset()

		r.Header.Add("authorization", "Bearer")
		handlers.HandleCreateResume(w, r, a, db)
		if w.StatusCode != 401 {
			t.Fatalf("expected %d, received %d", 401, w.StatusCode)
		}

		w.Reset()

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

		w.Reset()

		resume := newResume{
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

		r, err := http.NewRequest("POST", "", bytes.NewReader([]byte(test.MIN_RESUME)))
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

		var response resume.Resume
		err = json.Unmarshal(w.Body, &response)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		stored := database.GetResume(db, response.Id, user.Id)

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

		w.Reset()

		r.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

		otherUser, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		res, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &otherUser, &res)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r.SetPathValue("resumeId", res.Id)

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

		res, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &res)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		token, err := a.GenToken(&user)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := http.NewRequest("GET", "", nil)
		r.SetPathValue("resumeId", res.Id)
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

		if response.Resume.Id != res.Id {
			t.Fatalf("expected %s, received %s", res.Id, response.Resume.Id)
		}
	})
}