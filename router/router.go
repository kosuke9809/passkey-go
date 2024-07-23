package router

import (
	"net/http"
	"passkey-auth/auth"
)

type Router struct {
	http.ServeMux
}

func New() (*Router, error) {
	err := auth.InitWebAuthn()
	if err != nil {
		return nil, err
	}

	s := &Router{}
	s.HandleFunc("/register/begin", auth.BeginRegistration)
	s.HandleFunc("/register/finish", auth.FinishRegistration)
	s.HandleFunc("/login/begin", auth.BeginLogin)
	s.HandleFunc("/login/finish", auth.FinishLogin)
	s.HandleFunc("/", serveHTML)
	return s, nil
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func (s *Router) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}
