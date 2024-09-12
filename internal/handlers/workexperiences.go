package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"resumegenerator/internal/auth"
	"resumegenerator/internal/database"
	"resumegenerator/pkg/resume"
)

type NewWorkResponsibility struct {
	Responsibility string `json:"responsibility"`
}

type NewWorkExperience struct {
	Employer         string `json:"employer"`
	Title            string `json:"title"`
	Began            string `json:"began"`
	Current          bool   `json:"current"`
	Responsibilities []struct {
		Responsibility string `json:"responsibility"`
	} `json:"responsibilities"`
	Location *string `json:"location"`
	Finished *string `json:"finished"`
}

func HandleCreateWorkExperience(w http.ResponseWriter, r *http.Request, a *auth.Auth, db database.VersionedDatabase) {
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

	if r.Body == nil {
		http.Error(w, "bad request", 400)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	we, err := resume.NewWorkExperience(body)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if we.Employer == "" || we.Title == "" {
		http.Error(w, "bad request", 400)
		return
	}

	err = database.CreateWorkExperience(db, res, &we)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot create work experience")
		w.Write(internalError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	response, err := json.Marshal(we)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot send work experience")
		w.Write(internalError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	w.Write(response)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(201)
}

func HandleGetWorkExperiences(w http.ResponseWriter, r *http.Request, a *auth.Auth, db database.VersionedDatabase) {
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

	we, err := database.ResumeWorkExperiences(db, res)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot get work experiences")
		w.Write(internalError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	response, err := json.Marshal(we)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot send work experiences")
		w.Write(internalError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	w.Write(response)
	w.Header().Set("content-type", "application/json")
}
