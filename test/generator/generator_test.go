package generator_test

import (
	"resumegenerator/pkg/generator"
	"resumegenerator/pkg/resume"
	"testing"
)

const TEST_TEMPLATE string = `<html>
<p>{{.Name}}{{.Email}}{{.PhoneNumber}}{{.Location}}{{.LinkedIn}}{{.Github}}{{.Facebook}}{{.Instagram}}{{.Twitter}}{{.Portfolio}}</p>
<ul>{{range .Skills}}<li>{{.}}</li>{{end}}</ul>
<ul>{{range .Achievements}}<li>{{.}}</li>{{end}}</ul>
{{range .Educations}}<section>{{.Degree}}{{.FieldOfStudy}}{{.Institution}}{{month_year .Began}}{{if .Current}}current{{end}}{{month_year .Finished}}{{.Location}}{{.GPA}}</section>{{end}}
{{range .WorkExperiences}}<section>{{.Employer}}{{.Title}}{{month_year .Began}}{{if .Current}}current{{end}}{{month_year .Finished}}{{.Location}}<ul>{{range .Responsibilities}}<li>{{.}}</li>{{end}}</ul></section>{{end}}
{{range .Projects}}<section>{{.Name}}{{.Role}}<ul>{{range .Responsibilities}}<li>{{.}}</li>{{end}}</ul></section>{{end}}
</html>`

func TestGenerateHtml(t *testing.T) {
	t.Run("minimumRequredFields", func(t *testing.T) {
		expected := `<html>
<p>John Doejdoe@email.com&#43;1 (000) 000-0000</p>
<ul></ul>
<ul></ul>



</html>`

		r := resume.MinExample()

		received, err := generator.GenerateHtml(&r, []byte(TEST_TEMPLATE))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if expected != string(received) {
			t.Fatalf("expected %s, received %s", expected, string(received))
		}
	})

	t.Run("allFields", func(t *testing.T) {
		r := resume.Example()

		expected := `<html>
<p>John Doejdoe@email.com&#43;1 (000) 000-0000City, Statelinked-in-accountgithub-accountportfolio-website</p>
<ul><li>Lorem ipsum odor amet, consectetuer adipiscing elit</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit</li><li>Lorem ipsum odor amet</li><li>Lorem ipsum odor amet</li><li>Lorem ipsum odor amet</li></ul>
<ul><li>Lorem ipsum odor amet, consectetuer adipiscing elit</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit</li><li>Lorem ipsum odor amet</li><li>Lorem ipsum odor amet</li><li>Lorem ipsum odor amet</li></ul>
<section>DegreeFieldInstitutionJanuary 1970currentCity, State</section><section>Earlier DegreeFieldInstitutionJanuary 1970January 1970City, State</section>
<section>EmployerJob TitleJanuary 1970currentJanuary 1970City, State<ul><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li></ul></section><section>Earlier EmployerJob TitleJanuary 1970January 1970City, State<ul><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li></ul></section><section>Earliest EmployerJob TitleJanuary 1970January 1970City, State<ul><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li></ul></section>
<section>Project NameRole<ul><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li></ul></section><section>Project NameRole<ul><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li></ul></section><section>Project NameRole<ul><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li><li>Lorem ipsum odor amet, consectetuer adipiscing elit. Venenatis mi dignissim sem quisque iaculis.</li></ul></section>
</html>`

		received, err := generator.GenerateHtml(&r, []byte(TEST_TEMPLATE))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if expected != string(received) {
			t.Fatalf("expected %s, received %s", expected, string(received))
		}
	})
}
