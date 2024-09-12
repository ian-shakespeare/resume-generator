package database

import (
	"database/sql"
	"fmt"
	"resumegenerator/pkg/resume"
)

func CreateResume(
	db VersionedDatabase,
	u *User,
	r *resume.Resume,
) error {
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

	result, err := db.DB().Exec(
		query,
		r.Id,
		u.Id,
		r.Name,
		r.Email,
		r.PhoneNumber,
		r.Prelude,
		strToStrPtr(r.Location),
		strToStrPtr(r.LinkedIn),
		strToStrPtr(r.Github),
		strToStrPtr(r.Facebook),
		strToStrPtr(r.Instagram),
		strToStrPtr(r.Twitter),
		strToStrPtr(r.Portfolio),
		r.CreatedAt.Unix(),
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rowsAffected != 1 {
		return fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	return nil
}

func GetResume(db VersionedDatabase, resumeId string, userId string) *resume.Resume {
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
AND user_id = ?
  `

	row := db.DB().QueryRow(query, resumeId, userId)
	if row == nil {
		return nil
	}

	resume, err := rowToResume(row)
	if err != nil {
		return nil
	}

	return &resume
}

func ResumeEducations(db VersionedDatabase, r *resume.Resume) ([]resume.Education, error) {
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

func ResumeWorkExperiences(db VersionedDatabase, r *resume.Resume) ([]resume.WorkExperience, error) {
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

func ResumeProjects(db VersionedDatabase, r *resume.Resume) ([]resume.Project, error) {
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

func rowToResume(row *sql.Row) (resume.Resume, error) {
	var r struct {
		Id          string
		UserId      string
		Name        string
		Email       string
		PhoneNumber string
		Prelude     string
		Location    *string
		LinkedIn    *string
		Github      *string
		Facebook    *string
		Instagram   *string
		Twitter     *string
		Portfolio   *string
	}
	if err := row.Scan(
		&r.Id,
		&r.UserId,
		&r.Name,
		&r.Email,
		&r.PhoneNumber,
		&r.Prelude,
		&r.Location,
		&r.LinkedIn,
		&r.Github,
		&r.Facebook,
		&r.Instagram,
		&r.Twitter,
		&r.Portfolio,
	); err != nil {
		return resume.Resume{}, nil
	}

	return resume.Resume{
		Id:          r.Id,
		Name:        r.Name,
		Email:       r.Email,
		PhoneNumber: r.PhoneNumber,
		Prelude:     r.Prelude,
		Location:    strPtrToStr(r.Location),
		LinkedIn:    strPtrToStr(r.LinkedIn),
		Github:      strPtrToStr(r.Github),
		Facebook:    strPtrToStr(r.Facebook),
		Instagram:   strPtrToStr(r.Instagram),
		Twitter:     strPtrToStr(r.Twitter),
		Portfolio:   strPtrToStr(r.Portfolio),
	}, nil
}
