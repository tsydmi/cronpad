package report

import (
	"github.com/ts-dmitry/cronpad/backend/repository"
	"time"
)

type BarChartDataSet struct {
	Label           string         `json:"label"`
	Data            []BarChartData `json:"data"`
	BackgroundColor string         `json:"backgroundColor"`
}

type BarChartData struct {
	X time.Time `json:"x"`
	Y float64   `json:"y"`
}

func createDateRangeChartDataSet(days []repository.Day, hoursPerEventName hoursPerKeyMap, timeFrom time.Time, timeTo time.Time) []BarChartDataSet {
	barCharDataSetMap := make(map[string]BarChartDataSet)

	namePerAllValuesMap := make(map[string]hoursPerDayMap)

	for i := range days {
		day := days[i]
		for j := range day.Events {
			event := day.Events[j]

			hours := event.End.Sub(event.Start).Hours()

			dayPerValueMap, ok := namePerAllValuesMap[event.Name]
			if ok {
				dayPerValueMap.increment(day.Date, hours)
			} else {
				dayPerValueMap = make(map[time.Time]float64)
				dayPerValueMap[day.Date] = hours
			}

			namePerAllValuesMap[event.Name] = dayPerValueMap
		}
	}

	for name := range hoursPerEventName {
		dayPerValue, ok := namePerAllValuesMap[name]
		if ok {
			barCharDataSetMap[name] = BarChartDataSet{
				Label:           name,
				Data:            createDataArrayWithValues(timeFrom, timeTo, dayPerValue),
				BackgroundColor: getNextDateRangeDatasetColor(),
			}
		}
	}

	result := make([]BarChartDataSet, 0)
	for s := range barCharDataSetMap {
		result = append(result, barCharDataSetMap[s])
	}

	return result
}

func createDataArrayWithValues(from time.Time, to time.Time, dayPerValue map[time.Time]float64) []BarChartData {
	result := make([]BarChartData, 0)

	for d := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC); !d.After(to); d = d.AddDate(0, 0, 1) {
		sum, ok := dayPerValue[d]
		if ok {
			result = append(result, BarChartData{X: d, Y: sum})
		} else {
			result = append(result, BarChartData{X: d, Y: 0})
		}
	}

	return result
}

var nextDataRangeChartColorIndex = 1
var dataRangeChartColors = []string{"#80DEEA", "#4DD0E1", "#26C6DA"}

func getNextDateRangeDatasetColor() string {
	nextDataRangeChartColorIndex = (nextDataRangeChartColorIndex + 1) % len(dataRangeChartColors)
	return dataRangeChartColors[nextDataRangeChartColorIndex]
}
