package utils

import (
	"gorm.io/datatypes"
	"time"
)

func Time2Str(time time.Time) string {
	return time.Format("20060102")
}

func Date2Str(date datatypes.Date) string {
	dateTime := time.Time(date)
	return Time2Str(dateTime)
}

func Str2Time(dateStr string) time.Time {
	localTime, _ := time.ParseInLocation("20060102", dateStr, time.Local)
	return localTime
}

func Str2Date(dateStr string) datatypes.Date {
	localTime, _ := time.ParseInLocation("20060102", dateStr, time.Local)
	return datatypes.Date(localTime)
}
