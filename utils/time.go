package utils

import (
	"gorm.io/datatypes"
	"time"
)

func Str2Date(dateStr string) datatypes.Date {
	localTime, _ := time.ParseInLocation("20060102", dateStr, time.Local)
	return datatypes.Date(localTime)
}
