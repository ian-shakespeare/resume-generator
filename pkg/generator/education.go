package generator

type education struct {
	Degree       string
	FieldOfStudy string
	Institution  string
	Timeframe    dateRange
	Location     string
	GPA          string
}

type educationOption func(e *education)

func WithInstitutionLocation(location string) educationOption {
	return func(e *education) {
		e.Location = location
	}
}

func WithGPA(gpa string) educationOption {
	return func(e *education) {
		e.GPA = gpa
	}
}

func EducationData(degree string, fieldOfStudy string, institution string, timeframe dateRange, opts ...educationOption) *education {
	e := &education{
		Degree:       degree,
		FieldOfStudy: fieldOfStudy,
		Institution:  institution,
		Timeframe:    timeframe,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}
