package user

import (
	"net/http"
)

type Handler struct {
	service *Service
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) UpdateName(w http.ResponseWriter, r *http.Request) {

}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

}
