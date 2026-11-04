package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewUserHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Mail     string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := h.service.Create(r.Context(), req.Username, req.Mail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(u); err != nil {
		slog.Default().Error("json encode failed", err)
	}
}
func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) UpdateName(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

}
