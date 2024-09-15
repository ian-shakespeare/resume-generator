package resume

import "time"

type Resume struct {
	Name            string           `json,yaml:"name"`
	Title           string           `json,yaml:"title"`
	Email           string           `json,yaml:"email"`
	PhoneNumber     string           `json,yaml:"phonenumber"`
	Prelude         string           `json,yaml:"prelude"`
	Location        string           `json,yaml:"location"`
	LinkedIn        string           `json,yaml:"linkedin"`
	Github          string           `json,yaml:"github"`
	Facebook        string           `json,yaml:"facebook"`
	Instagram       string           `json,yaml:"instagram"`
	Twitter         string           `json,yaml:"twitter"`
	Portfolio       string           `json,yaml:"portfolio"`
	Skills          []string         `json,yaml:"skills"`
	Educations      []Education      `json,yaml:"educations"`
	WorkExperiences []WorkExperience `json,yaml:"workexperiences"`
	Projects        []Project        `json,yaml:"projects"`
}

func MinExample() Resume {
	return Resume{
		Name:            "John Doe",
		Title:           "Job Title",
		Email:           "jdoe@email.com",
		PhoneNumber:     "+1 (000) 000-0000",
		Prelude:         "Lorem ipsum odor amet, consectetuer adipiscing elit.",
		Skills:          []string{},
		Educations:      []Education{},
		WorkExperiences: []WorkExperience{},
		Projects:        []Project{},
	}
}

func Example() Resume {
  t := time.Unix(0, 0)
	return Resume{
		Name:        "John Doe",
		Title:       "Job Title",
		Email:       "jdoe@email.com",
		PhoneNumber: "+1 (000) 000-0000",
		Prelude:     "Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis. Tortor eleifend nec eleifend curabitur, urna non. Aliquet dictumst sapien himenaeos habitasse taciti fames accumsan porta suscipit. Dolor consequat venenatis fusce pellentesque diam blandit a senectus.",
		Location:    "City, State",
		LinkedIn:    "linked-in-account",
		Github:      "github-account",
		Portfolio:   "portfolio-website",
		Skills: []string{
			"Lorem ipsum odor amet",
			"Lorem ipsum odor amet",
			"Lorem ipsum odor amet",
			"Lorem ipsum odor amet",
			"Lorem ipsum odor amet",
			"Lorem ipsum odor amet",
		},
		Educations: []Education{
			{Degree: "Degree", FieldOfStudy: "Field", Institution: "Institution", Began: t, Current: true, Location: "City, State"},
			{Degree: "Earlier Degree", FieldOfStudy: "Field", Institution: "Institution", Began: t, Current: false, Finished: &t, Location: "City, State"},
		},
		WorkExperiences: []WorkExperience{
			{
				Employer: "Employer",
				Title:    "Job Title",
				Began:    t,
				Current:  true,
				Finished: &t,
        Location: "City, State",
				Responsibilities: []string{
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
				},
			},
			{
				Employer: "Earlier Employer",
				Title:    "Job Title",
				Began:    t,
				Current:  false,
				Finished: &t,
        Location: "City, State",
				Responsibilities: []string{
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
				},
			},
			{
				Employer: "Earliest Employer",
				Title:    "Job Title",
				Began:    t,
				Current:  false,
				Finished: &t,
        Location: "City, State",
				Responsibilities: []string{
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
				},
			},
		},
		Projects: []Project{
			{
				Name:        "Project Name",
				Role:        "Role",
				Description: "Lorem ipsum odor amet, consectetuer adipiscing elit.",
				Responsibilities: []string{
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
				},
			},
			{
				Name:        "Project Name",
				Role:        "Role",
				Description: "Lorem ipsum odor amet, consectetuer adipiscing elit.",
				Responsibilities: []string{
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
				},
			},
			{
				Name:        "Project Name",
				Role:        "Role",
				Description: "Lorem ipsum odor amet, consectetuer adipiscing elit.",
				Responsibilities: []string{
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
					"Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.",
				},
			},
		},
	}
}
