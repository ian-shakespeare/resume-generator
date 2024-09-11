package generator

type project struct {
	Name             string
	Role             string
	Responsibilities []string
}

func ProjectData(name string, role string, responsibilities ...string) *project {
	return &project{
		name,
		role,
		responsibilities,
	}
}
