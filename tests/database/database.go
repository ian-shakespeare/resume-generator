package database_test

import (
	"os"
	"testing"
)

const TEST_DB_NAME string = "testing.db"

func setup(t *testing.T) {
	file, err := os.Create(TEST_DB_NAME)
	if err != nil {
		t.Fatalf("setup %s", err.Error())
	}
	file.Close()
}

func tearDown(t *testing.T) {
	err := os.Remove(TEST_DB_NAME)
	if err != nil {
		t.Fatalf("teardown %s", err.Error())
	}
}
