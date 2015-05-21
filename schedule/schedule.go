package schedule

import (
	"time"
)

// These are guides for formatting - the numbers can't
// be arbitrary.
// See: http://golang.org/pkg/time/

const TIME_FULLY_QUALIFIED = "2006-01-02 3:04PM"
const TIME_RELATIVE = "3:04PM"

type Rotation [][]string

type Schedule struct {
	Entries   []ScheduleEntry
	Overrides []Override
}

type ScheduleEntry struct {
	Name      string
	Start     string
	End       string
	Location  string
	Rotations map[string]Rotation
}

type Override struct {
	Start     string
	End       string
	Rotations map[string][]string
	Comment   string
	Location  string
}

func ParseTimeInLocation(layout string, timeStr string, location string) (time.Time, error) {
	l, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation(layout, timeStr, l)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func adjustedDate(when time.Time, kitchenTime string, location string) (time.Time, error) {
	l, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation(TIME_RELATIVE, kitchenTime, l)
	if err != nil {
		return time.Time{}, err
	}
	now := when.In(l)
	d := time.Date(now.Year(), now.Month(), now.Day(),
		t.Hour(), t.Minute(), 0, 0, l)
	return d.UTC(), nil
}

// lookup current entry for a level and rotation
func (schedule *Schedule) Get(level string, when time.Time) ([]string, error) {
	var matches []string
	// find all possible matches from schedule
	for _, se := range schedule.Entries {
		// determine start/end of oncall periods that "when" belongs to
		start_time, err := adjustedDate(when, se.Start, se.Location)
		if err != nil {
			return nil, err
		}
		end_time, _ := adjustedDate(when, se.End, se.Location)
		if err != nil {
			return nil, err
		}
		if start_time.After(end_time) || start_time.Equal(end_time) {
			end_time = end_time.AddDate(0, 0, 1)
		}

		if start_time.After(when) || start_time.Equal(when) {
			start_time = start_time.AddDate(0, 0, -1)
			end_time = end_time.AddDate(0, 0, -1)
		}

		if when.Unix() >= start_time.Unix() && when.Unix() <= end_time.Unix() {
			val, ok := se.Rotations[level]
			if !ok {
				continue
			}
			idx := start_time.YearDay() % len(val)
			matches = append(matches, val[idx]...)
		}
	}

	// apply overrides
	for _, override := range schedule.Overrides {
		override_start, err := ParseTimeInLocation(TIME_FULLY_QUALIFIED, override.Start, override.Location)
		if err != nil {
			continue
		}
		override_end, err := ParseTimeInLocation(TIME_FULLY_QUALIFIED, override.End, override.Location)
		if err != nil {
			continue
		}
		if when.Unix() >= override_start.Unix() && when.Unix() <= override_end.Unix() {
			val, ok := override.Rotations[level]
			if !ok {
				continue
			}
			matches = val
		}
	}

	return matches, nil
}
