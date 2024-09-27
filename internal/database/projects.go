package database

import (
	"fmt"
	"resumegenerator/pkg/resume"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	resume.Project
	Id        string
	ResumeId  string
	CreatedAt time.Time
}

func (p *Project) AddResponsibility(c Connection, responsibilities ...string) error {
	row := c.Database().QueryRow(`SELECT MAX(position) FROM project_responsibilities WHERE project_id = ?`, p.Id)
	var startingPos int
	if err := row.Scan(&startingPos); err != nil {
		startingPos = -1
	}
	startingPos += 1

	queryStrs := make([]string, 0, len(responsibilities))
	var args []interface{}

	for i, responsibility := range responsibilities {
		queryStrs = append(queryStrs, "(?, ?, ?, ?)")
		args = append(args, uuid.NewString())
		args = append(args, p.Id)
		args = append(args, startingPos+i)
		args = append(args, responsibility)
	}

	query := fmt.Sprintf(`
  INSERT INTO project_responsibilities (
    project_responsibility_id,
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

	if rowsAffected != int64(len(responsibilities)) {
		return fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	p.Responsibilities = append(p.Responsibilities, responsibilities...)

	return nil
}
