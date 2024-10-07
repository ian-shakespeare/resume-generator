package expect

import "testing"

func NilError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("expected nil error, received %s", err.Error())
	}
}

func Error(t *testing.T, err error) {
	if err == nil {
		t.Fatal("expected error, received nil")
	}
}

func Equal[T comparable](t *testing.T, left T, right T) {
	if left != right {
		t.Fatalf("expected %v == %v", left, right)
	}
}

func Nil[T any](t *testing.T, p *T) {
	if p != nil {
		t.Fatalf("expected nil, received %p", p)
	}
}

func True(t *testing.T, expr bool) {
	if !expr {
		t.Fatal("expected true, received false")
	}
}

func Panic(t *testing.T) {
	if r := recover(); r == nil {
		t.Fatal("expected panic")
	}
}

func NoPanic(t *testing.T) {
	if r := recover(); r != nil {
		t.Fatal("received unexpected panic")
	}
}
