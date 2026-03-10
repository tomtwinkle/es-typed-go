package estype

import (
	"fmt"
	"time"
)

// RefreshInterval represents an Elasticsearch refresh interval duration.
// Special values are provided for common cases.
type RefreshInterval time.Duration

const (
	// RefreshIntervalNotSet means the refresh interval is not explicitly configured.
	// Callers should interpret this as "use default".
	RefreshIntervalNotSet RefreshInterval = 0

	// RefreshIntervalDisable disables automatic refreshing (-1).
	RefreshIntervalDisable RefreshInterval = -1

	// RefreshIntervalDefault is the Elasticsearch default refresh interval (1s).
	RefreshIntervalDefault RefreshInterval = RefreshInterval(time.Second)
)

// String returns the string representation of the RefreshInterval suitable for Elasticsearch.
// A value of -1 means disabled, 0 means not set, otherwise it formats as a duration string.
func (r RefreshInterval) String() string {
	switch r {
	case RefreshIntervalDisable:
		return "-1"
	case RefreshIntervalNotSet:
		return ""
	default:
		return time.Duration(r).String()
	}
}

// ESTypeDuration returns the value as a types.Duration (any) for use with the Elasticsearch typed API.
func (r RefreshInterval) ESTypeDuration() any {
	switch r {
	case RefreshIntervalDisable:
		return "-1"
	case RefreshIntervalNotSet:
		return ""
	default:
		return time.Duration(r).String()
	}
}

// ParseRefreshInterval parses a refresh interval string from Elasticsearch.
// Accepts "-1" for disabled, "" for not set, or a valid Go duration string (e.g. "1s", "500ms").
func ParseRefreshInterval(s string) (RefreshInterval, error) {
	switch s {
	case "-1":
		return RefreshIntervalDisable, nil
	case "":
		return RefreshIntervalNotSet, nil
	default:
		d, err := time.ParseDuration(s)
		if err != nil {
			return 0, fmt.Errorf("failed to parse refresh interval %q: %w", s, err)
		}
		return RefreshInterval(d), nil
	}
}
