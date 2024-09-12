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

func TestHandleCreateEducation(t *testing.T) {
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

		handlers.HandleCreateEducation(w, r, a, db)

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

		handlers.HandleCreateEducation(w, r, a, db)

		if w.StatusCode != 404 {
			t.Fatalf("expected %d, received %d", 404, w.StatusCode)
		}

		w.Reset()

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

		handlers.HandleCreateEducation(w, r, a, db)

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

		res, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &res)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

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
		r.SetPathValue("resumeId", res.Id)

		handlers.HandleCreateEducation(w, r, a, db)

		if w.StatusCode != 400 {
			t.Fatalf("expected %d, received %d", 400, w.StatusCode)
		}

		w.Reset()

		r.Body = io.NopCloser(bytes.NewReader([]byte(`{"degreeType": "degreeType"}`)))

		handlers.HandleCreateEducation(w, r, a, db)

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

		res, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &res)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)
		token, err := a.GenToken(&user)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		w := test.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", io.NopCloser(bytes.NewReader([]byte(test.MIN_EDUCATION))))
		r.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
		r.SetPathValue("resumeId", res.Id)

		handlers.HandleCreateEducation(w, r, a, db)

		if w.StatusCode != 201 {
			t.Fatalf("expected %d, received %d", 201, w.StatusCode)
		}

		var education resume.Education
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
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		a := auth.New(TEST_SIGNING_KEY)

		w := test.NewDummyResponseWriter()

		r, err := http.NewRequest("POST", "", nil)

		handlers.HandleGetEducations(w, r, a, db)

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

		handlers.HandleGetEducations(w, r, a, db)

		if w.StatusCode != 404 {
			t.Fatalf("expected %d, received %d", 404, w.StatusCode)
		}

		w.Reset()

		r.SetPathValue("resumeId", "BAD")

		handlers.HandleGetEducations(w, r, a, db)

		if w.StatusCode != 404 {
			t.Fatalf("expected %d, received %d", 404, w.StatusCode)
		}

		w.Reset()

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

		handlers.HandleGetEducations(w, r, a, db)

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

		res, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &res)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
		r.SetPathValue("resumeId", res.Id)

		e, err := resume.NewEducation([]byte(test.MIN_EDUCATION))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateEducation(db, &res, &e)
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

		var educations []resume.Education
		if err = json.Unmarshal(w.Body, &educations); err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(educations) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(educations))
		}

		if educations[0].Id != e.Id {
			t.Fatalf("expected %s, received %s", e.Id, educations[0].Id)
		}
	})
}
