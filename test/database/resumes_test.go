package database_test

import (
	"fmt"
	"resumegenerator/internal/database"
	"resumegenerator/pkg/resume"
	"resumegenerator/test"
	"testing"
)

func TestAddSkill(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddSkill(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM skills WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 0 {
			t.Fatalf("expected %d, received %d", 0, count)
		}

		if len(r.Skills) != 0 {
			t.Fatalf("expected %d, received %d", 0, len(r.Skills))
		}
	})

	t.Run("single", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddSkill(conn, "some skill")
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM skills WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 1 {
			t.Fatalf("expected %d, received %d", 1, count)
		}

		if len(r.Skills) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(r.Skills))
		}
	})

	t.Run("many", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddSkill(conn, "skill0", "skill1", "skill2")
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM skills WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 3 {
			t.Fatalf("expected %d, received %d", 3, count)
		}

		if len(r.Skills) != 3 {
			t.Fatalf("expected %d, received %d", 3, len(r.Skills))
		}
	})

	t.Run("ordered", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		skills := []string{"skill0", "skill1", "skill2", "skill3", "skill4"}

		err = r.AddSkill(conn, skills[0], skills[1], skills[2])
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddSkill(conn, skills[3], skills[4])
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM skills WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 5 {
			t.Fatalf("expected %d, received %d", 5, count)
		}

		if len(r.Skills) != 5 {
			t.Fatalf("expected %d, received %d", 5, len(r.Skills))
		}

		for i, skill := range skills {
			expected := fmt.Sprintf("skill%d", i)
			if skill != expected {
				t.Fatalf("expected %s, received %s", expected, skill)
			}
		}
	})
}

func TestAddAchievement(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddAchievement(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM achievements WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 0 {
			t.Fatalf("expected %d, received %d", 0, count)
		}

		if len(r.Achievements) != 0 {
			t.Fatalf("expected %d, received %d", 0, len(r.Achievements))
		}
	})

	t.Run("single", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddAchievement(conn, "some achievement")
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM achievements WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 1 {
			t.Fatalf("expected %d, received %d", 1, count)
		}

		if len(r.Achievements) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(r.Achievements))
		}
	})

	t.Run("many", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddAchievement(conn, "achievement0", "achievement1", "achievement2")
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM achievements WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 3 {
			t.Fatalf("expected %d, received %d", 3, count)
		}

		if len(r.Achievements) != 3 {
			t.Fatalf("expected %d, received %d", 3, len(r.Achievements))
		}
	})

	t.Run("ordered", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		achievements := []string{"achievement0", "achievement1", "achievement2", "achievement3", "achievement4"}

		err = r.AddAchievement(conn, achievements[0], achievements[1], achievements[2])
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddAchievement(conn, achievements[3], achievements[4])
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM achievements WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 5 {
			t.Fatalf("expected %d, received %d", 5, count)
		}

		if len(r.Achievements) != 5 {
			t.Fatalf("expected %d, received %d", 5, len(r.Achievements))
		}

		for i, achievement := range achievements {
			expected := fmt.Sprintf("achievement%d", i)
			if achievement != expected {
				t.Fatalf("expected %s, received %s", expected, achievement)
			}
		}
	})
}

func TestAddEducation(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddEducation(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM educations WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 0 {
			t.Fatalf("expected %d, received %d", 0, count)
		}

		if len(r.Educations) != 0 {
			t.Fatalf("expected %d, received %d", 0, len(r.Educations))
		}
	})

	t.Run("single", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddEducation(conn, resume.Example().Educations[0])
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM educations WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 1 {
			t.Fatalf("expected %d, received %d", 1, count)
		}

		if len(r.Educations) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(r.Educations))
		}
	})

	t.Run("many", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddEducation(conn, resume.Example().Educations...)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM educations WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != len(resume.Example().Educations) {
			t.Fatalf("expected %d, received %d", len(resume.Example().Educations), count)
		}

		if len(r.Educations) != len(resume.Example().Educations) {
			t.Fatalf("expected %d, received %d", len(resume.Example().Educations), len(r.Educations))
		}
	})
}

func TestAddWorkExperience(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddWorkExperience(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM work_experiences WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 0 {
			t.Fatalf("expected %d, received %d", 0, count)
		}

		if len(r.WorkExperiences) != 0 {
			t.Fatalf("expected %d, received %d", 0, len(r.WorkExperiences))
		}
	})

	t.Run("single", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddWorkExperience(conn, resume.Example().WorkExperiences[0])
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM work_experiences WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 1 {
			t.Fatalf("expected %d, received %d", 1, count)
		}

		if len(r.WorkExperiences) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(r.WorkExperiences))
		}

		for i, experience := range r.WorkExperiences {
			for j, responsibibility := range experience.Responsibilities {
				if responsibibility != resume.Example().WorkExperiences[i].Responsibilities[j] {
					t.Fatalf("expected %s, received %s", resume.Example().WorkExperiences[i].Responsibilities[j], responsibibility)
				}
			}
		}
	})

	t.Run("many", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddWorkExperience(conn, resume.Example().WorkExperiences...)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM work_experiences WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != len(resume.Example().WorkExperiences) {
			t.Fatalf("expected %d, received %d", len(resume.Example().WorkExperiences), count)
		}

		if len(r.WorkExperiences) != len(resume.Example().WorkExperiences) {
			t.Fatalf("expected %d, received %d", len(resume.Example().WorkExperiences), len(r.WorkExperiences))
		}

		for i, experience := range r.WorkExperiences {
			for j, responsibibility := range experience.Responsibilities {
				if responsibibility != resume.Example().WorkExperiences[i].Responsibilities[j] {
					t.Fatalf("expected %s, received %s", resume.Example().WorkExperiences[i].Responsibilities[j], responsibibility)
				}
			}
		}
	})
}

func TestAddProject(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddProject(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM projects WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 0 {
			t.Fatalf("expected %d, received %d", 0, count)
		}

		if len(r.Projects) != 0 {
			t.Fatalf("expected %d, received %d", 0, len(r.Projects))
		}
	})

	t.Run("single", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddProject(conn, resume.Example().Projects[0])
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM projects WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != 1 {
			t.Fatalf("expected %d, received %d", 1, count)
		}

		if len(r.Projects) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(r.Projects))
		}

		for i, project := range r.Projects {
			for j, responsibibility := range project.Responsibilities {
				if responsibibility != resume.Example().WorkExperiences[i].Responsibilities[j] {
					t.Fatalf("expected %s, received %s", resume.Example().WorkExperiences[i].Responsibilities[j], responsibibility)
				}
			}
		}
	})

	t.Run("many", func(t *testing.T) {
		conn := test.SetupDB(t)
		defer test.TearDownDB(t, conn)

		err := database.ApplyMigrations(conn, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		u, err := database.CreateUser(conn)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := u.AddResume(conn, resume.MinExample())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = r.AddProject(conn, resume.Example().Projects...)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		row := conn.Database().QueryRow("SELECT COUNT(1) FROM projects WHERE resume_id = ?", r.Id)

		var count int
		row.Scan(&count)

		if count != len(resume.Example().Projects) {
			t.Fatalf("expected %d, received %d", len(resume.Example().Projects), count)
		}

		if len(r.Projects) != len(resume.Example().Projects) {
			t.Fatalf("expected %d, received %d", len(resume.Example().Projects), len(r.Projects))
		}

		for i, project := range r.Projects {
			for j, responsibibility := range project.Responsibilities {
				if responsibibility != resume.Example().WorkExperiences[i].Responsibilities[j] {
					t.Fatalf("expected %s, received %s", resume.Example().WorkExperiences[i].Responsibilities[j], responsibibility)
				}
			}
		}
	})
}
