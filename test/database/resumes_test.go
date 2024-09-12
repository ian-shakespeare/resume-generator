package database_test

import (
	"resumegenerator/internal/database"
	"resumegenerator/pkg/resume"
	"resumegenerator/test"
	"testing"
)

func TestCreateResume(t *testing.T) {
	t.Run("minimumRequiredFields", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations(), 1)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(
			db,
			&user,
			&r,
		)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})

	t.Run("allFields", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.New([]byte(test.FULL_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(
			db,
			&user,
			&r,
		)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
	})
}

func TestGetResume(t *testing.T) {
	t.Run("minimumRequiredFields", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(
			db,
			&user,
			&r,
		)

		stored := database.GetResume(db, r.Id, user.Id)

		if stored == nil {
			t.Fatalf("expected %s, received %s", "resume", "nil")
		}
		if stored.Location != "" {
			t.Fatalf("expected %s, received %s", "nil", stored.Location)
		}
	})

	t.Run("allFields", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.New([]byte(test.FULL_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(
			db,
			&user,
			&r,
		)

		stored := database.GetResume(db, r.Id, user.Id)

		if stored == nil {
			t.Fatalf("expected %s, received %s", "resume", "nil")
		}
		if stored.Location != r.Location {
			t.Fatalf("expected %s, received %s", r.Location, stored.Location)
		}
	})
}

func TestResumeGetEducations(t *testing.T) {
	t.Run("nonExistant", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		educations, err := database.ResumeEducations(db, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}
		if len(educations) != 0 {
			t.Fatalf("expected %d, received %d", 0, len(educations))
		}
	})

	t.Run("existant", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		e, err := resume.NewEducation([]byte(test.MIN_EDUCATION))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateEducation(db, &r, &e)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		educations, err := database.ResumeEducations(db, &r)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(educations) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(educations))
		}

		if educations[0].Id != e.Id {
			t.Fatalf("expected %s, received %s", e.Id, educations[0].Id)
		}
	})
}

func TestResumeGetWorkExperiences(t *testing.T) {
	t.Run("nonExistant", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		workExperiences, err := database.ResumeWorkExperiences(db, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(workExperiences) != 0 {
			t.Fatalf("expected %d, received %d", 1, len(workExperiences))
		}
	})

	t.Run("existant", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		w, err := resume.NewWorkExperience([]byte(test.MIN_WORK_EXPERIENCE))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateWorkExperience(db, &r, &w)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		workExperiences, err := database.ResumeWorkExperiences(db, &r)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(workExperiences) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(workExperiences))
		}

		if workExperiences[0].Id != w.Id {
			t.Fatalf("expected %s, received %s", w.Id, workExperiences[0].Id)
		}
	})
}

func TestResumeGetProjects(t *testing.T) {
	t.Run("nonExistant", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		projects, err := database.ResumeProjects(db, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(projects) != 0 {
			t.Fatalf("expected %d, received %d", 1, len(projects))
		}
	})

	t.Run("existant", func(t *testing.T) {
		db := test.SetupDB(t)
		defer test.TearDownDB(t, db)

		err := database.ApplyMigrations(db, database.UpMigrations(), database.DownMigrations())
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		user, err := database.CreateUser(db)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		r, err := resume.New([]byte(test.MIN_RESUME))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateResume(db, &user, &r)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		p, err := resume.NewProject([]byte(test.PROJECT))
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		err = database.CreateProject(db, &r, &p)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		projects, err := database.ResumeProjects(db, &r)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if len(projects) != 1 {
			t.Fatalf("expected %d, received %d", 1, len(projects))
		}

		if projects[0].Id != p.Id {
			t.Fatalf("expected %s, received %s", p.Id, projects[0].Id)
		}
	})
}
