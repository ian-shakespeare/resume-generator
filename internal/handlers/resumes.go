package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"resumegenerator/internal/auth"
	"resumegenerator/internal/database"
	"resumegenerator/pkg/resume"
	"sync"
)

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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	newResume, err := resume.New(body)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if newResume.Name == "" || newResume.Email == "" || newResume.PhoneNumber == "" || newResume.Prelude == "" {
		http.Error(w, "bad request", 400)
		return
	}

	if err = database.CreateResume(db, user, &newResume); err != nil {
		internalErr, _ := newInternalError(err, "cannot create resume")
		w.Write(internalErr)
		w.WriteHeader(500)
		w.Header().Add("content-type", "application/json")
		return
	}

	response, err := json.Marshal(newResume)
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
	Resume         resume.Resume           `json:"resume"`
	Education      []resume.Education      `json:"education"`
	WorkExperience []resume.WorkExperience `json:"workExperiences"`
	Projects       []resume.Project        `json:"projects"`
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

	res := database.GetResume(db, resumeId, userId)
	if res == nil || res.Id == "" {
		http.Error(w, "not found", 404)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	var educations []resume.Education
	var workExperiences []resume.WorkExperience
	var projects []resume.Project
	var dbErr error

	go func() {
		defer wg.Done()
		educations, err = database.ResumeEducations(db, res)
		if err != nil && dbErr == nil {
			dbErr = err
		}
	}()

	go func() {
		defer wg.Done()
		workExperiences, err = database.ResumeWorkExperiences(db, res)
		if err != nil && dbErr == nil {
			dbErr = err
		}
	}()

	go func() {
		defer wg.Done()
		projects, err = database.ResumeProjects(db, res)
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
		Resume:         *res,
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
