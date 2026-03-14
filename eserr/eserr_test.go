package eserr_test

import (
	"errors"
	"fmt"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/eserr"
)

// sentinelError is a named error type used as a test helper.
type sentinelError struct{ msg string }

func (e *sentinelError) Error() string { return e.msg }

func TestUnwrapErr(t *testing.T) {
	t.Parallel()

	t.Run("nil error returns nil", func(t *testing.T) {
		t.Parallel()
		assert.Assert(t, eserr.UnwrapErr(nil) == nil)
	})

	t.Run("non-wrapping error returns nil", func(t *testing.T) {
		t.Parallel()
		assert.Assert(t, eserr.UnwrapErr(errors.New("plain")) == nil)
	})

	t.Run("wrapped error returns inner", func(t *testing.T) {
		t.Parallel()
		inner := errors.New("inner")
		outer := fmt.Errorf("outer: %w", inner)
		assert.Equal(t, inner, eserr.UnwrapErr(outer))
	})
}

func TestFindErrorInChain(t *testing.T) {
	t.Parallel()

	t.Run("nil error returns false", func(t *testing.T) {
		t.Parallel()
		var target *sentinelError
		assert.Assert(t, !eserr.FindErrorInChain[*sentinelError](nil, &target))
		assert.Assert(t, target == nil)
	})

	t.Run("direct match returns true and sets target", func(t *testing.T) {
		t.Parallel()
		se := &sentinelError{msg: "direct"}
		var target *sentinelError
		assert.Assert(t, eserr.FindErrorInChain[*sentinelError](se, &target))
		assert.Equal(t, se, target)
	})

	t.Run("wrapped match returns true and sets target", func(t *testing.T) {
		t.Parallel()
		se := &sentinelError{msg: "wrapped"}
		wrapped := fmt.Errorf("outer: %w", se)
		var target *sentinelError
		assert.Assert(t, eserr.FindErrorInChain[*sentinelError](wrapped, &target))
		assert.Equal(t, se, target)
	})

	t.Run("nil target pointer does not panic", func(t *testing.T) {
		t.Parallel()
		se := &sentinelError{msg: "no target"}
		assert.Assert(t, eserr.FindErrorInChain[*sentinelError](se, nil))
	})

	t.Run("unrelated error returns false", func(t *testing.T) {
		t.Parallel()
		var target *sentinelError
		assert.Assert(t, !eserr.FindErrorInChain[*sentinelError](errors.New("other"), &target))
		assert.Assert(t, target == nil)
	})
}
