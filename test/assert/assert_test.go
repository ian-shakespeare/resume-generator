package assert_test

import (
	"errors"
	"resumegenerator/internal/assert"
	"resumegenerator/test/expect"
	"testing"
)

func TestAssertEqual(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer expect.NoPanic(t)
		assert.Equal(1, 1)
	})

	t.Run("false", func(t *testing.T) {
		defer expect.Panic(t)
		assert.Equal(0, 1)
	})
}

func TestAssertInEqual(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer expect.NoPanic(t)
		assert.InEqual(0, 1)
	})

	t.Run("false", func(t *testing.T) {
		defer expect.Panic(t)
		assert.InEqual(1, 1)
	})
}

func TestAssertGreaterThan(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer expect.NoPanic(t)
		assert.GreaterThan(1, 0)
	})

	t.Run("false", func(t *testing.T) {
		defer expect.Panic(t)
		assert.GreaterThan(1, 1)
	})
}

func TestAssertGreaterThanOrEqual(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer expect.NoPanic(t)
		assert.GreaterThanOrEqual(1, 1)
	})

	t.Run("false", func(t *testing.T) {
		defer expect.Panic(t)
		assert.GreaterThanOrEqual(0, 1)
	})
}

func TestAssertLessThan(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer expect.NoPanic(t)
		assert.LessThan(0, 1)
	})

	t.Run("false", func(t *testing.T) {
		defer expect.Panic(t)
		assert.LessThan(1, 1)
	})
}

func TestAssertLessThanOrEqual(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer expect.NoPanic(t)
		assert.LessThanOrEqual(1, 1)
	})

	t.Run("false", func(t *testing.T) {
		defer expect.Panic(t)
		assert.LessThanOrEqual(1, 0)
	})
}

func TestAssertNilError(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		defer expect.NoPanic(t)
		var err error
		assert.NilError(err)
	})

	t.Run("false", func(t *testing.T) {
		defer expect.Panic(t)
		assert.NilError(errors.New("test"))
	})
}
