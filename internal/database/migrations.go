package database

func UpMigrations() []string {
	return []string{
		`
CREATE TABLE users (
  user_id VARCHAR(36) PRIMARY KEY,
  created_at INTEGER NOT NULL --> TIMESTAMP
);

CREATE TABLE resumes (
  resume_id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  created_at INTEGER NOT NULL, --> TIMESTAMP

  name VARCHAR(255) NOT NULL,
  title VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  phone_number VARCHAR(255) NOT NULL,
  summary TEXT NOT NULL,
  location VARCHAR(255),
  linkedin_username VARCHAR(255),
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

CREATE TABLE skills (
  skill_id VARCHAR(36) PRIMARY KEY,
  resume_id VARCHAR(36) NOT NULL,
  position INTEGER NOT NULL DEFAULT 0,
  skill TEXT NOT NULL,

  CONSTRAINT fk_resume
    FOREIGN KEY(resume_id)
      REFERENCES resumes(resume_id)
      ON DELETE CASCADE
);

CREATE TABLE achievements (
  achievement_id VARCHAR(36) PRIMARY KEY,
  resume_id VARCHAR(36) NOT NULL,
  position INTEGER NOT NULL DEFAULT 0,
  achievement TEXT NOT NULL,

  CONSTRAINT fk_resume
    FOREIGN KEY(resume_id)
      REFERENCES resumes(resume_id)
      ON DELETE CASCADE
);

CREATE TABLE educations (
  education_id VARCHAR(36) PRIMARY KEY,
  resume_id VARCHAR(36) NOT NULL,
  created_at INTEGER NOT NULL, --> TIMESTAMP

  degree VARCHAR(255) NOT NULL,
  field VARCHAR(255) NOT NULL,
  institution VARCHAR(255) NOT NULL,
  began INTEGER NOT NULL, --> TIMESTAMP
  current INTEGER DEFAULT 0 NOT NULL,
  location VARCHAR(255),
  finished INTEGER, --> TIMESTAMP
  gpa VARCHAR(255),

  CONSTRAINT fk_resume
    FOREIGN KEY(resume_id)
      REFERENCES resumes(resume_id)
      ON DELETE CASCADE
);

CREATE TABLE work_experiences (
  work_experience_id VARCHAR(36) PRIMARY KEY,
  resume_id VARCHAR(36) NOT NULL,
  created_at INTEGER NOT NULL, --> TIMESTAMP

  employer VARCHAR(255) NOT NULL,
  title VARCHAR(255) NOT NULL,
  began INTEGER NOT NULL, --> TIMESTAMP
  current INTEGER DEFAULT 0 NOT NULL,
  location VARCHAR(255),
  finished INTEGER, --> TIMESTAMP

  CONSTRAINT fk_resume
    FOREIGN KEY(resume_id)
      REFERENCES resumes(resume_id)
      ON DELETE CASCADE
);

CREATE TABLE work_responsibilities (
  work_responsibility_id VARCHAR(36) PRIMARY KEY,
  work_experience_id VARCHAR(36) NOT NULL,

  responsibility VARCHAR(255) NOT NULL,
  position INTEGER NOT NULL DEFAULT 0,

  CONSTRAINT fk_work_experience
    FOREIGN KEY(work_experience_id)
      REFERENCES work_experiences(work_experience_id)
      ON DELETE CASCADE
);

CREATE TABLE projects (
  project_id VARCHAR(36) PRIMARY KEY,
  resume_id VARCHAR(36) NOT NULL,
  created_at INTEGER NOT NULL, --> TIMESTAMP

  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  role VARCHAR(255) NOT NULL,

  CONSTRAINT fk_resume
    FOREIGN KEY(resume_id)
      REFERENCES resumes(resume_id)
      ON DELETE CASCADE
);

CREATE TABLE project_responsibilities (
  project_responsibility_id VARCHAR(36) PRIMARY KEY,
  project_id VARCHAR(36) NOT NULL,

  responsibility VARCHAR(255) NOT NULL,
  position INTEGER NOT NULL DEFAULT 0,

  CONSTRAINT fk_project
    FOREIGN KEY(project_id)
      REFERENCES projects(project_id)
      ON DELETE CASCADE
);
    `,
	}
}

func DownMigrations() []string {
	return []string{
		`
DROP TABLE project_responsibilities;
DROP TABLE projects;
DROP TABLE work_responsibilities;
DROP TABLE work_experiences;
DROP TABLE educations;
DROP TABLE achievements;
DROP TABLE skills;
DROP TABLE resumes;
DROP TABLE users;
    `,
	}
}
