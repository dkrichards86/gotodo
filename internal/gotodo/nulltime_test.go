package gotodo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidTime(t *testing.T) {
	now := time.Now()
	nullTime := ValidTime(now)
	assert.Equal(t, now, nullTime.Time)
	assert.Equal(t, true, nullTime.Valid)
}

func TestInvalidTime(t *testing.T) {
	nullTime := InvalidTime
	assert.Equal(t, time.Time{}, nullTime.Time)
	assert.Equal(t, false, nullTime.Valid)
}

func TestValidDisplay(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	t2 := time.Time{}
	assert.Equal(t, "2012-11-01", ValidTime(t1).Display())
	assert.Equal(t, "0001-01-01", ValidTime(t2).Display())
}
func TestInvalidDisplay(t *testing.T) {
	assert.Equal(t, "", InvalidTime.Display())
}
