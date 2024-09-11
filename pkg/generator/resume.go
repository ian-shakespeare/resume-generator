package generator

type resumeData struct {
	Name            string
	Role            string
	Email           string
	PhoneNumber     string
	Location        string
	LinkedIn        string
	Github          string
	Facebook        string
	Instagram       string
	Twitter         string
	Portfolio       string
	Educations      []education
	WorkExperiences []workExperience
	Projects        []project
}

type resumeOption func(r *resumeData)

func NewResumeData(personName string, title string, email string, phoneNumber string, opts ...resumeOption) *resumeData {
	r := &resumeData{
		Name:            personName,
		Role:            title,
		Email:           email,
		PhoneNumber:     phoneNumber,
		Educations:      []education{},
		WorkExperiences: []workExperience{},
		Projects:        []project{},
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func WithLocation(location string) resumeOption {
	return func(r *resumeData) {
		r.Location = location
	}
}

func WithLinkedIn(linkedInUsername string) resumeOption {
	return func(r *resumeData) {
		r.LinkedIn = linkedInUsername
	}
}

func WithGitHub(githubUsername string) resumeOption {
	return func(r *resumeData) {
		r.Github = githubUsername
	}
}

func WithFacebook(facebookUsername string) resumeOption {
	return func(r *resumeData) {
		r.Facebook = facebookUsername
	}
}

func WithInstagram(instagramHandle string) resumeOption {
	return func(r *resumeData) {
		r.Instagram = instagramHandle
	}
}

func WithTwitter(twitterHandle string) resumeOption {
	return func(r *resumeData) {
		r.Twitter = twitterHandle
	}
}

func WithPortfolio(portfolioUrl string) resumeOption {
	return func(r *resumeData) {
		r.Portfolio = portfolioUrl
	}
}

func WithEducation(e *education) resumeOption {
	return func(r *resumeData) {
		r.Educations = append(r.Educations, *e)
	}
}

func WithWorkExperience(w *workExperience) resumeOption {
	return func(r *resumeData) {
		r.WorkExperiences = append(r.WorkExperiences, *w)
	}
}

func WithProject(p *project) resumeOption {
	return func(r *resumeData) {
		r.Projects = append(r.Projects, *p)
	}
}

func (r *resumeData) AddEducation(e *education) {
	r.Educations = append(r.Educations, *e)
}

func (r *resumeData) AddWorkExperience(w *workExperience) {
	r.WorkExperiences = append(r.WorkExperiences, *w)
}

func (r *resumeData) AddProject(p *project) {
	r.Projects = append(r.Projects, *p)
}
