package report

import (
	"github.com/ts-dmitry/cronpad/backend/repository"
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
	FindAll() ([]repository.Tag, error)
}

type reportProjectStore interface {
	GetProjectByID(projectID string) (repository.Project, error)
	FindAllProjectsByUser(userID string) ([]repository.Project, error)
}

func CreateReportService(dayStore reportDayStore, tagStore reportTagStore, projectStore reportProjectStore) *ReportService {
	return &ReportService{dayStore: dayStore, tagStore: tagStore, projectStore: projectStore}
}

func (t *ReportService) CalculateUserReport(form repository.DaySearchForm) (UserReport, error) {
	days, err := t.dayStore.Search(form)
	if err != nil {
		return UserReport{}, err
	}

	tags, err := t.tagStore.FindAll()
	if err != nil {
		return UserReport{}, err
	}

	projects, err := t.projectStore.FindAllProjectsByUser(form.UserID)
	if err != nil {
		return UserReport{}, err
	}

	return CreateUserReport(days, tags, projects, form.From.UTC(), form.To.UTC()), nil
}

func (t *ReportService) CalculateProjectReport(projectID string) (ProjectReport, error) {
	project, err := t.projectStore.GetProjectByID(projectID)
	if err != nil {
		return ProjectReport{}, err
	}

	var days = make([]repository.Day, 0) //TODO filter days here or return projection form mongo?
	if len(project.Users) > 0 {
		days, err = t.dayStore.Search(daySearhcFormFromProject(project))
		if err != nil {
			return ProjectReport{}, err
		}
	}

	tags, err := t.tagStore.FindAll()
	if err != nil {
		return ProjectReport{}, err
	}

	return CreateProjectReport(project, days, tags), nil
}

func daySearhcFormFromProject(project repository.Project) repository.DaySearchForm {
	var form = repository.DaySearchForm{UserIDs: project.Users}
	if project.Start != nil && !project.Start.IsZero() {
		form.From = *project.Start
	}
	if project.End != nil && !project.End.IsZero() {
		form.To = *project.End
	}
	return form
}
