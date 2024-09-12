package resume

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	Id               string    `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Role             string    `json:"role"`
	Responsibilities []string  `json:"responsibilities"`
}

func NewProject(b []byte) (Project, error) {
	var p struct {
		Name             string   `json:"name"`
		Description      string   `json:"description"`
		Role             string   `json:"role"`
		Responsibilities []string `json:"responsibilities"`
	}
	if err := json.Unmarshal(b, &p); err != nil {
		return Project{}, err
	}

	id := uuid.New().String()
	now := time.Now()

	return Project{
		Id:               id,
		CreatedAt:        now,
		Name:             p.Name,
		Description:      p.Description,
		Role:             p.Role,
		Responsibilities: p.Responsibilities,
	}, nil
}
