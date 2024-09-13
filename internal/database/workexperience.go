package database

import (
	"database/sql"
	"fmt"
	"resumegenerator/pkg/utils"
	"resumegenerator/pkg/resume"
	"time"
)

func CreateWorkExperience(
	db VersionedDatabase,
	r *resume.Resume,
	w *resume.WorkExperience,
) error {
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

	// TODO: Create responsibilities
	result, err := db.DB().Exec(
		query,
		w.Id,
		r.Id,
		w.Employer,
		w.Title,
		w.Began.Unix(),
		w.Current,
		w.CreatedAt.Unix(),
		w.Location,
		timePtrToIntPtr(w.Finished),
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

func GetWorkExperience(db VersionedDatabase, id string) *resume.WorkExperience {
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

	workExperiences, err := rowsToWorkExperience(rows)
	if err != nil || len(workExperiences) != 1 {
		return nil
	}

	return &workExperiences[0]
}

func rowsToWorkExperience(rows *sql.Rows) ([]resume.WorkExperience, error) {
	experiences := make([]resume.WorkExperience, 0)

	for rows.Next() {
		var sw struct {
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
			&sw.Id,
			&sw.ResumeId,
			&sw.Employer,
			&sw.Title,
			&sw.Began,
			&sw.Current,
			&sw.CreatedAt,
			&sw.Location,
			&sw.Finished,
			&r.Id,
			&r.Responsibility,
			&r.CreatedAt,
		); err != nil {
			return nil, err
		}

		existingIndex := utils.Find(experiences, func(w resume.WorkExperience) bool {
			return sw.Id == w.Id
		})

		if existingIndex == -1 {
			current := false
			if sw.Current != 0 {
				current = true
			}
			var finished *time.Time
			if sw.Finished != nil {
				val := time.Unix(*sw.Finished, 0)
				finished = &val
			}

			w := resume.WorkExperience{
				Id:        sw.Id,
				Employer:  sw.Employer,
				Title:     sw.Title,
				Began:     time.Unix(sw.Began, 0),
				Current:   current,
				CreatedAt: time.Unix(sw.CreatedAt, 0),
				Location:  *sw.Location,
				Finished:  finished,
			}

			experiences = append(experiences, w)
			existingIndex = len(experiences) - 1
		}

		if r.Id != nil {
			experiences[existingIndex].Responsibilities = append(experiences[existingIndex].Responsibilities, *r.Responsibility)
		}
	}

	return experiences, nil
}
