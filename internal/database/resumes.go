package database

import (
	"resumegenerator/pkg/resume"
	"time"
)

type Resume struct {
	resume.Resume

	Id        string
  UserId string
	CreatedAt time.Time
}

func GetResume(c Connection, id string) (Resume, error) {
  return Resume{}, nil
}
