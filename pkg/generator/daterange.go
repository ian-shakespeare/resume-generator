package generator

import (
	"fmt"
	"time"
)

type dateRange struct {
	Began    string
	Current  bool
	Finished string
}

func NewDateRange(began time.Time, finished time.Time) dateRange {
	return dateRange{
		Began:    timeToMonthYearStr(began),
		Current:  false,
		Finished: timeToMonthYearStr(finished),
	}
}

func CurrentDateRange(began time.Time) dateRange {
	return dateRange{
		Began:    timeToMonthYearStr(began),
		Current:  true,
		Finished: "",
	}
}

func timeToMonthYearStr(t time.Time) string {
	month := t.Month().String()
	year := t.Year()
	return fmt.Sprintf("%s %d", month, year)
}
