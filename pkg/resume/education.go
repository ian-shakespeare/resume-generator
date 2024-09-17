package resume

import "time"

type Education struct {
	Degree       string     `json,yaml:"degree"`
	FieldOfStudy string     `json,yaml:"fieldofstudy"`
	Institution  string     `json,yaml:"institution"`
	Began        time.Time  `json,yaml:"began"`
	Current      bool       `json,yaml:"current"`
	Finished     *time.Time `json,yaml:"finished"`
	Location     string     `json,yaml:"location"`
	GPA          string     `json,yaml:"gpa"`
}
