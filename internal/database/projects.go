package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProjectResponsibility struct {
	Id             string    `json:"id"`
	Responsibility string    `json:"responsibility"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Project struct {
	Id               string                  `json:"id"`
	ResumeId         string                  `json:"resumeId"`
	Name             string                  `json:"name"`
	Description      string                  `json:"description"`
	Role             string                  `json:"role"`
	Responsibilities []ProjectResponsibility `json:"responsibilities"`
	CreatedAt        time.Time               `json:"createdAt"`
}

func CreateProject(
	db VersionedDatabase,
	resume *Resume,
	name string,
	description string,
	role string,
) (Project, error) {
	id := uuid.New().String()
	createdAt := time.Now()

	query := `
INSERT INTO projects (
  project_id,
  resume_id,
  name,
  description,
  role,
  created_at
) VALUES (
  ?, ?, ?,
  ?, ?, ?
)
  `

	result, err := db.DB().Exec(query, id, resume.Id, name, description, role, createdAt.Unix())
	if err != nil {
		return Project{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Project{}, err
	} else if rowsAffected != 1 {
		return Project{}, fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	responsibilities := make([]ProjectResponsibility, 0)
	return Project{
		Id:               id,
		ResumeId:         resume.Id,
		Name:             name,
		Description:      description,
		Role:             role,
		CreatedAt:        createdAt,
		Responsibilities: responsibilities,
	}, nil
}

func GetProject(db VersionedDatabase, id string) *Project {
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
WHERE p.project_id = ?
  `

	rows, err := db.DB().Query(query, id)
	if err != nil || rows == nil {
		return nil
	}

	var project *Project
	responsibilities := make([]ProjectResponsibility, 0)

	for rows.Next() {
		p, r, err := rowsToProject(rows)
		if err != nil {
			return nil
		}

		if project == nil {
			project = &p
		}

		if r != nil {
			responsibilities = append(responsibilities, *r)
		}
	}

	if project != nil {
		project.Responsibilities = responsibilities
	}

	return project
}

func CreateProjectResponsibility(
	db VersionedDatabase,
	project *Project,
	responsibility string,
) (ProjectResponsibility, error) {
	id := uuid.New().String()
	createdAt := time.Now()

	query := `
INSERT INTO project_responsibilities (
  project_responsibility_id,
  project_id,
  responsibility,
  created_at
) VALUES (
  ?, ?, ?,
  ?
)
  `

	result, err := db.DB().Exec(query, id, project.Id, responsibility, createdAt.Unix())
	if err != nil {
		return ProjectResponsibility{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ProjectResponsibility{}, err
	} else if rowsAffected != 1 {
		return ProjectResponsibility{}, fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	r := ProjectResponsibility{
		Id:             id,
		Responsibility: responsibility,
		CreatedAt:      createdAt,
	}

	project.Responsibilities = append(project.Responsibilities, r)

	return r, nil
}

func rowsToProject(rows *sql.Rows) (Project, *ProjectResponsibility, error) {
	var p struct {
		Id          string
		ResumeId    string
		Name        string
		Description string
		Role        string
		CreatedAt   int64
	}

	var r struct {
		Id             *string
		Responsibility *string
		CreatedAt      *int64
	}

	if err := rows.Scan(
		&p.Id,
		&p.ResumeId,
		&p.Name,
		&p.Description,
		&p.Role,
		&p.CreatedAt,
		&r.Id,
		&r.Responsibility,
		&r.CreatedAt,
	); err != nil {
		return Project{}, nil, err
	}

	project := Project{
		Id:          p.Id,
		ResumeId:    p.ResumeId,
		Name:        p.Name,
		Description: p.Description,
		Role:        p.Role,
		CreatedAt:   time.Unix(p.CreatedAt, 0),
	}

	if r.Id == nil {
		return project, nil, nil
	}

	createdAt := time.Unix(*r.CreatedAt, 0)

	responsibility := ProjectResponsibility{
		Id:             *r.Id,
		Responsibility: *r.Responsibility,
		CreatedAt:      createdAt,
	}

	return project, &responsibility, nil
}
