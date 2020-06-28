package repository

import (
	"encoding/base64"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

// DecodeCursor func
func DecodeCursor(encodedTime string) (time.Time, error) {
	bytecode, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeStr := string(bytecode)
	t, err := time.Parse(timeFormat, timeStr)
	return t, err
}

// EncodeCursor func
func EncodeCursor(t time.Time) string {
	timeStr := t.Format(timeFormat)
	return base64.StdEncoding.EncodeToString([]byte(timeStr))
}
