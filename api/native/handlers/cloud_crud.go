package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mr-tasker/internal/services/user"
	"mr-tasker/internal/services/user/model"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type CloudCrudHandler struct {
	cloud user.UserService
}

func NewCloudCrudHandler(cloud user.UserService) *CloudCrudHandler {
	return &CloudCrudHandler{
		cloud: cloud,
	}
}

func (c *CloudCrudHandler) GetHandlers() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Route("/cloudapi", func(r chi.Router) {
		r.Use(TimerMiddleware)
		r.Get("/{id}", c.Get)
		r.Post("/", c.Create)
		r.Put("/", c.Update)
		r.Delete("/{id}", c.Remove)
	})
	return r
}

func (c *CloudCrudHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "panic" {
		panic("panic id")
	}

	// TODO: check if id uuid
	user, err := c.cloud.Read(id)
	if err != nil {
		log.Printf("[ERROR] failed to read user, err:%s", err)
		http.Error(w, "failed to read user", http.StatusInternalServerError)
		return
	}
	renderJson(w, user)
}

func (c *CloudCrudHandler) Create(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, "failed to parse media type", http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	byteBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	var User model.User
	err = json.Unmarshal(byteBody, &User)
	if err != nil {
		http.Error(w, "failed to unmarshal", http.StatusInternalServerError)
		return
	}
	id, err := c.cloud.Create(&User)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}
	renderPlainText(w, fmt.Sprintf("user with id = %s", id))
}

func (c *CloudCrudHandler) Update(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, "failed to parse media type", http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	byteBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	var User model.User
	err = json.Unmarshal(byteBody, &User)
	if err != nil {
		log.Printf("[ERROR] failed to unmarshal err:%s", err)
		http.Error(w, "failed to unmarshal", http.StatusInternalServerError)
		return
	}
	u, err := c.cloud.Update(&User)
	if err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}
	renderJson(w, u)
}

func (c *CloudCrudHandler) Remove(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// TODO: check if id uuid
	err := c.cloud.Delete(id)
	if err != nil {
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}
	renderPlainText(w, fmt.Sprintf("user is deleted, id = %s", id))
}

func renderJson(w http.ResponseWriter, resp any) {
	json, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to marshal user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func renderPlainText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}
