package server

import (
	"net/http"

	"github.com/CALEXCO/minisocial/internal/config"
	"github.com/CALEXCO/minisocial/internal/user"
)

func Init(conf config.Config, users *user.Handler) error {
	srv := http.Server{
		Addr:    conf.HTTPAddr,
		Handler: routes(users),
	}

	return srv.ListenAndServe()
}
