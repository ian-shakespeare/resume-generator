package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/tests"
	"testing"
	"time"
)

func TestCreateResume(t *testing.T) {
	t.Run("minimumRequiredFields", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations(), 1)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		_, err = database.CreateResume(
			db,
			&user,
			"John Doe",
			"jdoe@email.com",
			"+1 (000) 000-0000",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})

	t.Run("allFields", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		location := "SOME_PLACE"
		linkedIn := "SOME_LINKEDIN"
		github := "SOME_GITHUB"
		facebook := "SOME_FACEBOOK"
		instagram := "SOME_INSTAGRAM"
		twitter := "SOME_TWITTER"
		portfolio := "SOME_PORTFOLIO"

		_, err = database.CreateResume(
			db,
			&user,
			"John Doe",
			"jdoe@email.com",
			"+1 (000) 000-0000",
			"prelude",
			&location,
			&linkedIn,
			&github,
			&facebook,
			&instagram,
			&twitter,
			&portfolio,
		)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})
}

func TestGetResume(t *testing.T) {
	t.Run("minimumRequiredFields", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		created, err := database.CreateResume(
			db,
			&user,
			"John Doe",
			"jdoe@email.com",
			"+1 (000) 000-0000",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		resume := database.GetResume(db, created.Id)

		if resume == nil {
			t.Fatalf("expected %s, received %s", "resume", "nil")
		}
		if resume.Location != nil {
			t.Fatalf("expected %s, received %p", "nil", resume.Location)
		}
	})

	t.Run("allFields", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		location := "SOME_PLACE"
		linkedIn := "SOME_LINKEDIN"
		github := "SOME_GITHUB"
		facebook := "SOME_FACEBOOK"
		instagram := "SOME_INSTAGRAM"
		twitter := "SOME_TWITTER"
		portfolio := "SOME_PORTFOLIO"

		created, err := database.CreateResume(
			db,
			&user,
			"John Doe",
			"jdoe@email.com",
			"+1 (000) 000-0000",
			"prelude",
			&location,
			&linkedIn,
			&github,
			&facebook,
			&instagram,
			&twitter,
			&portfolio,
		)

		resume := database.GetResume(db, created.Id)

		if resume == nil {
			t.Fatalf("expected %s, received %s", "resume", "nil")
		}
		if *resume.Location != location {
			t.Fatalf("expected %s, received %s", location, *resume.Location)
		}
	})
}

func TestResumeGetEducations(t *testing.T) {
	t.Run("nonExistant", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&user,
			"name",
			"email",
			"phoneNumber",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		educations, err := resume.Educations(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
		if len(educations) != 0 {
			t.Fatalf("expected %d, received %d", 0, len(educations))
		}
	})

	t.Run("existant", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&user,
			"name",
			"email",
			"phoneNumber",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		education, err := database.CreateEducation(
			db,
			&resume,
			"degree",
			"fieldOfStudy",
			"institution",
			time.Now(),
			true,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		educations, err := resume.Educations(db)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(educations) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(educations))
		}

		if educations[0].Id != education.Id {
			t.Fatalf("expected %s, received %s", education.Id, educations[0].Id)
		}
	})
}

func TestResumeGetWorkExperiences(t *testing.T) {
	t.Run("nonExistant", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&user,
			"name",
			"email",
			"phoneNumber",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		workExperiences, err := resume.WorkExperiences(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(workExperiences) != 0 {
			t.Fatalf("expected %d, received %d", 1, len(workExperiences))
		}
	})

	t.Run("existant", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&user,
			"name",
			"email",
			"phoneNumber",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		workExperience, err := database.CreateWorkExperience(
			db,
			&resume,
			"employer",
			"title",
			time.Now(),
			true,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		workExperiences, err := resume.WorkExperiences(db)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(workExperiences) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(workExperiences))
		}

		if workExperiences[0].Id != workExperience.Id {
			t.Fatalf("expected %s, received %s", workExperience.Id, workExperiences[0].Id)
		}
	})
}

func TestResumeGetProjects(t *testing.T) {
	t.Run("nonExistant", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&user,
			"name",
			"email",
			"phoneNumber",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		projects, err := resume.Projects(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(projects) != 0 {
			t.Fatalf("expected %d, received %d", 1, len(projects))
		}
	})

	t.Run("existant", func(t *testing.T) {
		db := tests.SetupDB(t)
		defer tests.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		resume, err := database.CreateResume(
			db,
			&user,
			"name",
			"email",
			"phoneNumber",
			"prelude",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		project, err := database.CreateProject(
			db,
			&resume,
			"name",
			"description",
			"role",
		)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		projects, err := resume.Projects(db)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(projects) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(projects))
		}

		if projects[0].Id != project.Id {
			t.Fatalf("expected %s, received %s", project.Id, projects[0].Id)
		}
	})
}
