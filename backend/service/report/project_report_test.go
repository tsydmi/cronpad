package report

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestCreateProjectReport(t *testing.T) {
	t.Run("Empty arrays", func(t *testing.T) {
		now := time.Now()
		start := now.AddDate(0, 0, -2)
		end := now.AddDate(0, 0, 5)

		project := repository.Project{
			ID:    uuid.New().String(),
			Start: &start,
			End:   &end,
		}

		result := CreateProjectReport(project, []repository.Day{}, []repository.Tag{})

		assert.Equal(t, float64(0), result.HoursSum, "wrong HoursSum")
		assert.Equal(t, int64(2), result.DaysPassed, "wrong DaysPassed")
		assert.Equal(t, int64(5), result.DaysAhead, "wrong DaysAhead")

	})

	t.Run("Base tags only", func(t *testing.T) {
		daysAmount := rand.Intn(30)
		eventsPerDay := 5

		now := time.Now()
		start := now.AddDate(0, 0, -(daysAmount - 1))

		project := repository.Project{
			ID:    uuid.New().String(),
			Start: &start,
			End:   &now,
		}
		users := []string{uuid.New().String()}
		baseTags := generateBaseTags(3)
		days := generateRandomDays(project, users, baseTags, eventsPerDay)

		result := CreateProjectReport(project, days, baseTags)

		assert.Equal(t, float64(daysAmount*eventsPerDay), result.HoursSum, "wrong HoursSum")
		assert.Equal(t, int64(daysAmount-1), result.DaysPassed, "wrong DaysPassed")
		assert.Equal(t, int64(0), result.DaysAhead, "wrong DaysAhead")

		assert.Len(t, result.TagChart.Datasets, 1)
		assert.Equal(t, result.HoursSum, sumFloats(result.TagChart.Datasets[0].Data), "sum of TagChart dataset should be equal to HoursSum")
		assert.Equal(t, result.HoursSum, sumTagSummariesHours(result.TagSummaries), "sum of TagSummaries dataset should be equal to HoursSum")
	})

	t.Run("Different tags", func(t *testing.T) {
		daysAmount := rand.Intn(30)
		eventsPerDay := 5

		now := time.Now()
		start := now.AddDate(0, 0, -(daysAmount - 1))

		project := repository.Project{
			ID:    uuid.New().String(),
			Start: &start,
			End:   &now,
		}
		users := []string{uuid.New().String()}

		baseTags := generateBaseTags(3)
		tags := generateTags(10, baseTags)
		allTags := append(tags, baseTags...)

		days := generateRandomDays(project, users, allTags, eventsPerDay)

		result := CreateProjectReport(project, days, allTags)

		assert.Equal(t, float64(daysAmount*eventsPerDay), result.HoursSum, "wrong HoursSum")
		assert.Equal(t, int64(daysAmount-1), result.DaysPassed, "wrong DaysPassed")
		assert.Equal(t, int64(0), result.DaysAhead, "wrong DaysAhead")

		assert.Len(t, result.TagChart.Datasets, 1)
		assert.Equal(t, result.HoursSum, sumFloats(result.TagChart.Datasets[0].Data), "sum of TagChart dataset should be equal to HoursSum")
		assert.Equal(t, result.HoursSum, sumTagSummariesHours(result.TagSummaries), "sum of TagSummaries dataset should be equal to HoursSum")
	})

	t.Run("Without passed to the method tags", func(t *testing.T) {
		daysAmount := rand.Intn(30)
		eventsPerDay := 5

		now := time.Now()
		start := now.AddDate(0, 0, -(daysAmount - 1))

		project := repository.Project{
			ID:    uuid.New().String(),
			Start: &start,
			End:   &now,
		}
		users := []string{uuid.New().String()}

		baseTags := generateBaseTags(3)
		tags := generateTags(10, baseTags)
		allTags := append(tags, baseTags...)

		days := generateRandomDays(project, users, allTags, eventsPerDay)

		result := CreateProjectReport(project, days, []repository.Tag{})

		assert.Equal(t, float64(daysAmount*eventsPerDay), result.HoursSum, "wrong HoursSum")
		assert.Equal(t, int64(daysAmount-1), result.DaysPassed, "wrong DaysPassed")
		assert.Equal(t, int64(0), result.DaysAhead, "wrong DaysAhead")

		assert.Len(t, result.TagChart.Datasets, 1)
		assert.Equal(t, result.HoursSum, sumFloats(result.TagChart.Datasets[0].Data), "sum of TagChart dataset should be equal to HoursSum")
		assert.Equal(t, result.HoursSum, sumTagSummariesHours(result.TagSummaries), "sum of TagSummaries dataset should be equal to HoursSum")
	})
}

func BenchmarkCreateProjectReport(b *testing.B) {
	start := time.Now()
	end := start.AddDate(0, 3, 0)

	users := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	project := repository.Project{ID: uuid.New().String(), Start: &start, End: &end}
	baseTags := generateBaseTags(5)
	tags := generateTags(50, baseTags)

	allTags := append(tags, baseTags...)

	days := generateRandomDays(project, users, allTags, 5)

	for i := 0; i < b.N; i++ {
		CreateProjectReport(project, days, allTags)
	}
}

func generateBaseTags(count int) []repository.Tag {
	tags := make([]repository.Tag, 0)
	for i := 1; i <= count; i++ {
		tags = append(tags, repository.Tag{
			ID:          uuid.New().String(),
			Name:        "name",
			Description: "description",
			Color:       "#FFFFFF",
			ParentID:    "",
			ProjectID:   "",
			Basic:       true,
		})
	}

	return tags
}

func generateTags(count int, basicTags []repository.Tag) []repository.Tag {
	tags := make([]repository.Tag, 0)
	for i := 1; i <= count; i++ {
		tags = append(tags, repository.Tag{
			ID:          uuid.New().String(),
			Name:        "name",
			Description: "description",
			Color:       "#FFFFFF",
			ParentID:    basicTags[rand.Intn(len(basicTags))].ID,
			ProjectID:   "",
			Basic:       true,
		})
	}

	return tags
}

func generateRandomDays(project repository.Project, users []string, tags []repository.Tag, eventsPerDay int) []repository.Day {
	currentDay := time.Date(project.Start.Year(), project.Start.Month(), project.Start.Day(), 0, 0, 0, 0, time.UTC)
	days := make([]repository.Day, 0)

	for ; !currentDay.After(*project.End); currentDay = currentDay.AddDate(0, 0, 1) {
		for _, userID := range users {
			days = append(days, repository.Day{
				ID:     uuid.New().String(),
				Events: generateRandomEvents(currentDay, eventsPerDay, tags, project),
				Date:   currentDay,
				UserID: userID,
			})
		}
	}

	return days
}

func generateRandomEvents(day time.Time, eventsPerDay int, tags []repository.Tag, project repository.Project) []repository.Event {
	events := make([]repository.Event, 0)

	startOfTheDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)

	for i := 1; i <= eventsPerDay; i++ {
		events = append(events, repository.Event{
			ID:        uuid.New().String(),
			Name:      "",
			Start:     startOfTheDay.Add(time.Duration(i * 60 * 60 * 1000_000_000)),
			End:       startOfTheDay.Add(time.Duration((i + 1) * 60 * 60 * 1000_000_000)),
			TagID:     tags[rand.Intn(len(tags))].ID,
			ProjectID: project.ID,
			Timed:     false,
		})
	}

	return events
}

func sumFloats(data []float64) float64 {
	sum := new(big.Float)
	for _, f := range data {
		sum.Add(sum, new(big.Float).SetFloat64(f))
	}
	result, _ := sum.Float64()
	return result
}

func sumTagSummariesHours(data []TagSummary) float64 {
	sum := new(big.Float)
	for _, f := range data {
		sum.Add(sum, new(big.Float).SetFloat64(f.Hours))
	}
	result, _ := sum.Float64()
	return result
}
