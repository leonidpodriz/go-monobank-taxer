package main

import "time"

func ReversedTimePeriods(from, to time.Time, step time.Duration) []TimePeriod {
	var periods []TimePeriod
	var cursor = to

	if to.Before(from) {
		return periods
	}

	for {
		period := TimePeriod{
			From: cursor.Add(-step),
			To:   cursor,
		}

		if period.From.Before(from) {
			period.From = from
		}

		periods = append(periods, period)

		if period.From.Equal(from) {
			break
		}

		cursor = period.From
	}

	return periods
}
