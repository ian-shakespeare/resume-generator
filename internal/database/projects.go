package database

import (
	"database/sql"
	"fmt"
	"resumegenerator/pkg/utils"
	"resumegenerator/pkg/resume"
	"time"
)

func CreateProject(
	db VersionedDatabase,
	r *resume.Resume,
	p *resume.Project,
) error {
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

	result, err := db.DB().Exec(
		query,
		p.Id,
		r.Id,
		p.Name,
		p.Description,
		p.Role,
		p.CreatedAt.Unix(),
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

func GetProject(db VersionedDatabase, id string) *resume.Project {
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

	projects, err := rowsToProject(rows)
	if err != nil || len(projects) != 1 {
		return nil
	}

	return &projects[0]
}

func rowsToProject(rows *sql.Rows) ([]resume.Project, error) {
	projects := make([]resume.Project, 0)

	for rows.Next() {
		var sp struct {
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
			&sp.Id,
			&sp.ResumeId,
			&sp.Name,
			&sp.Description,
			&sp.Role,
			&sp.CreatedAt,
			&r.Id,
			&r.Responsibility,
			&r.CreatedAt,
		); err != nil {
			return nil, err
		}

		existingIndex := utils.Find(projects, func(p resume.Project) bool {
			return sp.Id == p.Id
		})

		if existingIndex == -1 {
			p := resume.Project{
				Id:          sp.Id,
				Name:        sp.Name,
				Description: sp.Description,
				Role:        sp.Role,
				CreatedAt:   time.Unix(sp.CreatedAt, 0),
			}
			projects = append(projects, p)
			existingIndex = len(projects) - 1
		}

		if r.Id != nil {
			projects[existingIndex].Responsibilities = append(projects[existingIndex].Responsibilities, *r.Responsibility)
		}
	}

	return projects, nil
}
