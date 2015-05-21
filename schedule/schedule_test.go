package schedule

import (
	"reflect"
	"testing"
)

var us = Rotation{{"us1", "us2"}, {"us3"}, {"us4"}, {"us5"}}
var ind = Rotation{{"ind1", "ind2"}, {"ind3"}, {"ind4"}}

var se1 = ScheduleEntry{
	Name:     "uswest",
	Start:    "6:00AM",
	End:      "6:00PM",
	Location: "US/Pacific",
	Rotations: map[string]Rotation{
		"primary":   us,
		"secondary": ind,
	},
}

var se2 = ScheduleEntry{
	Name:     "india",
	Start:    "6:30AM",
	End:      "6:30PM",
	Location: "Asia/Kolkata",
	Rotations: map[string]Rotation{
		"primary":   ind,
		"secondary": us,
	},
}

func TestScheduleSimple(t *testing.T) {
	sc := Schedule{Entries: []ScheduleEntry{se1, se2}}
	now, err := ParseTimeInLocation(TIME_FULLY_QUALIFIED, "2015-03-09 9:15PM", "US/Pacific")
	if err != nil {
		t.Errorf("Failed to parse time!")
	}
	matches, err := sc.Get("primary", now.UTC())
	if err != nil {
		t.Errorf("Schedule simple", err)
	}
	if !reflect.DeepEqual(matches, []string{"ind1", "ind2"}) {
		t.Errorf("Expected=%v got=%v", []string{"ind1", "ind2"}, matches)
	}
}

func TestSingleOverride(t *testing.T) {
	override := Override{
		Location:  "US/Pacific",
		Start:     "2015-03-09 9:00PM",
		End:       "2015-03-09 9:32PM",
		Rotations: map[string][]string{"primary": {"us1"}},
		Comment:   "Out for dinner",
	}
	sc := Schedule{
		Entries:   []ScheduleEntry{se1},
		Overrides: []Override{override},
	}
	now, err := ParseTimeInLocation(TIME_FULLY_QUALIFIED, "2015-03-09 9:15PM", "US/Pacific")
	if err != nil {
		t.Errorf("Failed to parse time!")
	}
	matches, err := sc.Get("primary", now.UTC())
	if err != nil {
		t.Errorf("Overrides: ", err)
	}
	if !reflect.DeepEqual(matches, []string{"us1"}) {
		t.Errorf("Applying overrides failed: expected=%v got=%v", []string{"us1"}, matches)
	}
}

func TestMultipleOverrides1(t *testing.T) {
	override1 := Override{
		Location:  "US/Pacific",
		Start:     "2015-03-09 9:00PM",
		End:       "2015-03-09 9:32PM",
		Rotations: map[string][]string{"primary": {"us1"}},
		Comment:   "Out for dinner",
	}
	override2 := Override{
		Location:  "US/Pacific",
		Start:     "2015-03-09 9:18PM",
		End:       "2015-03-09 9:32PM",
		Rotations: map[string][]string{"primary": {"ind1"}},
		Comment:   "Came back from dinner",
	}
	sc := Schedule{
		Entries:   []ScheduleEntry{se1},
		Overrides: []Override{override1, override2},
	}
	now, err := ParseTimeInLocation(TIME_FULLY_QUALIFIED, "2015-03-09 9:18PM", "US/Pacific")
	if err != nil {
		t.Errorf("Failed to parse time!")
	}
	matches, err := sc.Get("primary", now.UTC())
	if err != nil {
		t.Errorf("Overrides: ", err)
	}
	if !reflect.DeepEqual(matches, []string{"ind1"}) {
		t.Errorf("Applying overrides failed: expected=%v got=%v", []string{"ind1"}, matches)
	}
}

func TestMultipleOverrides2(t *testing.T) {
	override1 := Override{
		Location:  "US/Pacific",
		Start:     "2015-03-09 9:00PM",
		End:       "2015-03-09 9:32PM",
		Rotations: map[string][]string{"primary": {"us1"}},
		Comment:   "covering for india",
	}
	override2 := Override{
		Location:  "Asia/Calcutta",
		Start:     "2015-03-10 9:48AM",
		End:       "2015-03-10 10:02AM",
		Rotations: map[string][]string{"primary": {"ind1"}},
		Comment:   "electricity is back!",
	}
	sc := Schedule{
		Entries:   []ScheduleEntry{se1},
		Overrides: []Override{override1, override2},
	}
	now, err := ParseTimeInLocation(TIME_FULLY_QUALIFIED, "2015-03-09 9:18PM", "US/Pacific")
	if err != nil {
		t.Errorf("Failed to parse time!")
	}
	matches, err := sc.Get("primary", now.UTC())
	if err != nil {
		t.Errorf("Overrides: ", err)
	}
	if !reflect.DeepEqual(matches, []string{"ind1"}) {
		t.Errorf("Applying overrides failed: expected=%v got=%v", []string{"ind1"}, matches)
	}
}
