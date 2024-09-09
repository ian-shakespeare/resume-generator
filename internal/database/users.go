package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func CreateUser(db VersionedDatabase) (User, error) {
	id := uuid.New().String()
	createdAt := time.Now()

	query := `
INSERT INTO users (user_id, created_at)
VALUES (?, ?)
  `

	result, err := db.DB().Exec(query, id, createdAt.Unix())
	if err != nil {
		return User{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return User{}, err
	} else if rowsAffected != 1 {
		return User{}, errors.New(fmt.Sprintf("affected an unexpected number of rows (%d)", rowsAffected))
	}

	return User{Id: id, CreatedAt: createdAt}, nil
}

func GetUser(db VersionedDatabase, id string) *User {
	query := `
SELECT 
  user_id,
  created_at
FROM users
WHERE user_id = ?
  `

	row := db.DB().QueryRow(query, id)
	if row == nil {
		return nil
	}

	user, err := rowToUser(row)
	if err != nil {
		return nil
	}

	return &user
}

func rowToUser(row *sql.Row) (User, error) {
	var user struct {
		Id        string
		CreatedAt int64
	}
	if err := row.Scan(&user.Id, &user.CreatedAt); err != nil {
		return User{}, err
	}

	return User{
		Id:        user.Id,
		CreatedAt: time.Unix(user.CreatedAt, 0),
	}, nil
}
