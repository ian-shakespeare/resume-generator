package generator_test

import (
	"resumegenerator/pkg/generator"
	"testing"
	"time"
)

const TEST_TEMPLATE string = `<html>
<p>{{.Name}}{{.Role}}{{.Email}}{{.PhoneNumber}}{{.Location}}{{.LinkedIn}}{{.Github}}{{.Facebook}}{{.Instagram}}{{.Twitter}}{{.Portfolio}}</p>
{{range .Educations}}<section>{{.Degree}}{{.FieldOfStudy}}{{.Institution}}{{.Timeframe.Began}}{{if .Timeframe.Current}}current{{end}}{{.Timeframe.Finished}}{{.Location}}{{.GPA}}</section>{{end}}
{{range .WorkExperiences}}<section>{{.Employer}}{{.Role}}{{.Timeframe.Began}}{{if .Timeframe.Current}}current{{end}}{{.Timeframe.Finished}}{{.Location}}<ul>{{range .Responsibilities}}<li>{{.}}</li>{{end}}</ul></section>{{end}}
{{range .Projects}}<section>{{.Name}}{{.Role}}<ul>{{range .Responsibilities}}<li>{{.}}</li>{{end}}</ul></section>{{end}}
</html>`

func TestGenerateHtml(t *testing.T) {
	t.Run("minimumRequredFields", func(t *testing.T) {
		data := generator.NewResumeData("name", "role", "email", "phoneNumber")

		expected := `<html>
<p>nameroleemailphoneNumber</p>



</html>`

		received, err := data.GenerateHtml(TEST_TEMPLATE)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if expected != received {
			t.Fatalf("expected %s, received %s", expected, received)
		}
	})

	t.Run("allFields", func(t *testing.T) {
		d, err := time.Parse(time.RFC3339, "1970-01-01T00:00:00.000Z")
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		data := generator.NewResumeData(
			"name",
			"role",
			"email",
			"phoneNumber",
			generator.WithLocation("location"),
			generator.WithLinkedIn("linkedIn"),
			generator.WithGitHub("github"),
			generator.WithFacebook("facebook"),
			generator.WithInstagram("instagram"),
			generator.WithTwitter("twitter"),
			generator.WithPortfolio("portfolio"),
			generator.WithEducation(
				generator.NewEducation(
					"degree",
					"fieldOfStudy",
					"institution",
					generator.CurrentDateRange(d),
					generator.WithInstitutionLocation("location"),
					generator.WithGPA("gpa"),
				),
			),
			generator.WithWorkExperience(
				generator.NewWorkExperience(
					"employer",
					"role",
					"location",
					generator.CurrentDateRange(d),
					"responsibility",
				),
			),
			generator.WithProject(
				generator.NewProject(
					"name",
					"role",
					"responsibility",
				),
			),
		)

		expected := `<html>
<p>nameroleemailphoneNumberlocationlinkedIngithubfacebookinstagramtwitterportfolio</p>
<section>degreefieldOfStudyinstitutionJanuary 1970currentlocationgpa</section>
<section>employerroleJanuary 1970currentlocation<ul><li>responsibility</li></ul></section>
<section>namerole<ul><li>responsibility</li></ul></section>
</html>`

		received, err := data.GenerateHtml(TEST_TEMPLATE)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if expected != received {
			t.Fatalf("expected %s, received %s", expected, received)
		}
	})
}
