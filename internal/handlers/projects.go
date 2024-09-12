package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"resumegenerator/internal/auth"
	"resumegenerator/internal/database"
	"resumegenerator/pkg/resume"
)

func HandleCreateProject(w http.ResponseWriter, r *http.Request, a *auth.Auth, db database.VersionedDatabase) {
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

	p, err := resume.NewProject(body)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if p.Name == "" || p.Role == "" || p.Description == "" {
		http.Error(w, "bad request", 400)
		return
	}

	err = database.CreateProject(db, res, &p)
	if err != nil {
		internalErr, _ := newInternalError(err, "cannot create project")
		w.Write(internalErr)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	response, err := json.Marshal(&p)
	if err != nil {
		internalErr, _ := newInternalError(err, "cannot send project")
		w.Write(internalErr)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	w.Write(response)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(201)
}

func HandleGetProjects(w http.ResponseWriter, r *http.Request, a *auth.Auth, db database.VersionedDatabase) {
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

	projects, err := database.ResumeProjects(db, res)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot get projects")
		w.Write(internalError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	response, err := json.Marshal(projects)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot send projects")
		w.Write(internalError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	w.Write(response)
	w.Header().Set("content-type", "application/json")
}
