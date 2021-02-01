package report

import (
	"github.com/ts-dmitry/cronpad/backend/utils"
	"time"
)

type hoursPerKeyMap map[string]float64

func (h hoursPerKeyMap) increment(key string, hours float64) {
	value, ok := h[key]
	if ok {
		h[key] = value + hours
	} else {
		h[key] = hours
	}
}

type hoursPerDayMap map[time.Time]float64

func (h hoursPerDayMap) increment(day time.Time, hours float64) {
	value, ok := h[day]
	if ok {
		h[day] = value + hours
	} else {
		h[day] = hours
	}
}

type chartSource struct {
	ID    string
	Name  string
	Color string
}

type eventSummaryMap map[string]UserReportEventSummary

func (e eventSummaryMap) increment(key string, hours float64, tagID string) {
	value, ok := e[key]
	if ok {
		value.Hours += hours
		if !utils.Contains(value.Tags, tagID) {
			value.Tags = append(value.Tags, tagID)
		}

		e[key] = value
	} else {
		e[key] = UserReportEventSummary{Name: key, Hours: hours, Tags: []string{tagID}}
	}
}
