package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"resumegenerator/internal/auth"
	"resumegenerator/internal/database"
	"sync"
)

type NewResume struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	Prelude     string  `json:"prelude"`
	Location    *string `json:"location"`
	LinkedIn    *string `json:"linkedIn"`
	Github      *string `json:"github"`
	Facebook    *string `json:"facebook"`
	Instagram   *string `json:"instagram"`
	Twitter     *string `json:"twitter"`
	Portfolio   *string `json:"portfolio"`
}

func HandleCreateResume(w http.ResponseWriter, r *http.Request, a *auth.Auth, db database.VersionedDatabase) {
	token, err := a.DecodeAuthHeader(r.Header)
	if err != nil {
		http.Error(w, "unauthorized", 401)
		return
	}

	userId, err := a.TokenUserId(token)
	if err != nil {
		http.Error(w, "unauthorized", 401)
		return
	}

	user := database.GetUser(db, userId)
	if user == nil {
		http.Error(w, "unauthorized", 401)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	var newResume NewResume
	if err = json.Unmarshal(reqBody, &newResume); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if newResume.Name == "" || newResume.Email == "" || newResume.PhoneNumber == "" || newResume.Prelude == "" {
		http.Error(w, "bad request", 400)
		return
	}

	resume, err := database.CreateResume(
		db,
		user,
		newResume.Name,
		newResume.Email,
		newResume.PhoneNumber,
		newResume.Prelude,
		newResume.Location,
		newResume.LinkedIn,
		newResume.Github,
		newResume.Facebook,
		newResume.Instagram,
		newResume.Twitter,
		newResume.Portfolio,
	)
	if err != nil {
		internalErr, _ := newInternalError(err, "cannot create resume")
		w.Write(internalErr)
		w.WriteHeader(500)
		w.Header().Add("content-type", "application/json")
		return
	}

	response, err := json.Marshal(resume)
	if err != nil {
		internalErr, _ := newInternalError(err, "cannot send resume")
		w.Write(internalErr)
		w.WriteHeader(500)
		w.Header().Add("content-type", "application/json")
		return
	}

	w.Write(response)
	w.WriteHeader(201)
	w.Header().Add("content-type", "application/json")
}

type FullResume struct {
	Resume         database.Resume           `json:"resume"`
	Education      []database.Education      `json:"education"`
	WorkExperience []database.WorkExperience `json:"workExperiences"`
	Projects       []database.Project        `json:"projects"`
}

func HandleGetResume(w http.ResponseWriter, r *http.Request, a *auth.Auth, db database.VersionedDatabase) {
	token, err := a.DecodeAuthHeader(r.Header)
	if err != nil {
		http.Error(w, "unauthorized", 401)
		return
	}

	userId, err := a.TokenUserId(token)
	if err != nil {
		http.Error(w, "unauthorized", 401)
		return
	}

	user := database.GetUser(db, userId)
	if user == nil {
		http.Error(w, "unauthorized", 401)
		return
	}

	resumeId := r.PathValue("resumeId")
	if resumeId == "" {
		http.Error(w, "not found", 404)
		return
	}

	resume := database.GetResume(db, resumeId)
	if resume == nil || resume.Id == "" {
		http.Error(w, "not found", 404)
		return
	}

	if resume.UserId != user.Id {
		http.Error(w, "unauthorized", 401)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	var educations []database.Education
	var workExperiences []database.WorkExperience
	var projects []database.Project
	var dbErr error

	go func() {
		defer wg.Done()
		educations, err = resume.Educations(db)
		if err != nil && dbErr == nil {
			dbErr = err
		}
	}()

	go func() {
		defer wg.Done()
		workExperiences, err = resume.WorkExperiences(db)
		if err != nil && dbErr == nil {
			dbErr = err
		}
	}()

	go func() {
		defer wg.Done()
		projects, err = resume.Projects(db)
		if err != nil && dbErr == nil {
			dbErr = err
		}
	}()

	wg.Wait()

	if dbErr != nil {
		internalErr, _ := newInternalError(err, "cannot get full resume")
		w.Write(internalErr)
		w.WriteHeader(500)
		w.Header().Add("content-type", "application/json")
		return
	}

	fullResume := FullResume{
		Resume:         *resume,
		Education:      educations,
		WorkExperience: workExperiences,
		Projects:       projects,
	}

	response, err := json.Marshal(fullResume)
	if err != nil {
		internalErr, _ := newInternalError(err, "cannot send resume")
		w.Write(internalErr)
		w.WriteHeader(500)
		w.Header().Add("content-type", "application/json")
		return
	}

	w.Write(response)
	w.Header().Add("content-type", "application/json")
}
