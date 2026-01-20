package helpers

import "time"

func GetCurrentTimeStampUTC() time.Time {
	return time.Now().UTC()
}
