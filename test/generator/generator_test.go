package generator_test

import (
	"resumegenerator/pkg/generator"
	"resumegenerator/pkg/resume"
	"resumegenerator/test"
	"testing"
)

const TEST_TEMPLATE string = `<html>
<p>{{.Name}}{{.Email}}{{.PhoneNumber}}{{.Location}}{{.LinkedIn}}{{.Github}}{{.Facebook}}{{.Instagram}}{{.Twitter}}{{.Portfolio}}</p>
{{range .Educations}}<section>{{.DegreeType}}{{.FieldOfStudy}}{{.Institution}}{{.Timeframe.Began}}{{if .Timeframe.Current}}current{{end}}{{.Timeframe.Finished}}{{.Location}}{{.GPA}}</section>{{end}}
{{range .WorkExperiences}}<section>{{.Employer}}{{.Role}}{{.Timeframe.Began}}{{if .Timeframe.Current}}current{{end}}{{.Timeframe.Finished}}{{.Location}}<ul>{{range .Responsibilities}}<li>{{.}}</li>{{end}}</ul></section>{{end}}
{{range .Projects}}<section>{{.Name}}{{.Role}}<ul>{{range .Responsibilities}}<li>{{.}}</li>{{end}}</ul></section>{{end}}
</html>`

func TestGenerateHtml(t *testing.T) {
	t.Run("minimumRequredFields", func(t *testing.T) {
		expected := `<html>
<p>nameemailphoneNumber</p>



</html>`

		r, err := resume.FromJSON([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		received, err := generator.GenerateHtml(&r, TEST_TEMPLATE)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if expected != received {
			t.Fatalf("expected %s, received %s", expected, received)
		}
	})

	t.Run("allFields", func(t *testing.T) {
		r, err := resume.FromJSON([]byte(test.FULL_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		e, err := resume.EducationFromJSON([]byte(test.FULL_EDUCATION))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		w, err := resume.WorkExperienceFromJSON([]byte(test.FULL_WORK_EXPERIENCE))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		p, err := resume.ProjectFromJSON([]byte(test.PROJECT))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r.AddEducation(e)
		r.AddWorkExperiences(w)
		r.AddProject(p)

		expected := `<html>
<p>nameemailphoneNumberlocationlinkedIngithubfacebookinstagramtwitterportfolio</p>
<section>degreefieldOfStudyinstitutionJanuary 1970currentlocationgpa</section>
<section>employerroleJanuary 1970currentlocation<ul><li>responsibility</li></ul></section>
<section>namerole<ul><li>responsibility</li></ul></section>
</html>`

		received, err := generator.GenerateHtml(&r, TEST_TEMPLATE)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if expected != received {
			t.Fatalf("expected %s, received %s", expected, received)
		}
	})
}
