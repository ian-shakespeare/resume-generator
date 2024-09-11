package generator

type workExperience struct {
	Employer         string
	Role             string
	Timeframe        dateRange
	Location         string
	Responsibilities []string
}

func WorkExperienceData(employer string, role string, location string, timeframe dateRange, responsibilities ...string) *workExperience {
	return &workExperience{
		Employer:         employer,
		Role:             role,
		Timeframe:        timeframe,
		Location:         location,
		Responsibilities: responsibilities,
	}
}
