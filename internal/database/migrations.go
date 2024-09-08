package database

func ForwardMigrations() []string {
	return []string{
		`
CREATE TABLE IF NOT EXISTS users (
  user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4()
);

CREATE TABLE IF NOT EXISTS resumes (
  resume_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  phone_number VARCHAR(255) NOT NULL,

  location VARCHAR(255),
  linked_in_username VARCHAR(255),
  github_username VARCHAR(255),
  facebook_username VARCHAR(255),
  instagram_username VARCHAR(255),
  twitter_handle VARCHAR(255),
  portfolio VARCHAR(255),

  CONSTRAINT fk_user
    FOREIGN KEY(user_id)
      REFERENCES users(user_id)
      ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS educations (
  education_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  resume_id UUID NOT NULL,
  degree_type VARCHAR(255) NOT NULL,
  field_of_study VARCHAR(255) NOT NULL,
  institution VARCHAR(255) NOT NULL,
  began DATE NOT NULL,
  current BOOLEAN DEFAULT false NOT NULL,

  location VARCHAR(255),
  finished DATE,
  gpa VARCHAR(255),

  CONSTRAINT fk_resume
    FOREIGN KEY(resume_id)
      REFERENCES resumes(resume_id)
      ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS work_experiences (
  work_experience_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  resume_id UUID NOT NULL,
  employer VARCHAR(255) NOT NULL,
  title VARCHAR(255) NOT NULL,
  began DATE NOT NULL,
  current BOOLEAN DEFAULT false NOT NULL,

  location VARCHAR(255),
  finished DATE,

  CONSTRAINT fk_resume
    FOREIGN KEY(resume_id)
      REFERENCES resumes(resume_id)
      ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS work_responsibilities (
  work_responsibility_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  work_experience_id UUID NOT NULL,
  responsibility VARCHAR(255) NOT NULL,

  CONSTRAINT fk_work_experience
    FOREIGN KEY(work_experience_id)
      REFERENCES work_experiences(work_experience_id)
      ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS projects (
  project_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  resume_id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  role VARCHAR(255) NOT NULL,

  CONSTRAINT fk_resume
    FOREIGN KEY(resume_id)
      REFERENCES resumes(resume_id)
      ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_responsibilities (
  project_responsibility_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  project_id UUID NOT NULL,
  responsibility VARCHAR(255) NOT NULL,

  CONSTRAINT fk_project
    FOREIGN KEY(project_id)
      REFERENCES projects(project_id)
      ON DELETE CASCADE
);
    `,
	}
}

func BackwardMigrations() []string {
	return []string{
		`
DROP TABLE IF EXISTS project_responsibilities;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS work_responsibilities;
DROP TABLE IF EXISTS work_experiences;
DROP TABLE IF EXISTS educations;
DROP TABLE IF EXISTS resumes;
DROP TABLE IF EXISTS users;
    `,
	}
}
