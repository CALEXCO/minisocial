package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

func Init() {
	svr := http.DefaultServeMux
	svr.HandleFunc("/hello", hello)

	if err := http.ListenAndServe(":8080", svr); err != nil {
		panic(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	log := slog.Default()
	log.Debug("hello")

	fmt.Fprint(w, "hello world")
}
