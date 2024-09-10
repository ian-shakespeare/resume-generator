package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Resume struct {
	Id          string    `json:"id"`
	UserId      string    `json:"userId"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Prelude     string    `json:"prelude"`
	CreatedAt   time.Time `json:"createdAt"`
	Location    *string   `json:"location"`
	LinkedIn    *string   `json:"linkedIn"`
	Github      *string   `json:"github"`
	Facebook    *string   `json:"facebook"`
	Instagram   *string   `json:"instagram"`
	Twitter     *string   `json:"twitter"`
	Portfolio   *string   `json:"portfolio"`
}

func CreateResume(
	db VersionedDatabase,
	user *User,
	name string,
	email string,
	phoneNumber string,
	prelude string,
	location *string,
	linkedIn *string,
	github *string,
	facebook *string,
	instagram *string,
	twitter *string,
	portfolio *string,
) (Resume, error) {
	id := uuid.New().String()
	createdAt := time.Now()

	query := `
INSERT INTO resumes (
  resume_id,
  user_id,
  name,
  email,
  phone_number,
  prelude,
  location,
  linkedin_username,
  github_username,
  facebook_username,
  instagram_username,
  twitter_handle,
  portfolio,
  created_at
) VALUES (
  ?, ?, ?,
  ?, ?, ?,
  ?, ?, ?,
  ?, ?, ?,
  ?, ?
)
  `

	result, err := db.DB().Exec(query, id, user.Id, name, email, phoneNumber, prelude, location, linkedIn, github, facebook, instagram, twitter, portfolio, createdAt.Unix())
	if err != nil {
		return Resume{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Resume{}, err
	} else if rowsAffected != 1 {
		return Resume{}, fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	return Resume{
		Id:          id,
		UserId:      user.Id,
		Name:        name,
		Email:       email,
		PhoneNumber: phoneNumber,
		Prelude:     prelude,
		Location:    location,
		LinkedIn:    linkedIn,
		Github:      github,
		Facebook:    facebook,
		Instagram:   instagram,
		Twitter:     twitter,
		Portfolio:   portfolio,
		CreatedAt:   createdAt,
	}, nil
}

func GetResume(db VersionedDatabase, id string) *Resume {
	query := `
SELECT
  resume_id,
  user_id,
  name,
  email,
  phone_number,
  prelude,
  location,
  linkedin_username,
  github_username,
  facebook_username,
  instagram_username,
  twitter_handle,
  portfolio
FROM resumes
WHERE resume_id = ?
  `

	row := db.DB().QueryRow(query, id)
	if row == nil {
		return nil
	}

	resume, err := rowToResume(row)
	if err != nil {
		return nil
	}

	return &resume
}

func (r *Resume) Educations(db VersionedDatabase) ([]Education, error) {
	query := `
SELECT
  education_id,
  resume_id,
  degree_type,
  field_of_study,
  institution,
  began,
  current,
  created_at,
  location,
  finished,
  gpa
FROM educations
WHERE resume_id = ?
  `

	rows, err := db.DB().Query(query, r.Id)
	if err != nil {
		return nil, err
	}

	educations, err := rowsToEducation(rows)

	return educations, err
}

func (r *Resume) WorkExperiences(db VersionedDatabase) ([]WorkExperience, error) {
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
WHERE we.resume_id = ?
  `

	rows, err := db.DB().Query(query, r.Id)
	if err != nil {
		return nil, err
	}

	workExperiences, err := rowsToWorkExperience(rows)
	return workExperiences, err
}

func (r *Resume) Projects(db VersionedDatabase) ([]Project, error) {
	query := `
SELECT
  p.project_id,
  p.resume_id,
  p.name,
  p.description,
  p.role,
  p.created_at,
  r.project_responsibility_id,
  r.responsibility,
  r.created_at AS responsibility_created_at
FROM projects AS p
LEFT JOIN project_responsibilities AS r ON (p.project_id = r.project_id)
WHERE p.resume_id = ?
  `

	rows, err := db.DB().Query(query, r.Id)
	if err != nil {
		return nil, err
	}

	projects, err := rowsToProject(rows)
	return projects, err
}

func rowToResume(row *sql.Row) (Resume, error) {
	var resume Resume
	if err := row.Scan(
		&resume.Id,
		&resume.UserId,
		&resume.Name,
		&resume.Email,
		&resume.PhoneNumber,
		&resume.Prelude,
		&resume.Location,
		&resume.LinkedIn,
		&resume.Github,
		&resume.Facebook,
		&resume.Instagram,
		&resume.Twitter,
		&resume.Portfolio,
	); err != nil {
		return Resume{}, nil
	}

	return resume, nil
}
