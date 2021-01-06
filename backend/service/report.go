package service

type ChartReport struct {
	TagChart               PieChart          `json:"tagChart"`
	ProjectChart           PieChart          `json:"projectChart"`
	DateRangeChartDataSets []BarChartDataSet `json:"dateRangeChartDataSets"`
	HoursSum               float64           `json:"hoursSum"`
	EventSummaryTable      []EventSummary    `json:"eventSummaryTable"`
}

type EventSummary struct {
	Name    string  `json:"name"`
	Hours   float64 `json:"hours"`
	Percent string  `json:"percent"`
}
