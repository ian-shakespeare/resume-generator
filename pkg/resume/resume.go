package resume

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Resume struct {
	Id              string           `json:"id"`
	CreatedAt       time.Time        `json:"createdAt"`
	Name            string           `json:"name"`
	Email           string           `json:"email"`
	PhoneNumber     string           `json:"phoneNumber"`
	Prelude         string           `json:"prelude"`
	Location        string           `json:"location"`
	LinkedIn        string           `json:"linkedIn"`
	Github          string           `json:"github"`
	Facebook        string           `json:"facebook"`
	Instagram       string           `json:"instagram"`
	Twitter         string           `json:"twitter"`
	Portfolio       string           `json:"portfolio"`
	Educations      []Education      `json:"educations"`
	WorkExperiences []WorkExperience `json:"workExperiences"`
	Projects        []Project        `json:"projects"`
}

func FromJSON(b []byte) (Resume, error) {
	var r struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
		Prelude     string `json:"prelude"`
		Location    string `json:"location"`
		LinkedIn    string `json:"linkedIn"`
		Github      string `json:"github"`
		Facebook    string `json:"facebook"`
		Instagram   string `json:"instagram"`
		Twitter     string `json:"twitter"`
		Portfolio   string `json:"portfolio"`
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return Resume{}, err
	}

	id := uuid.New().String()
	now := time.Now()

	return Resume{
		Id:          id,
		CreatedAt:   now,
		Name:        r.Name,
		Email:       r.Email,
		PhoneNumber: r.PhoneNumber,
		Prelude:     r.Prelude,
		Location:    r.Location,
		LinkedIn:    r.LinkedIn,
		Github:      r.Github,
		Facebook:    r.Facebook,
		Instagram:   r.Instagram,
		Twitter:     r.Twitter,
		Portfolio:   r.Portfolio,
		Educations:  []Education{},
	}, nil
}

func (r *Resume) AddEducation(e Education) {
	r.Educations = append(r.Educations, e)
}

func (r *Resume) AddWorkExperiences(w WorkExperience) {
	r.WorkExperiences = append(r.WorkExperiences, w)
}

func (r *Resume) AddProject(p Project) {
	r.Projects = append(r.Projects, p)
}
