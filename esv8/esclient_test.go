package esv8

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_taskIDToString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		taskID  any
		want    string
		wantErr bool
	}{
		"string":              {taskID: "abc:123", want: "abc:123"},
		"int":                 {taskID: int(42), want: "42"},
		"int8":                {taskID: int8(8), want: "8"},
		"int16":               {taskID: int16(16), want: "16"},
		"int32":               {taskID: int32(32), want: "32"},
		"int64":               {taskID: int64(64), want: "64"},
		"uint":                {taskID: uint(1), want: "1"},
		"uint8":               {taskID: uint8(2), want: "2"},
		"uint16":              {taskID: uint16(3), want: "3"},
		"uint32":              {taskID: uint32(4), want: "4"},
		"uint64":              {taskID: uint64(5), want: "5"},
		"unsupported float64": {taskID: float64(1.5), wantErr: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
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
