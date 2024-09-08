package database

import (
	"database/sql"
	"errors"
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
	CreatedAt   time.Time `json:"createdAt"`
	Location    *string   `json:"location"`
	LinkedIn    *string   `json:"linkedIn"`
	Github      *string   `json:"github"`
	Facebook    *string   `json:"facebook"`
	Instagram   *string   `json:"instagram"`
	Twitter     *string   `json:"twitter"`
	Portfolio   *string   `json:"portfolio"`
}

func resumeFromRow(row *sql.Row) (Resume, error) {
	var resume Resume
	if err := row.Scan(
		&resume.Id,
		&resume.UserId,
		&resume.Name,
		&resume.Email,
		&resume.PhoneNumber,
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

func CreateResume(
	db VersionedDatabase,
	user *User,
	name string,
	email string,
	phoneNumber string,
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
  ?
)
  `

	result, err := db.DB().Exec(query, id, user.Id, name, email, phoneNumber, location, linkedIn, github, facebook, instagram, twitter, portfolio, createdAt.Unix())
	if err != nil {
		return Resume{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Resume{}, err
	} else if rowsAffected != 1 {
		return Resume{}, errors.New(fmt.Sprintf("affected an unexpected number of rows (%d)", rowsAffected))
	}

	return Resume{
		Id:          id,
		UserId:      user.Id,
		Name:        name,
		Email:       email,
		PhoneNumber: phoneNumber,
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

	resume, err := resumeFromRow(row)
	if err != nil {
		return nil
	}

	return &resume
}
