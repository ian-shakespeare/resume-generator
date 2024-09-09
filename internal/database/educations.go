package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Education struct {
	Id           string     `json:"id"`
	ResumeId     string     `json:"resumeId"`
	DegreeType   string     `json:"degreeType"`
	FieldOfStudy string     `json:"fieldOfStudy"`
	Institution  string     `json:"institution"`
	Began        time.Time  `json:"began"`
	Current      bool       `json:"current"`
	CreatedAt    time.Time  `json:"createdAt"`
	Location     *string    `json:"location"`
	Finished     *time.Time `json:"finished"`
	GPA          *string    `json:"gpa"`
}

func CreateEducation(
	db VersionedDatabase,
	resume *Resume,
	degreeType string,
	fieldOfStudy string,
	institution string,
	began time.Time,
	current bool,
	location *string,
	finished *time.Time,
	gpa *string,
) (Education, error) {
	id := uuid.New().String()
	createdAt := time.Now()

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
	if current {
		currentInt = 1
	}

	var finishedInt *int64
	if finished != nil {
		val := finished.Unix()
		finishedInt = &val
	}

	result, err := db.DB().Exec(
		query,
		id,
		resume.Id,
		degreeType,
		fieldOfStudy,
		institution,
		began.Unix(),
		currentInt,
		createdAt.Unix(),
		location,
		finishedInt,
		gpa,
	)

	if err != nil {
		return Education{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Education{}, err
	} else if rowsAffected != 1 {
		return Education{}, fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	return Education{
		Id:           id,
		ResumeId:     resume.Id,
		DegreeType:   degreeType,
		FieldOfStudy: fieldOfStudy,
		Institution:  institution,
		Began:        began,
		Current:      current,
		CreatedAt:    createdAt,
		Location:     location,
		Finished:     finished,
		GPA:          gpa,
	}, nil
}

func GetEducation(db VersionedDatabase, id string) *Education {
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

	row := db.DB().QueryRow(query, id)
	if row == nil {
		return nil
	}

	education, err := rowToEducation(row)
	if err != nil {
		return nil
	}

	return &education
}

func rowToEducation(row *sql.Row) (Education, error) {
	var education struct {
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
	if err := row.Scan(
		&education.Id,
		&education.ResumeId,
		&education.DegreeType,
		&education.FieldOfStudy,
		&education.Institution,
		&education.Began,
		&education.Current,
		&education.CreatedAt,
		&education.Location,
		&education.Finished,
		&education.GPA,
	); err != nil {
		return Education{}, nil
	}

	began := time.Unix(education.Began, 0)
	current := false
	if education.Current != 0 {
		current = true
	}
	createdAt := time.Unix(education.CreatedAt, 0)
	var finished *time.Time
	if education.Finished != nil {
		val := time.Unix(*education.Finished, 0)
		finished = &val
	}

	return Education{
		Id:           education.Id,
		ResumeId:     education.ResumeId,
		DegreeType:   education.DegreeType,
		FieldOfStudy: education.FieldOfStudy,
		Institution:  education.Institution,
		Began:        began,
		Current:      current,
		CreatedAt:    createdAt,
		Location:     education.Location,
		Finished:     finished,
		GPA:          education.GPA,
	}, nil
}
