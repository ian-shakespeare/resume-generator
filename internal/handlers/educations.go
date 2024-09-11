package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"resumegenerator/internal/auth"
	"resumegenerator/internal/database"
	"time"
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

	var ne struct {
		DegreeType   string  `json:"degreeType"`
		FieldOfStudy string  `json:"fieldOfStudy"`
		Institution  string  `json:"institution"`
		Began        string  `json:"began"`
		Current      bool    `json:"current"`
		Location     *string `json:"location"`
		Finished     *string `json:"finished"`
		GPA          *string `json:"gpa"`
	}
	if err = json.Unmarshal(body, &ne); err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	if ne.DegreeType == "" || ne.FieldOfStudy == "" || ne.Institution == "" || ne.Began == "" {
		http.Error(w, "bad request", 400)
		return
	}

	began, err := time.Parse(time.RFC3339, ne.Began)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	var finished *time.Time
	if ne.Finished != nil {
		val, err := time.Parse(time.RFC3339, *ne.Finished)
		if err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		finished = &val
	}

	education, err := database.CreateEducation(
		db,
		resume,
		ne.DegreeType,
		ne.FieldOfStudy,
		ne.Institution,
		began,
		ne.Current,
		ne.Location,
		finished,
		ne.GPA,
	)
	if err != nil {
		internalError, _ := newInternalError(err, "cannot create education")
		w.Write(internalError)
		w.WriteHeader(500)
		w.Header().Set("content-type", "application/json")
		return
	}

	response, err := json.Marshal(education)
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

	resume := database.GetResume(db, resumeId)
	if resume == nil || resume.Id == "" {
		http.Error(w, "not found", 404)
		return
	}

	if userId != resume.UserId {
		http.Error(w, "unauthorized", 401)
		return
	}

	educations, err := resume.Educations(db)
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
