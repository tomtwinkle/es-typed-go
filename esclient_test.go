package estypedgo

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_taskIDToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		taskID  any
		want    string
		wantErr bool
	}{
		{name: "string", taskID: "abc:123", want: "abc:123"},
		{name: "int", taskID: int(42), want: "42"},
		{name: "int8", taskID: int8(8), want: "8"},
		{name: "int16", taskID: int16(16), want: "16"},
		{name: "int32", taskID: int32(32), want: "32"},
		{name: "int64", taskID: int64(64), want: "64"},
		{name: "uint", taskID: uint(1), want: "1"},
		{name: "uint8", taskID: uint8(2), want: "2"},
		{name: "uint16", taskID: uint16(3), want: "3"},
		{name: "uint32", taskID: uint32(4), want: "4"},
		{name: "uint64", taskID: uint64(5), want: "5"},
		{name: "unsupported float64", taskID: float64(1.5), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := taskIDToString(tt.taskID)
			if tt.wantErr {
				assert.Assert(t, err != nil)
				return
			}
			assert.NilError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_isElasticsearchError_nil(t *testing.T) {
	t.Parallel()
	assert.Assert(t, !isElasticsearchError(nil, nil))
}
