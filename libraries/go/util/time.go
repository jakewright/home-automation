package util

import "time"

// TimeToProto returns an RFC3339Nano string
func TimeToProto(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(time.RFC3339Nano)
}

// PTimeToProto returns an RFC3339Nano string
func PTimeToProto(t *time.Time) string {
	if t == nil {
		return ""
	}
	return TimeToProto(*t)
}
