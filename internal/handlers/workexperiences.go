package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"resumegenerator/internal/auth"
	"resumegenerator/internal/database"
	"time"
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

	resume := database.GetResume(db, resumeId)
	if resume == nil || resume.Id == "" {
		http.Error(w, "not found", 404)
		return
	}

	if userId != resume.UserId {
		http.Error(w, "unauthorized", 401)
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

	var we struct {
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
	if err = json.Unmarshal(body, &we); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if we.Employer == "" || we.Title == "" || we.Began == "" {
		http.Error(w, "bad request", 400)
		return
	}

	began, err := time.Parse(time.RFC3339, we.Began)
	if err = json.Unmarshal(body, &we); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	var finished *time.Time
	if we.Finished != nil {
		val, err := time.Parse(time.RFC3339, *we.Finished)
		if err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		finished = &val
	}

	workExperience, err := database.CreateWorkExperience(
		db,
		resume,
		we.Employer,
		we.Title,
		began,
		we.Current,
		we.Location,
		finished,
	)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot create work experience")
		w.Write(internalError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	for _, responsibility := range we.Responsibilities {
		_, err = database.CreateWorkResponsibility(db, &workExperience, responsibility.Responsibility)
		if err != nil {
			internalError, _ := newInternalError(err, "cannot create work responsibility")
			w.Write(internalError)
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(500)
			return
		}
	}

	response, err := json.Marshal(workExperience)
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

	resume := database.GetResume(db, resumeId)
	if resume == nil || resume.Id == "" {
		http.Error(w, "not found", 404)
		return
	}

	if userId != resume.UserId {
		http.Error(w, "unauthorized", 401)
		return
	}

	workExperiences, err := resume.WorkExperiences(db)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot get work experiences")
		w.Write(internalError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	response, err := json.Marshal(workExperiences)
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
