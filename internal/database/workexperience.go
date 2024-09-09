package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type WorkResponsibility struct {
	Id             string    `json:"id"`
	Responsibility string    `json:"responsibility"`
	CreatedAt      time.Time `json:"createdAt"`
}

type WorkExperience struct {
	Id               string               `json:"id"`
	ResumeId         string               `json:"resumeId"`
	Employer         string               `json:"employer"`
	Title            string               `json:"title"`
	Began            time.Time            `json:"began"`
	Current          bool                 `json:"current"`
	CreatedAt        time.Time            `json:"createdAt"`
	Responsibilities []WorkResponsibility `json:"responsibilities"`
	Location         *string              `json:"location"`
	Finished         *time.Time           `json:"finished"`
}

func CreateWorkExperience(
	db VersionedDatabase,
	resume *Resume,
	employer string,
	title string,
	began time.Time,
	current bool,
	location *string,
	finished *time.Time,
) (WorkExperience, error) {
	id := uuid.New().String()
	createdAt := time.Now()

	query := `
INSERT INTO work_experiences (
  work_experience_id,
  resume_id,
  employer,
  title,
  began,
  current,
  created_at,
  location,
  finished
) VALUES (
  ?, ?, ?,
  ?, ?, ?,
  ?, ?, ?
)
  `

	var finishedInt int64
	if finished != nil {
		finishedInt = finished.Unix()
	}

	result, err := db.DB().Exec(query, id, resume.Id, employer, title, began.Unix(), current, createdAt.Unix(), location, &finishedInt)
	if err != nil {
		return WorkExperience{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return WorkExperience{}, err
	} else if rowsAffected != 1 {
		return WorkExperience{}, fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	responsibilities := make([]WorkResponsibility, 0)
	return WorkExperience{
		Id:               id,
		ResumeId:         resume.Id,
		Employer:         employer,
		Title:            title,
		Began:            began,
		Current:          current,
		CreatedAt:        createdAt,
		Location:         location,
		Finished:         finished,
		Responsibilities: responsibilities,
	}, nil
}

func GetWorkExperience(db VersionedDatabase, id string) *WorkExperience {
	query := `
SELECT
  we.work_experience_id,
  we.resume_id,
  we.employer,
  we.title,
  we.began,
  we.current,
  we.created_at,
  we.location,
  we.finished,
  wr.work_responsibility_id,
  wr.responsibility,
  wr.created_at AS responsibility_created_at
FROM work_experiences AS we
LEFT JOIN work_responsibilities AS wr ON (we.work_experience_id = wr.work_experience_id)
WHERE we.work_experience_id = ?
  `

	rows, err := db.DB().Query(query, id)
	if err != nil || rows == nil {
		return nil
	}

	var workExperience *WorkExperience
	responsibilities := make([]WorkResponsibility, 0)

	for rows.Next() {
		w, r, err := rowsToWorkExperience(rows)
		if err != nil {
			return nil
		}

		if workExperience == nil {
			workExperience = &w
		}

		if r != nil {
			responsibilities = append(responsibilities, *r)
		}
	}

	if workExperience != nil {
		workExperience.Responsibilities = responsibilities
	}

	return workExperience
}

func CreateWorkResponsibility(
	db VersionedDatabase,
	workExperience *WorkExperience,
	responsibility string,
) (WorkResponsibility, error) {
	id := uuid.New().String()
	createdAt := time.Now()

	query := `
INSERT INTO work_responsibilities (
  work_responsibility_id,
  work_experience_id,
  responsibility,
  created_at
) VALUES (
  ?, ?, ?,
  ?
)
  `

	result, err := db.DB().Exec(query, id, workExperience.Id, responsibility, createdAt.Unix())
	if err != nil {
		return WorkResponsibility{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return WorkResponsibility{}, err
	} else if rowsAffected != 1 {
		return WorkResponsibility{}, fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	r := WorkResponsibility{
		Id:             id,
		Responsibility: responsibility,
		CreatedAt:      createdAt,
	}

	workExperience.Responsibilities = append(workExperience.Responsibilities, r)

	return r, nil
}

func rowsToWorkExperience(rows *sql.Rows) (WorkExperience, *WorkResponsibility, error) {
	var we struct {
		Id        string
		ResumeId  string
		Employer  string
		Title     string
		Began     int64
		Current   int64
		CreatedAt int64
		Location  *string
		Finished  *int64
	}

	var r struct {
		Id             *string
		Responsibility *string
		CreatedAt      *int64
	}

	if err := rows.Scan(
		&we.Id,
		&we.ResumeId,
		&we.Employer,
		&we.Title,
		&we.Began,
		&we.Current,
		&we.CreatedAt,
		&we.Location,
		&we.Finished,
		&r.Id,
		&r.Responsibility,
		&r.CreatedAt,
	); err != nil {
		return WorkExperience{}, nil, err
	}

	current := false
	if we.Current != 0 {
		current = true
	}
	var finished *time.Time
	if we.Finished != nil {
		val := time.Unix(*we.Finished, 0)
		finished = &val
	}

	workExperience := WorkExperience{
		Id:        we.Id,
		ResumeId:  we.ResumeId,
		Employer:  we.Employer,
		Title:     we.Title,
		Began:     time.Unix(we.Began, 0),
		Current:   current,
		CreatedAt: time.Unix(we.CreatedAt, 0),
		Location:  we.Location,
		Finished:  finished,
	}

	if r.Id == nil {
		return workExperience, nil, nil
	}

	createdAt := time.Unix(*r.CreatedAt, 0)

	responsibility := WorkResponsibility{
		Id:             *r.Id,
		Responsibility: *r.Responsibility,
		CreatedAt:      createdAt,
	}

	return workExperience, &responsibility, nil
}
