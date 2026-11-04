package server

import (
	"net/http"

	"github.com/CALEXCO/minisocial/internal/user"
)

func routes(users *user.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("POST /users", users.Create)

	return mux
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
