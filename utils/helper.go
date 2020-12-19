package utils

import (
	"gorm.io/datatypes"
	"strconv"
	"time"
)

func DateToFloat64(date datatypes.Date) float64 {
	float, err := strconv.ParseFloat(time.Time(date).Format("20060102"), 64)
	if err != nil {
		return 0
	}
	return float
}
