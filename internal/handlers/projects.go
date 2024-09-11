package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"resumegenerator/internal/auth"
	"resumegenerator/internal/database"
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

	var p struct {
		Name             string `json:"name"`
		Description      string `json:"description"`
		Role             string `json:"role"`
		Responsibilities []struct {
			Responsibility string `json:"responsibility"`
		} `json:"responsibilities"`
	}
	if err = json.Unmarshal(body, &p); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if p.Name == "" || p.Role == "" || p.Description == "" {
		http.Error(w, "bad request", 400)
		return
	}

	project, err := database.CreateProject(db, resume, p.Name, p.Description, p.Role)
	if err != nil {
		internalErr, _ := newInternalError(err, "cannot create project")
		w.Write(internalErr)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	for _, responsibility := range p.Responsibilities {
		_, err = database.CreateProjectResponsibility(db, &project, responsibility.Responsibility)
		if err != nil {
			internalErr, _ := newInternalError(err, "cannot create project responsibility")
			w.Write(internalErr)
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(500)
			return
		}
	}

	response, err := json.Marshal(project)
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

	resume := database.GetResume(db, resumeId)
	if resume == nil || resume.Id == "" {
		http.Error(w, "not found", 404)
		return
	}

	if userId != resume.UserId {
		http.Error(w, "unauthorized", 401)
		return
	}

	projects, err := resume.Projects(db)
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
