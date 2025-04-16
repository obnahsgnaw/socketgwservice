package conflict

import (
	"strings"
	"time"
)

func Gen(t time.Time) string {
	s := t.Format("20060102150405.99")
	return strings.Replace(s, ".", "", 1)
}

func Validate(t time.Time, conflict string) bool {
	return Gen(t) == conflict
}
