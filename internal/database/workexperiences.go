package database

import (
	"fmt"
	"resumegenerator/pkg/resume"
	"strings"
	"time"

	"github.com/google/uuid"
)

type WorkExperience struct {
	resume.WorkExperience
	Id        string
	ResumeId  string
	CreatedAt time.Time
}

func (w *WorkExperience) AddResponsibility(c Connection, responsibilities ...string) error {
	row := c.Database().QueryRow(`SELECT MAX(position) FROM work_responsibilities WHERE work_experience_id = ?`, w.Id)
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
		args = append(args, w.Id)
		args = append(args, startingPos+i)
		args = append(args, responsibility)
	}

	query := fmt.Sprintf(`
  INSERT INTO work_responsibilities (
    work_responsibility_id,
    work_experience_id,
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

	w.Responsibilities = append(w.Responsibilities, responsibilities...)

	return nil
}
