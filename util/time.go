package util

import "time"

func FormatTime(ts time.Time) string {
	return ts.Format("2006-01-02T15:04:05Z")
	//return ts.Format(time.RFC3339)
}

func FormatUtcNow() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}
