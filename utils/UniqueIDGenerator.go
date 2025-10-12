package utils

import (
	"time"
)

// GenerateID generates a unique ID based on the current timestamp in format yyMMddHHmmss
// Returns an integer ID like 250110123045 for 2025-01-10 12:30:45
func GenerateID() int64 {
	now := time.Now()
	// Format: yyMMddHHmmss
	id := int64(now.Year()%100)*10000000000 +
		int64(now.Month())*100000000 +
		int64(now.Day())*1000000 +
		int64(now.Hour())*10000 +
		int64(now.Minute())*100 +
		int64(now.Second())
	return id
}
