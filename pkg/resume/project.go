package resume

type Project struct {
	Name             string    `json,yaml:"name"`
	Description      string    `json,yaml:"description"`
	Role             string    `json,yaml:"role"`
	Responsibilities []string  `json,yaml:"responsibilities"`
}
