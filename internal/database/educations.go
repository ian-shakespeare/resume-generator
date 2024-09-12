package database

import (
	"database/sql"
	"fmt"
	"resumegenerator/pkg/resume"
	"time"
)

func CreateEducation(
	db VersionedDatabase,
	r *resume.Resume,
	e *resume.Education,
) error {
	query := `
INSERT INTO educations (
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
) VALUES (
  ?, ?, ?,
  ?, ?, ?,
  ?, ?, ?,
  ?, ?
)
  `

	currentInt := 0
	if e.Current {
		currentInt = 1
	}

	result, err := db.DB().Exec(
		query,
		e.Id,
		r.Id,
		e.DegreeType,
		e.FieldOfStudy,
		e.Institution,
		e.Began.Unix(),
		currentInt,
		e.CreatedAt.Unix(),
		strToStrPtr(e.Location),
		timePtrToIntPtr(e.Finished),
		strToStrPtr(e.GPA),
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

func GetEducation(db VersionedDatabase, id string) *resume.Education {
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
WHERE education_id = ?
  `

	rows, err := db.DB().Query(query, id)
	if err != nil {
		return nil
	}

	e, err := rowsToEducation(rows)
	if err != nil {
		return nil
	}

	if len(e) < 1 {
		return nil
	}

	return &e[0]
}

func rowsToEducation(rows *sql.Rows) ([]resume.Education, error) {
	educations := make([]resume.Education, 0)

	for rows.Next() {
		var stored struct {
			Id           string
			ResumeId     string
			DegreeType   string
			FieldOfStudy string
			Institution  string
			Began        int64
			Current      int
			CreatedAt    int64
			Location     *string
			Finished     *int64
			GPA          *string
		}
		if err := rows.Scan(
			&stored.Id,
			&stored.ResumeId,
			&stored.DegreeType,
			&stored.FieldOfStudy,
			&stored.Institution,
			&stored.Began,
			&stored.Current,
			&stored.CreatedAt,
			&stored.Location,
			&stored.Finished,
			&stored.GPA,
		); err != nil {
			return nil, err
		}

		// TODO: Resume id

		began := time.Unix(stored.Began, 0)
		current := false
		if stored.Current != 0 {
			current = true
		}
		createdAt := time.Unix(stored.CreatedAt, 0)
		var finished *time.Time
		if stored.Finished != nil {
			val := time.Unix(*stored.Finished, 0)
			finished = &val
		}
		educations = append(educations, resume.Education{
			Id:           stored.Id,
			DegreeType:   stored.DegreeType,
			FieldOfStudy: stored.FieldOfStudy,
			Institution:  stored.Institution,
			Began:        began,
			Current:      current,
			CreatedAt:    createdAt,
			Location:     strPtrToStr(stored.Location),
			Finished:     finished,
			GPA:          strPtrToStr(stored.GPA),
		})
	}

	return educations, nil
}
