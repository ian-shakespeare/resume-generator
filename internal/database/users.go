package database

import (
	"fmt"
	"resumegenerator/pkg/resume"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        string
	CreatedAt time.Time
	Resumes   []Resume
}

func CreateUser(c Connection) (User, error) {
	query := `
  INSERT INTO users (user_id, created_at)
  VALUES (?, ?)
  `

	id := uuid.New().String()
	createdAt := time.Now()

	result, err := c.Database().Exec(query, id, createdAt.Unix())
	if err != nil {
		return User{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return User{}, err
	}

	if rowsAffected != 1 {
		return User{}, fmt.Errorf("affected an unexpected number of rows (%d)", rowsAffected)
	}

	return User{
		Id:        id,
		CreatedAt: createdAt,
		Resumes:   []Resume{},
	}, nil
}

func GetUser(c Connection, id string) (User, error) {
	query := `
  SELECT
    user_id
    , created_at
  FROM users
  WHERE user_id = ?
  `

	row := c.Database().QueryRow(query, id)

	var u struct {
		id        string
		createdAt int64
	}
	if err := row.Scan(&u.id, &u.createdAt); err != nil {
		return User{}, err
	}

	return User{
		Id:        u.id,
		CreatedAt: time.Unix(u.createdAt, 0),
	}, nil
}

func (u *User) CreateResume(c Connection, r resume.Resume) (Resume, error) {
	return Resume{}, nil
}
