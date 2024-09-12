package resume

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type WorkExperience struct {
	Id               string     `json:"id"`
	CreatedAt        time.Time  `json:"createdAt"`
	Employer         string     `json:"employer"`
	Title            string     `json:"title"`
	Began            time.Time  `json:"began"`
	Current          bool       `json:"current"`
	Finished         *time.Time `json:"finished"`
	Responsibilities []string   `json:"responsibilities"`
	Location         string     `json:"location"`
}

func WorkExperienceFromJSON(b []byte) (WorkExperience, error) {
	var w struct {
		Employer         string   `json:"employer"`
		Title            string   `json:"title"`
		Began            string   `json:"began"`
		Current          bool     `json:"current"`
		Responsibilities []string `json:"responsibilities"`
		Location         string   `json:"location"`
		Finished         string   `json:"finished"`
	}
	if err := json.Unmarshal(b, &w); err != nil {
		return WorkExperience{}, err
	}

	id := uuid.New().String()
	now := time.Now()

	began, err := time.Parse(time.RFC3339, w.Began)
	if err != nil {
		return WorkExperience{}, err
	}

	var finished *time.Time
	if w.Finished != "" {
		f, err := time.Parse(time.RFC3339, w.Finished)
		if err != nil {
			return WorkExperience{}, err
		}
		finished = &f
	}

	return WorkExperience{
		Id:               id,
		CreatedAt:        now,
		Employer:         w.Employer,
		Title:            w.Title,
		Began:            began,
		Current:          w.Current,
		Finished:         finished,
		Responsibilities: w.Responsibilities,
		Location:         w.Location,
	}, nil
}
