package report

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ts-dmitry/cronpad/backend/repository"
	"math/rand"
	"testing"
	"time"
)

func TestCreateUserReport(t *testing.T) {
	t.Run("Empty arrays", func(t *testing.T) {
		now := time.Now()
		start := now.AddDate(0, 0, -2)
		end := now.AddDate(0, 0, 5)

		result := CreateUserReport([]repository.Day{}, []repository.Tag{}, []repository.Project{}, start, end)

		assert.Equal(t, float64(0), result.HoursSum, "wrong HoursSum")
	})

	t.Run("Base tags only", func(t *testing.T) {
		daysAmount := rand.Intn(30)
		eventsPerDay := 5

		now := time.Now()
		start := now.AddDate(0, 0, -(daysAmount - 1))

		projects := []repository.Project{{ID: uuid.New().String()}, {ID: uuid.New().String()}}
		baseTags := generateBaseTags(3)
		days := generateDaysForUser(uuid.New().String(), baseTags, projects, 5, start, now)

		result := CreateUserReport(days, baseTags, projects, start, now)

		assert.Equal(t, float64(len(projects)*daysAmount*eventsPerDay), result.HoursSum, "wrong HoursSum")

		assert.Len(t, result.TagChart.Datasets, 1)
		assert.Equal(t, result.HoursSum, sumFloats(result.TagChart.Datasets[0].Data), "sum of TagChart dataset should be equal to HoursSum")

		assert.Len(t, result.ProjectChart.Datasets, 1)
		assert.Equal(t, result.HoursSum, sumFloats(result.ProjectChart.Datasets[0].Data), "sum of ProjectChart dataset should be equal to HoursSum")
	})

	t.Run("Different tags", func(t *testing.T) {
		daysAmount := rand.Intn(30)
		eventsPerDay := 5

		now := time.Now()
		start := now.AddDate(0, 0, -(daysAmount - 1))

		projects := []repository.Project{{ID: uuid.New().String()}, {ID: uuid.New().String()}}
		baseTags := generateBaseTags(3)
		tags := generateTags(10, baseTags)
		allTags := append(tags, baseTags...)

		days := generateDaysForUser(uuid.New().String(), allTags, projects, 5, start, now)

		result := CreateUserReport(days, baseTags, projects, start, now)

		assert.Equal(t, float64(len(projects)*daysAmount*eventsPerDay), result.HoursSum, "wrong HoursSum")

		assert.Len(t, result.TagChart.Datasets, 1)
		assert.Equal(t, result.HoursSum, sumFloats(result.TagChart.Datasets[0].Data), "sum of TagChart dataset should be equal to HoursSum")

		assert.Len(t, result.ProjectChart.Datasets, 1)
		assert.Equal(t, result.HoursSum, sumFloats(result.ProjectChart.Datasets[0].Data), "sum of ProjectChart dataset should be equal to HoursSum")
	})

	t.Run("Without passed to the method tags", func(t *testing.T) {
		daysAmount := rand.Intn(30)
		eventsPerDay := 5

		now := time.Now()
		start := now.AddDate(0, 0, -(daysAmount - 1))

		projects := []repository.Project{{ID: uuid.New().String()}, {ID: uuid.New().String()}}
		baseTags := generateBaseTags(3)
		tags := generateTags(10, baseTags)
		allTags := append(tags, baseTags...)

		days := generateDaysForUser(uuid.New().String(), allTags, projects, 5, start, now)

		result := CreateUserReport(days, baseTags, projects, start, now)

		assert.Equal(t, float64(len(projects)*daysAmount*eventsPerDay), result.HoursSum, "wrong HoursSum")

		assert.Len(t, result.TagChart.Datasets, 1)
		assert.Equal(t, result.HoursSum, sumFloats(result.TagChart.Datasets[0].Data), "sum of TagChart dataset should be equal to HoursSum")

		assert.Len(t, result.ProjectChart.Datasets, 1)
		assert.Equal(t, result.HoursSum, sumFloats(result.ProjectChart.Datasets[0].Data), "sum of ProjectChart dataset should be equal to HoursSum")
	})
}

func BenchmarkCreateUserReport(b *testing.B) {
	start := time.Now()
	end := start.AddDate(0, 3, 0)

	projects := []repository.Project{
		{ID: uuid.New().String(), Start: &start, End: &end},
		{ID: uuid.New().String(), Start: &start, End: &end},
	}
	baseTags := generateBaseTags(5)
	tags := generateTags(50, baseTags)

	allTags := append(tags, baseTags...)

	days := generateDaysForUser(uuid.New().String(), allTags, projects, 5, start, end)

	for i := 0; i < b.N; i++ {
		CreateUserReport(days, allTags, projects, start, end)
	}
}

func generateDaysForUser(userID string, tags []repository.Tag, projects []repository.Project, eventsPerDay int, start time.Time, end time.Time) []repository.Day {
	currentDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	days := make([]repository.Day, 0)

	for ; !currentDay.After(end); currentDay = currentDay.AddDate(0, 0, 1) {
		for _, project := range projects {
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
