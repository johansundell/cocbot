package main

import (
	"fmt"
	"math"
	"time"
)

const (
	leapSince = 37
)

func getMarsTime() string {
	utc := time.Now().UTC()
	millis := utc.Unix()
	jd_ut := 2440587.5 + (float64(millis) / 86400)
	jd_tt := jd_ut + (leapSince+32.184)/86400
	j2 := jd_tt - 2451545
	msd := ((j2 - 4.5) / 1.027491252) + 44796.0 - 0.00096
	mtc := math.Mod((24 * msd), 24)
	return toTime(mtc)
}

func toTime(h float64) string {
	x := h * 3600
	hh := int(math.Floor(x / 3600))
	y := math.Mod(x, 3600)
	mm := int(math.Floor(y / 60))
	out := ""
	if hh < 10 {
		out += fmt.Sprintf("0%d", hh)
	} else {
		out += fmt.Sprintf("%d", hh)
	}
	out += ":"
	if mm < 10 {
		out += fmt.Sprintf("0%d", mm)
	}
	return out
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
}
