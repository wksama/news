package utils

import "strconv"

func Float64ToString(float float64) string {
	return strconv.FormatFloat(float, 'f', -1, 64)
}
