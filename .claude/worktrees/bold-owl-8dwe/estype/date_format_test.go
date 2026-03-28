package estype_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

func TestDateFormat_String(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "strict_date", estype.DateFormatStrictDate.String())
	assert.Equal(t, "epoch_millis", estype.DateFormatEpochMillis.String())
}

func TestJoinDateFormats(t *testing.T) {
	t.Parallel()

	t.Run("single format", func(t *testing.T) {
		t.Parallel()
		result := estype.JoinDateFormats(estype.DateFormatStrictDate)
		assert.Equal(t, "strict_date", result)
	})

	t.Run("multiple formats", func(t *testing.T) {
		t.Parallel()
		result := estype.JoinDateFormats(
			estype.DateFormatStrictDateOptionalTime,
			estype.DateFormatEpochMillis,
		)
		assert.Equal(t, "strict_date_optional_time||epoch_millis", result)
	})

	t.Run("no formats", func(t *testing.T) {
		t.Parallel()
		result := estype.JoinDateFormats()
		assert.Equal(t, "", result)
	})
}
