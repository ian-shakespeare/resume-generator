package resume

import "time"

type WorkExperience struct {
	Employer         string     `json,yaml:"employer"`
	Title            string     `json,yaml:"title"`
	Began            time.Time  `json,yaml:"began"`
	Current          bool       `json,yaml:"current"`
	Finished         *time.Time `json,yaml:"finished"`
	Responsibilities []string   `json,yaml:"responsibilities"`
	Location         string     `json,yaml:"location"`
}
