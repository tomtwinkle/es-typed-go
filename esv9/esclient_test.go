package esv9

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_isElasticsearchError_nil(t *testing.T) {
	t.Parallel()
	assert.Assert(t, !isElasticsearchError(nil, nil))
}
