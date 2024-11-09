package util

import "time"

func GetTime(args ...time.Time) float64 {
	t := time.Now()
	if len(args) > 0 {
		t = args[0]
	}

	return float64(t.UnixMicro()) / 1000000
}
