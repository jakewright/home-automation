package domain

import (
	"math"
	"time"
)

// Rule represents a periodic time
type Rule struct {
	Second  *int // 0-59
	Minute  *int // 0-59
	Hour    *int // 0-23
	Weekday *int // 0-6 (Sunday = 0)
	Day     *int // 1-31
	Month   *int // 1-12
}

// CalculateNextRunAfterTime returns the first time that the rule describes that is later than the given time t
func (r *Rule) CalculateNextRunAfterTime(t time.Time) (time.Time, error) {
	// Make sure the time is advanced by at least one second
	n := t.Add(time.Second)

	var stage float64
	var d time.Duration
	done := false

	for i := 0; i < 400 && !done; i++ {
		// If a second is set
		if r.Second != nil {
			d = mod(*r.Second-n.Second(), 60)
			n = n.Add(d * time.Second)
		} else if stage > 0 {
			n = n.Add(time.Duration(n.Second()) * time.Second * -1)
		}

		stage = math.Max(stage, 1)

		// If a minute is set
		if r.Minute != nil {
			d = mod(*r.Minute-n.Minute(), 60)
			n = n.Add(d * time.Minute)
			if d > 0 {
				continue
			}
		} else if stage > 1 {
			n = n.Add(time.Duration(n.Minute()) * time.Minute * -1)
		}

		stage = math.Max(stage, 2)

		// If an hour is set
		if r.Hour != nil {
			d = mod(*r.Hour-n.Hour(), 24)
			n = n.Add(d * time.Hour)
			if d > 0 {
				continue
			}
		} else if stage > 2 {
			n = n.Add(time.Duration(n.Hour()) * time.Hour * -1)
		}

		stage = math.Max(stage, 3)

		// If a weekday is set
		if r.Weekday != nil {
			d = mod(*r.Weekday-int(n.Weekday()), 7)
			n = n.AddDate(0, 0, int(d))
			if d > 0 {
				continue
			}
		}
		// There's no need to reset the day here because the month
		// adding section will only increment to the first day of a month.

		// If a day is set
		if r.Day != nil {
			d = mod(*r.Day-n.Day(), daysInMonth(n.Year(), n.Month()))
			n = n.AddDate(0, 0, int(d))
			if d > 0 {
				continue
			}
		}

		// If a month is set
		if r.Month != nil {
			// If this is the wrong month
			if int(n.Month()) != *r.Month {
				// Rollover to the first of the next month. This is obviously not as efficient as going
				// directly to the required month, but adding months is tricky because their lengths aren't
				// equal. For example, adding one month to October 31 yields December 1.
				n = n.AddDate(0, 0, daysInMonth(n.Year(), n.Month())-n.Day()+1)
				continue
			}
		}

		done = true
	}

	if !done {
		panic("Failed to calculate next run time")
	}

	return n, nil
}

// mod calculates a non-negative modulus
func mod(x int, m int) time.Duration {
	result := x % m
	if result < 0 {
		result = result + m
	}
	return time.Duration(result)
}

func daysInMonth(year int, m time.Month) int {
	// The zeroth day of month m+1 will be normalised to the last day of month m
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
