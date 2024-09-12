package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"resumegenerator/internal/auth"
	"resumegenerator/internal/database"
	"resumegenerator/pkg/resume"
)

func HandleCreateEducation(w http.ResponseWriter, r *http.Request, a *auth.Auth, db database.VersionedDatabase) {
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

	e, err := resume.EducationFromJSON(body)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if e.DegreeType == "" || e.FieldOfStudy == "" || e.Institution == "" {
		http.Error(w, "bad request", 400)
		return
	}

	err = database.CreateEducation(db, res, &e)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot create education")
		w.Write(internalError)
		w.WriteHeader(500)
		w.Header().Set("content-type", "application/json")
		return
	}

	response, err := json.Marshal(e)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot send education")
		w.Write(internalError)
		w.WriteHeader(500)
		w.Header().Set("content-type", "application/json")
		return
	}

	w.Write(response)
	w.WriteHeader(201)
	w.Header().Set("content-type", "application/json")
}

func HandleGetEducations(w http.ResponseWriter, r *http.Request, a *auth.Auth, db database.VersionedDatabase) {
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

	educations, err := database.ResumeEducations(db, res)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot get educations")
		w.Write(internalError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	response, err := json.Marshal(educations)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot send educations")
		w.Write(internalError)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(500)
		return
	}

	w.Write(response)
	w.Header().Set("content-type", "application/json")
}
