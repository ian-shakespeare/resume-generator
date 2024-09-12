package test

import (
	"errors"
	"net/http"
	"os"
	"resumegenerator/internal/database"
	"testing"
)

const TEST_DB_NAME string = "testing.db"

func SetupDB(t *testing.T) database.VersionedDatabase {
	file, err := os.Create(TEST_DB_NAME)
	if err != nil {
		t.Fatalf("setup %s", err.Error())
	}
	file.Close()

	db, err := database.NewSQLite(TEST_DB_NAME)
	if err != nil {
		t.Fatalf("setup %s", err.Error())
	}

	return db
}

func TearDownDB(t *testing.T, db database.VersionedDatabase) {
	db.DB().Close()
	err := os.Remove(TEST_DB_NAME)
	if err != nil {
		t.Fatalf("teardown %s", err.Error())
	}
}

type DummyResponseWriter struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

func NewDummyResponseWriter() *DummyResponseWriter {
	return &DummyResponseWriter{
		StatusCode: 200,
		Body:       nil,
		Headers:    http.Header{},
	}
}

func (d *DummyResponseWriter) Header() http.Header {
	return d.Headers
}

func (d *DummyResponseWriter) Write(b []byte) (int, error) {
	if d.Body != nil {
		return 0, errors.New("body has already been written")
	}

	d.Body = b
	return len(b), nil
}

func (d *DummyResponseWriter) WriteHeader(statusCode int) {
	d.StatusCode = statusCode
}

func (d *DummyResponseWriter) Reset() {
	d.StatusCode = 200
	d.Body = nil
	d.Headers = http.Header{}
}
