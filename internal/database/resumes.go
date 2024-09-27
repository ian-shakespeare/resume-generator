package database

import (
	"fmt"
	"resumegenerator/pkg/resume"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Resume struct {
	resume.Resume

	Id        string
	UserId    string
	CreatedAt time.Time
}

func GetResume(c Connection, id string) (Resume, error) {
	return Resume{}, nil
}

func (r *Resume) AddSkill(c Connection, skills ...string) error {
	if len(skills) < 1 {
		return nil
	}

	row := c.Database().QueryRow(`SELECT MAX(position) FROM skills WHERE resume_id = ?`, r.Id)
	var startingPos int
	if err := row.Scan(&startingPos); err != nil {
		startingPos = -1
	}
	startingPos += 1

	queryStrs := make([]string, 0, len(skills))
	var args []interface{}

	for i, skill := range skills {
		queryStrs = append(queryStrs, "(?, ?, ?, ?)")
		args = append(args, uuid.NewString())
		args = append(args, r.Id)
		args = append(args, startingPos+i)
		args = append(args, skill)
	}

	query := fmt.Sprintf(`
  INSERT INTO skills (
    skill_id,
    resume_id,
    position,
    skill
  ) VALUES %s
  `, strings.Join(queryStrs, ","))

	result, err := c.Database().Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != int64(len(skills)) {
		return fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	r.Skills = append(r.Skills, skills...)

	return nil
}

func (r *Resume) AddAchievement(c Connection, achievements ...string) error {
	if len(achievements) < 1 {
		return nil
	}

	row := c.Database().QueryRow(`SELECT MAX(position) FROM achievements WHERE resume_id = ?`, r.Id)
	var startingPos int
	if err := row.Scan(&startingPos); err != nil {
		startingPos = -1
	}
	startingPos += 1

	queryStrs := make([]string, 0, len(achievements))
	var args []interface{}

	for i, achievement := range achievements {
		queryStrs = append(queryStrs, "(?, ?, ?, ?)")
		args = append(args, uuid.NewString())
		args = append(args, r.Id)
		args = append(args, startingPos+i)
		args = append(args, achievement)
	}

	query := fmt.Sprintf(`
  INSERT INTO achievements (
    achievement_id,
    resume_id,
    position,
    achievement
  ) VALUES %s
  `, strings.Join(queryStrs, ","))

	result, err := c.Database().Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != int64(len(achievements)) {
		return fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	r.Achievements = append(r.Achievements, achievements...)

	return nil
}

func (r *Resume) AddEducation(c Connection, educations ...resume.Education) error {
	if len(educations) < 1 {
		return nil
	}

	queryStrs := make([]string, 0, len(educations))
	var args []interface{}

	for _, e := range educations {
		var finished *int64
		if e.Finished != nil {
			value := e.Finished.Unix()
			finished = &value
		}

		queryStrs = append(queryStrs, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args, uuid.NewString())
		args = append(args, r.Id)
		args = append(args, time.Now().Unix())
		args = append(args, e.Degree)
		args = append(args, e.FieldOfStudy)
		args = append(args, e.Institution)
		args = append(args, e.Began.Unix())
		args = append(args, e.Current)
		args = append(args, e.Location)
		args = append(args, finished)
		args = append(args, e.GPA)
	}

	query := fmt.Sprintf(`
  INSERT INTO educations (
    education_id,
    resume_id,
    created_at,
    degree,
    field,
    institution,
    began,
    current,
    location,
    finished,
    gpa
  ) VALUES %s
  `, strings.Join(queryStrs, ","))

	result, err := c.Database().Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != int64(len(educations)) {
		return fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	r.Educations = append(r.Educations, educations...)

	return nil
}

func (r *Resume) AddWorkExperience(c Connection, workExperiences ...resume.WorkExperience) error {
	if len(workExperiences) < 1 {
		return nil
	}

	var createds []WorkExperience

	queryStrs := make([]string, 0, len(workExperiences))
	var args []interface{}

	for _, w := range workExperiences {
		var finished *int64
		if w.Finished != nil {
			value := w.Finished.Unix()
			finished = &value
		}

		created := WorkExperience{
			Id:             uuid.NewString(),
			ResumeId:       r.Id,
			CreatedAt:      time.Now(),
			WorkExperience: w,
		}
		created.Responsibilities = []string{}
		createds = append(createds, created)

		queryStrs = append(queryStrs, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args, created.Id)
		args = append(args, created.ResumeId)
		args = append(args, created.CreatedAt.Unix())
		args = append(args, created.Employer)
		args = append(args, created.Title)
		args = append(args, created.Began.Unix())
		args = append(args, created.Current)
		args = append(args, created.Location)
		args = append(args, finished)
	}

	query := fmt.Sprintf(`
  INSERT INTO work_experiences (
    work_experience_id,
    resume_id,
    created_at,
    employer,
    title,
    began,
    current,
    location,
    finished
  ) VALUES %s
  `, strings.Join(queryStrs, ","))

	result, err := c.Database().Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != int64(len(workExperiences)) {
		return fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	r.WorkExperiences = append(r.WorkExperiences, workExperiences...)

	for i, created := range createds {
		created.AddResponsibility(c, workExperiences[i].Responsibilities...)
	}

	return nil
}

func (r *Resume) AddProject(c Connection, projects ...resume.Project) error {
	if len(projects) < 1 {
		return nil
	}

	var createds []Project

	queryStrs := make([]string, 0, len(projects))
	var args []interface{}

	for _, p := range projects {
		created := Project{
			Id:        uuid.NewString(),
			ResumeId:  r.Id,
			CreatedAt: time.Now(),
			Project:   p,
		}
		created.Responsibilities = []string{}
		createds = append(createds, created)

		queryStrs = append(queryStrs, "(?, ?, ?, ?, ?, ?)")
		args = append(args, created.Id)
		args = append(args, created.ResumeId)
		args = append(args, created.CreatedAt.Unix())
		args = append(args, created.Name)
		args = append(args, created.Description)
		args = append(args, created.Role)
	}

	query := fmt.Sprintf(`
  INSERT INTO projects (
    project_id,
    resume_id,
    created_at,
    name,
    description,
    role
  ) VALUES %s
  `, strings.Join(queryStrs, ","))

	result, err := c.Database().Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != int64(len(projects)) {
		return fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	r.Projects = append(r.Projects, projects...)

	for i, created := range createds {
		created.AddResponsibility(c, projects[i].Responsibilities...)
	}

	return nil
}
