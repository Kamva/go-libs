package utils

import "time"

type DateTime struct {
	timezone *time.Location
	time     time.Time
}

func NewDateTime() *DateTime {
	loc, _ := time.LoadLocation("Asia/Tehran")

	return &DateTime{
		time: time.Now().In(loc),
		timezone: loc,
	}
}

func (dt *DateTime) ToW3cString() string {
	dt.renewTime()

	return dt.time.Format("2006-01-02T15:04:05-07:00")
}

func (dt *DateTime) Timestamp() int64 {
	dt.renewTime()

	return dt.time.Unix()
}

func (dt *DateTime) renewTime() {
	dt.time = time.Now().In(dt.timezone)
}
