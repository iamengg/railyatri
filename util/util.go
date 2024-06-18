package util

import (
	"fmt"
	"time"
)

// Returns date by adding afterdays to today's date
func GetDate(afterDays int) string {
	year, mon, day := time.Now().AddDate(0, 0, afterDays).Date() //time.Now().Date()
	date := fmt.Sprintf("%v-%v-%v", year, int(mon), day)
	return date
}
