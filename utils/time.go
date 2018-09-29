package utils

import "time"

type DateTime struct {
	time time.Time
}

func NewDateTime() *DateTime {
	loc, _ := time.LoadLocation("Asia/Tehran")
	return &DateTime{time: time.Now().In(loc)}
}

func (dt *DateTime) ToW3cString() string {
	return dt.time.Format("2006-01-02T15:04:05-07:00")
}
