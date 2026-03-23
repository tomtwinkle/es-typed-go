package estype_test

import (
	"testing"

	"github.com/tomtwinkle/es-typed-go/estype"
	"gotest.tools/v3/assert"
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

func TestFieldNames(t *testing.T) {
	t.Parallel()
	names := estype.FieldNames(estype.Field("title"), estype.Field("name"), estype.Field("status"))
	assert.Assert(t, len(names) == 3)
	assert.Equal(t, "title", names[0])
	assert.Equal(t, "name", names[1])
	assert.Equal(t, "status", names[2])
}

func TestFieldNames_Empty(t *testing.T) {
	t.Parallel()
	names := estype.FieldNames()
	assert.Assert(t, len(names) == 0)
}

func TestFieldNames_Single(t *testing.T) {
	t.Parallel()
	names := estype.FieldNames(estype.Field("status"))
	assert.Assert(t, len(names) == 1)
	assert.Equal(t, "status", names[0])
}

func TestRefreshInterval_String(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		interval estype.RefreshInterval
		want     string
	}{
		"disabled": {
			interval: estype.RefreshIntervalDisable,
			want:     "-1",
		},
		"not set": {
			interval: estype.RefreshIntervalNotSet,
			want:     "",
		},
		"default 1s": {
			interval: estype.RefreshIntervalDefault,
			want:     "1s",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.interval.String())
		})
	}
}

func TestRefreshInterval_ESTypeDuration(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		interval estype.RefreshInterval
		want     any
	}{
		"disabled returns -1 string": {
			interval: estype.RefreshIntervalDisable,
			want:     "-1",
		},
		"not set returns empty string": {
			interval: estype.RefreshIntervalNotSet,
			want:     "",
		},
		"default returns 1s string": {
			interval: estype.RefreshIntervalDefault,
			want:     "1s",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.interval.ESTypeDuration())
		})
	}
}

func TestParseRefreshInterval(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		want    estype.RefreshInterval
		wantErr bool
	}{
		"disabled": {
			input: "-1",
			want:  estype.RefreshIntervalDisable,
		},
		"not set": {
			input: "",
			want:  estype.RefreshIntervalNotSet,
		},
		"1 second": {
			input: "1s",
			want:  estype.RefreshIntervalDefault,
		},
		"invalid": {
			input:   "not-a-duration",
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
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
