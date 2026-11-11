package server

import (
	"net/http"

	"github.com/CALEXCO/minisocial/internal/user"
)

func routes(users *user.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("POST /users", users.Create)
	mux.HandleFunc("GET /users", users.GetAll)
	mux.HandleFunc("GET /users/{id}", users.GetById)
	mux.HandleFunc("PUT /users/{id}", users.GetById)

	return mux
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
