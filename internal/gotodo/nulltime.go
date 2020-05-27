package gotodo

import "time"

// NullTime combines time.Time with a flag to indicate its validity
type NullTime struct {
	Time  time.Time
	Valid bool
}

// InvalidTime is a NullTime with Valid set to false and Time set to time.Time{}
var InvalidTime = NullTime{Valid: false}

// NewNullTime takes a timestamp and returns a NullTime
func NewNullTime(timeStr string) NullTime {
	result := InvalidTime
	if ts, err := time.Parse(TimeFormat, string(timeStr)); err == nil {
		result = ValidTime(ts)
	}

	return result
}

// ValidTime returns a NullTime with Valid set to true and Time set to given time
func ValidTime(time time.Time) NullTime {
	return NullTime{time, true}
}

// Display formats a NullTime in YYYY-MM-DD format
func (me NullTime) Display() string {
	if me.Valid {
		return me.Time.Format(TimeFormat)
	}
	return ""
}
