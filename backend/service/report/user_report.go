package report

import (
	"fmt"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"log"
	"time"
)

type UserReport struct {
	TagChart               PieChart                 `json:"tagChart"`
	ProjectChart           PieChart                 `json:"projectChart"`
	DateRangeChartDataSets []BarChartDataSet        `json:"dateRangeChartDataSets"`
	HoursSum               float64                  `json:"hoursSum"`
	EventSummaryTable      []UserReportEventSummary `json:"eventSummaryTable"`
}

type UserReportEventSummary struct {
	Name    string   `json:"name"`
	Hours   float64  `json:"hours"`
	Tags    []string `json:"tags"`
	Percent string   `json:"percent"`
}

func CreateUserReport(days []repository.Day, tags []repository.Tag, projects []repository.Project, timeFrom time.Time, timeTo time.Time) UserReport {
	basicTagsChartSourceMap := make(map[string]chartSource)
	nonBasicTagsMap := make(map[string]repository.Tag)
	for _, tag := range tags {
		if tag.Basic {
			basicTagsChartSourceMap[tag.ID] = chartSource{ID: tag.ID, Name: tag.Name, Color: tag.Color}
		} else {
			nonBasicTagsMap[tag.ID] = tag
		}
	}

	var hoursSum float64 = 0
	hoursPerBaseTag := make(hoursPerKeyMap)
	hoursPerProject := make(hoursPerKeyMap)
	hoursPerEventName := make(hoursPerKeyMap)

	eventSummaryPerEventName := make(eventSummaryMap, 0)

	for _, day := range days {
		for _, event := range day.Events {
			hours := event.End.Sub(event.Start).Hours()
			hoursSum += hours

			baseTagId := getBaseTagId(event.TagID, basicTagsChartSourceMap, nonBasicTagsMap)
			hoursPerBaseTag.increment(baseTagId, hours)

			hoursPerProject.increment(event.ProjectID, hours)
			hoursPerEventName.increment(event.Name, hours)

			eventSummaryPerEventName.increment(event.Name, hours, event.TagID)
		}
	}

	return UserReport{
		TagChart:               createPieChart(hoursPerBaseTag, basicTagsChartSourceMap),
		ProjectChart:           createPieChart(hoursPerProject, createProjectChartSourceMap(projects)),
		DateRangeChartDataSets: createDateRangeChartDataSet(days, hoursPerEventName, timeFrom, timeTo),
		HoursSum:               hoursSum,
		EventSummaryTable:      createEventSummaries(eventSummaryPerEventName, hoursSum),
	}
}

func createProjectChartSourceMap(projects []repository.Project) map[string]chartSource {
	projectMap := make(map[string]chartSource)
	for i := range projects {
		project := projects[i]
		projectMap[project.ID] = chartSource{ID: project.ID, Name: project.Name, Color: getNextProjectDatasetColor()}
	}
	return projectMap
}

func getBaseTagId(tagID string, basicTagsMap map[string]chartSource, projectTagsMap map[string]repository.Tag) string {
	_, ok := basicTagsMap[tagID]
	if !ok {
		tag, ok := projectTagsMap[tagID]
		if ok {
			if len(tag.ParentID) == 0 {
				log.Printf("[WARN] tag %v should extend basic tag", tag.ID)
				return tagID
			}

			_, ok := basicTagsMap[tag.ParentID]
			if !ok {
				log.Printf("[WARN] tag %v shouldn't extend non-basic tag %v", tag.ID, tag.ParentID)
			}

			return tag.ParentID
		} else {
			log.Printf("[WARN] tag %v not exists", tag.ID)
		}
	}
	return tagID
}

func createEventSummaries(eventSummaryMap eventSummaryMap, hoursSum float64) []UserReportEventSummary {
	result := make([]UserReportEventSummary, 0)

	for _, summary := range eventSummaryMap {
		summary.Percent = fmt.Sprintf("%.2f%%", summary.Hours*100/hoursSum)
		result = append(result, summary)
	}

	return result
}
