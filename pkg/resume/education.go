package resume

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Education struct {
	Id           string     `json:"id"`
	CreatedAt    time.Time  `json:"createdAt"`
	DegreeType   string     `json:"degreeType"`
	FieldOfStudy string     `json:"fieldOfStudy"`
	Institution  string     `json:"institution"`
	Began        time.Time  `json:"began"`
	Current      bool       `json:"current"`
	Finished     *time.Time `json:"finished"`
	Location     string     `json:"location"`
	GPA          string     `json:"gpa"`
}

func NewEducation(b []byte) (Education, error) {
	var e struct {
		DegreeType   string `json:"degreeType"`
		FieldOfStudy string `json:"fieldOfStudy"`
		Institution  string `json:"institution"`
		Began        string `json:"began"`
		Current      bool   `json:"current"`
		Location     string `json:"location"`
		Finished     string `json:"finished"`
		GPA          string `json:"gpa"`
	}
	if err := json.Unmarshal(b, &e); err != nil {
		return Education{}, err
	}

	id := uuid.New().String()
	now := time.Now()

	began, err := time.Parse(time.RFC3339, e.Began)
	if err != nil {
		return Education{}, err
	}

	var finished *time.Time
	if e.Finished != "" {
		f, err := time.Parse(time.RFC3339, e.Finished)
		if err != nil {
			return Education{}, err
		}
		finished = &f
	}

	return Education{
		Id:           id,
		CreatedAt:    now,
		DegreeType:   e.DegreeType,
		FieldOfStudy: e.FieldOfStudy,
		Institution:  e.Institution,
		Began:        began,
		Current:      e.Current,
		Finished:     finished,
		Location:     e.Location,
		GPA:          e.GPA,
	}, nil
}
