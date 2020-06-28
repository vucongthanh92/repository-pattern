package repository

import (
	"encoding/base64"
	"time"

	"github.com/spf13/viper"
)

// DecodeCursor func
func DecodeCursor(encodedTime string) (time.Time, error) {
	bytecode, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeStr := string(bytecode)
	t, err := time.Parse(viper.GetString("timeformat"), timeStr)
	return t, err
}

// EncodeCursor func
func EncodeCursor(t time.Time) string {
	timeStr := t.Format(viper.GetString("timeformat"))
	return base64.StdEncoding.EncodeToString([]byte(timeStr))
}
