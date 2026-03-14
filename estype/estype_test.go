package estype_test

import (
	"errors"
	"fmt"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

func TestIndex_String(t *testing.T) {
	t.Parallel()
	idx := estype.Index("my-index")
	assert.Equal(t, "my-index", idx.String())
}

func TestParseESIndex(t *testing.T) {
	t.Parallel()

	t.Run("valid index name", func(t *testing.T) {
		t.Parallel()
		idx, err := estype.ParseESIndex("my-index")
		assert.NilError(t, err)
		assert.Equal(t, estype.Index("my-index"), idx)
	})

	t.Run("empty index name returns error", func(t *testing.T) {
		t.Parallel()
		_, err := estype.ParseESIndex("")
		assert.Assert(t, err != nil)
	})
}

func TestAlias_String(t *testing.T) {
	t.Parallel()
	alias := estype.Alias("my-alias")
	assert.Equal(t, "my-alias", alias.String())
}

func TestParseESAlias(t *testing.T) {
	t.Parallel()

	t.Run("valid alias name", func(t *testing.T) {
		t.Parallel()
		alias, err := estype.ParseESAlias("my-alias")
		assert.NilError(t, err)
		assert.Equal(t, estype.Alias("my-alias"), alias)
	})

	t.Run("empty alias name returns error", func(t *testing.T) {
		t.Parallel()
		_, err := estype.ParseESAlias("")
		assert.Assert(t, err != nil)
	})
}

func TestRefreshInterval_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		interval estype.RefreshInterval
		want     string
	}{
		{
			name:     "disabled",
			interval: estype.RefreshIntervalDisable,
			want:     "-1",
		},
		{
			name:     "not set",
			interval: estype.RefreshIntervalNotSet,
			want:     "",
		},
		{
			name:     "default 1s",
			interval: estype.RefreshIntervalDefault,
			want:     "1s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.interval.String())
		})
	}
}

func TestRefreshInterval_ESTypeDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		interval estype.RefreshInterval
		want     any
	}{
		{
			name:     "disabled returns -1 string",
			interval: estype.RefreshIntervalDisable,
			want:     "-1",
		},
		{
			name:     "not set returns empty string",
			interval: estype.RefreshIntervalNotSet,
			want:     "",
		},
		{
			name:     "default returns 1s string",
			interval: estype.RefreshIntervalDefault,
			want:     "1s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.interval.ESTypeDuration())
		})
	}
}

// sentinelError is a named error type used as a test helper.
type sentinelError struct{ msg string }

func (e *sentinelError) Error() string { return e.msg }

func TestUnwrapErr(t *testing.T) {
	t.Parallel()

	t.Run("nil error returns nil", func(t *testing.T) {
		t.Parallel()
		assert.Assert(t, estype.UnwrapErr(nil) == nil)
	})

	t.Run("non-wrapping error returns nil", func(t *testing.T) {
		t.Parallel()
		assert.Assert(t, estype.UnwrapErr(errors.New("plain")) == nil)
	})

	t.Run("wrapped error returns inner", func(t *testing.T) {
		t.Parallel()
		inner := errors.New("inner")
		outer := fmt.Errorf("outer: %w", inner)
		assert.Equal(t, inner, estype.UnwrapErr(outer))
	})
}

func TestFindErrorInChain(t *testing.T) {
	t.Parallel()

	t.Run("nil error returns false", func(t *testing.T) {
		t.Parallel()
		var target *sentinelError
		assert.Assert(t, !estype.FindErrorInChain[*sentinelError](nil, &target))
		assert.Assert(t, target == nil)
	})

	t.Run("direct match returns true and sets target", func(t *testing.T) {
		t.Parallel()
		se := &sentinelError{msg: "direct"}
		var target *sentinelError
		assert.Assert(t, estype.FindErrorInChain[*sentinelError](se, &target))
		assert.Equal(t, se, target)
	})

	t.Run("wrapped match returns true and sets target", func(t *testing.T) {
		t.Parallel()
		se := &sentinelError{msg: "wrapped"}
		wrapped := fmt.Errorf("outer: %w", se)
		var target *sentinelError
		assert.Assert(t, estype.FindErrorInChain[*sentinelError](wrapped, &target))
		assert.Equal(t, se, target)
	})

	t.Run("nil target pointer does not panic", func(t *testing.T) {
		t.Parallel()
		se := &sentinelError{msg: "no target"}
		assert.Assert(t, estype.FindErrorInChain[*sentinelError](se, nil))
	})

	t.Run("unrelated error returns false", func(t *testing.T) {
		t.Parallel()
		var target *sentinelError
		assert.Assert(t, !estype.FindErrorInChain[*sentinelError](errors.New("other"), &target))
		assert.Assert(t, target == nil)
	})
}

func TestParseRefreshInterval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    estype.RefreshInterval
		wantErr bool
	}{
		{
			name:  "disabled",
			input: "-1",
			want:  estype.RefreshIntervalDisable,
		},
		{
			name:  "not set",
			input: "",
			want:  estype.RefreshIntervalNotSet,
		},
		{
			name:  "1 second",
			input: "1s",
			want:  estype.RefreshIntervalDefault,
		},
		{
			name:    "invalid",
			input:   "not-a-duration",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := estype.ParseRefreshInterval(tt.input)
			if tt.wantErr {
				assert.Assert(t, err != nil)
				return
			}
			assert.NilError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
