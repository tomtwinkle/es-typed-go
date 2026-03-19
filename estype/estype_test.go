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
