package service

const titleForUnknown = "unknown"
const colorForUnknown = "#E0E0E0"

type PieChart struct {
	Datasets []PieChartDataSet `json:"datasets"`
	Labels   []string          `json:"labels"`
}

type PieChartDataSet struct {
	Data            []float64 `json:"data"`
	BackgroundColor []string  `json:"backgroundColor"`
}

func createPieChart(hoursPerTag map[string]float64, tagMap map[string]ChartSource) PieChart {
	data := make([]float64, 0)
	colors := make([]string, 0)
	labels := make([]string, 0)

	for tagID := range hoursPerTag {
		hours := hoursPerTag[tagID]
		data = append(data, hours)

		tag, ok := tagMap[tagID]
		if ok == false {
			colors = append(colors, colorForUnknown)
			labels = append(labels, titleForUnknown)
		} else {
			colors = append(colors, tag.Color)
			labels = append(labels, tag.Name)
		}
	}

	return PieChart{
		Datasets: []PieChartDataSet{{Data: data, BackgroundColor: colors}},
		Labels:   labels,
	}
}

var nextProjectsChartColorIndex = 1
var projectsChartColors = []string{"#FFD54F", "#CDDC39", "#BA68C8", "#7986CB", "#4DD0E1", "#4DB6AC"}

func getNextProjectDatasetColor() string {
	nextProjectsChartColorIndex = (nextProjectsChartColorIndex + 1) % len(projectsChartColors)
	return projectsChartColors[nextProjectsChartColorIndex]
}
