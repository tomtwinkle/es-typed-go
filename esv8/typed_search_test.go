package esv8

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestSearchParams_ToRequest_DefaultSizeUsesElasticsearchDefault(t *testing.T) {
	t.Parallel()

	req := (SearchParams{}).ToRequest()

	assert.Assert(t, req != nil)
	assert.Assert(t, req.Size == nil)
	assert.Equal(t, req.Source_, true)
	assert.Assert(t, req.Timeout != nil)
	assert.Equal(t, *req.Timeout, "10s")
}

func TestSearchParams_ToRequest_ExplicitSizeIsApplied(t *testing.T) {
	t.Parallel()

	req := (SearchParams{
		Size: 25,
	}).ToRequest()

	assert.Assert(t, req != nil)
	assert.Assert(t, req.Size != nil)
	assert.Equal(t, *req.Size, 25)
	assert.Equal(t, req.Source_, true)
	assert.Assert(t, req.Timeout != nil)
	assert.Equal(t, *req.Timeout, "10s")
}
