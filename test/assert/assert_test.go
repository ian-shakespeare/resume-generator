package assert_test

import (
	"errors"
	"resumegenerator/internal/assert"
	"testing"
)

func TestAssertEqual(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("expected success, received panic")
			}
		}()

		assert.Equal(1, 1)
	})

	t.Run("false", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, received success")
			}
		}()

		assert.Equal(0, 1)
	})
}

func TestAssertInEqual(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("expected success, received panic")
			}
		}()

		assert.InEqual(0, 1)
	})

	t.Run("false", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, received success")
			}
		}()

		assert.InEqual(1, 1)
	})
}

func TestAssertGreaterThan(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("expected success, received panic")
			}
		}()

		assert.GreaterThan(1, 0)
	})

	t.Run("false", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, received success")
			}
		}()

		assert.GreaterThan(1, 1)
	})
}

func TestAssertGreaterThanOrEqual(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("expected success, received panic")
			}
		}()

		assert.GreaterThanOrEqual(1, 1)
	})

	t.Run("false", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, received success")
			}
		}()

		assert.GreaterThanOrEqual(0, 1)
	})
}

func TestAssertLessThan(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("expected success, received panic")
			}
		}()

		assert.LessThan(0, 1)
	})

	t.Run("false", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, received success")
			}
		}()

		assert.LessThan(1, 1)
	})
}

func TestAssertLessThanOrEqual(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("expected success, received panic")
			}
		}()

		assert.LessThanOrEqual(1, 1)
	})

	t.Run("false", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, received success")
			}
		}()

		assert.LessThanOrEqual(1, 0)
	})
}

func TestAssertNilError(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("expected success, received panic")
			}
		}()

		var err error
		assert.NilError(err)
	})

	t.Run("false", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, received success")
			}
		}()

		assert.NilError(errors.New("test"))
	})
}
