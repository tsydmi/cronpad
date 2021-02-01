package report

import (
	"fmt"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"time"
)

type ProjectReport struct {
	HoursSum     float64      `json:"hoursSum"`
	DaysPassed   int64        `json:"daysPassed"`
	DaysAhead    int64        `json:"daysAhead"`
	TagChart     PieChart     `json:"tagChart"`
	TagSummaries []TagSummary `json:"tagSummaryTable"`
}

type TagSummary struct {
	Name    string         `json:"name"`
	Parent  repository.Tag `json:"parent"`
	Hours   float64        `json:"hours"`
	Percent string         `json:"percent"`
}

func CreateProjectReport(project repository.Project, days []repository.Day, tags []repository.Tag) ProjectReport {
	basicTagsChartSourceMap := make(map[string]chartSource)
	nonBasicTagsMap := make(map[string]repository.Tag)
	allTagsMap := make(map[string]repository.Tag)

	for _, tag := range tags {
		allTagsMap[tag.ID] = tag
		if tag.Basic {
			basicTagsChartSourceMap[tag.ID] = chartSource{ID: tag.ID, Name: tag.Name, Color: tag.Color}
		} else {
			nonBasicTagsMap[tag.ID] = tag
		}
	}

	var hoursSum float64 = 0
	hoursPerTag := make(hoursPerKeyMap)
	hoursPerBaseTag := make(hoursPerKeyMap)

	for _, day := range days {
		for _, event := range day.Events {
			if event.ProjectID == project.ID {
				hours := event.End.Sub(event.Start).Hours()
				hoursSum += hours

				hoursPerTag.increment(event.TagID, hours)
				baseTagId := getBaseTagId(event.TagID, basicTagsChartSourceMap, nonBasicTagsMap)
				hoursPerBaseTag.increment(baseTagId, hours)
			}
		}
	}

	now := time.Now()
	var daysPassed int64
	var daysAhead int64

	if project.Start != nil && project.End != nil {
		if project.Start.After(now) {
			daysPassed = 0
			daysAhead = int64(project.End.Sub(*project.Start).Hours() / 24)
		} else {
			if project.End.Before(now) {
				daysPassed = int64(project.End.Sub(*project.Start).Hours() / 24)
				daysAhead = 0
			} else {
				daysPassed = int64(now.Sub(*project.Start).Hours() / 24)
				daysAhead = int64(project.End.Sub(now).Hours() / 24)
			}
		}
	}

	tagMap := make(map[string]chartSource)
	for i := range tags {
		tag := tags[i]
		tagMap[tag.ID] = chartSource{ID: tag.ID, Name: tag.Name, Color: tag.Color}
	}

	return ProjectReport{
		HoursSum:     hoursSum,
		DaysPassed:   daysPassed,
		DaysAhead:    daysAhead,
		TagChart:     createPieChart(hoursPerBaseTag, basicTagsChartSourceMap),
		TagSummaries: createTagSummaries(hoursPerTag, allTagsMap, hoursSum),
	}
}

func createTagSummaries(hoursPerTagName hoursPerKeyMap, tagsMap map[string]repository.Tag, hoursSum float64) []TagSummary {
	result := make([]TagSummary, 0)

	for tagID := range hoursPerTagName {
		hours := hoursPerTagName[tagID]

		eventSummary := TagSummary{
			Hours:   hours,
			Percent: fmt.Sprintf("%.2f%%", hours*100/hoursSum),
		}

		tag, ok := tagsMap[tagID]
		if ok {
			eventSummary.Name = tag.Name
			if len(tag.ParentID) > 0 {
				parentTag, ok := tagsMap[tag.ParentID]
				if ok {
					eventSummary.Parent = parentTag
				}
			}
		}

		result = append(result, eventSummary)
	}

	return result
}
