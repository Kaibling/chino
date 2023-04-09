package utils

import "time"

var Layout = "02.01.2006"

func TimeToFormat(t time.Time) string {
	return t.Format(Layout)
}

func TimeFromFormat(s string) (time.Time, error) {
	t, err := time.Parse(Layout, s)
	return t, err
}
