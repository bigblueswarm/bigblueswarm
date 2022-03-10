package utils

import "math"

// Round2Digits round a float64 with 2 decimals
func Round2Digits(number float64) float64 {
	return math.Round(number*100) / 100
}
