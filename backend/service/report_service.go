package service

import (
	"fmt"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"time"
)

type ReportService struct {
	dayStore     reportDayStore
	tagStore     reportTagStore
	projectStore reportProjectStore
}

type reportDayStore interface {
	Search(form repository.DaySearchForm) ([]repository.Day, error)
}

type reportTagStore interface {
	FindAll(userID string) ([]repository.Tag, error)
}

type reportProjectStore interface {
	FindAllProjectsByUser(userID string) ([]repository.Project, error)
}

func CreateReportService(dayStore reportDayStore, tagStore reportTagStore, projectStore reportProjectStore) *ReportService {
	return &ReportService{dayStore: dayStore, tagStore: tagStore, projectStore: projectStore}
}

func (t *ReportService) Search(form repository.DaySearchForm) (ChartReport, error) {
	days, err := t.dayStore.Search(form)
	if err != nil {
		return ChartReport{}, err
	}

	tags, err := t.tagStore.FindAll(form.UserID)
	if err != nil {
		return ChartReport{}, err
	}

	tagsChartSourceMap := make(map[string]ChartSource)
	for i := range tags {
		tag := tags[i]
		tagsChartSourceMap[tag.ID] = ChartSource{ID: tag.ID, Name: tag.Name, Color: tag.Color}
	}

	projects, err := t.projectStore.FindAllProjectsByUser(form.UserID)
	if err != nil {
		return ChartReport{}, err
	}

	projectChartSourceMap := make(map[string]ChartSource)
	for i := range projects {
		project := projects[i]
		projectChartSourceMap[project.ID] = ChartSource{ID: project.ID, Name: project.Name, Color: getNextProjectDatasetColor()}
	}

	return convertToChartReport(days, tagsChartSourceMap, projectChartSourceMap, form.From.UTC(), form.To.UTC()), nil
}

func convertToChartReport(days []repository.Day, tagMap map[string]ChartSource, projectMap map[string]ChartSource, timeFrom time.Time, timeTo time.Time) ChartReport {
	var hoursSum float64 = 0
	hoursPerTag := make(hoursPerKeyMap)
	hoursPerProject := make(hoursPerKeyMap)
	hoursPerEventName := make(hoursPerKeyMap)

	for i := range days {
		day := days[i]
		for j := range day.Events {
			event := day.Events[j]

			hours := event.End.Sub(event.Start).Hours()
			hoursSum += hours

			hoursPerTag.increment(event.TagID, hours)
			hoursPerProject.increment(event.ProjectID, hours)
			hoursPerEventName.increment(event.Name, hours)
		}
	}

	return ChartReport{
		TagChart:               createPieChart(hoursPerTag, tagMap),
		ProjectChart:           createPieChart(hoursPerProject, projectMap),
		DateRangeChartDataSets: createDateRangeChartDataSet(days, hoursPerEventName, timeFrom, timeTo),
		HoursSum:               hoursSum,
		EventSummaryTable:      createEventSummaries(hoursPerEventName, hoursSum),
	}
}

func createEventSummaries(hoursPerEventName map[string]float64, hoursSum float64) []EventSummary {
	result := make([]EventSummary, 0)

	for name := range hoursPerEventName {
		hours := hoursPerEventName[name]
		percent := fmt.Sprintf("%.2f%%", hours*100/hoursSum)
		result = append(result, EventSummary{Name: name, Hours: hours, Percent: percent})
	}

	return result
}

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

type ChartSource struct {
	ID    string
	Name  string
	Color string
}
