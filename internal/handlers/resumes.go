package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"resumegenerator/internal/auth"
	"resumegenerator/internal/database"
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
		// TODO: Should maybe be a different error.
		http.Error(w, "bad request", 400)
		return
	}

	response, err := json.Marshal(resume)
	if err != nil {
		// TODO: Should maybe be a different error.
		http.Error(w, "bad request", 400)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		// TODO: Should maybe be a different error.
		http.Error(w, "bad request", 400)
		return
	}
	w.WriteHeader(201)
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
}
