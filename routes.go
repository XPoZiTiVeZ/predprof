package main

import (
	"net/http"
)

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/",         HomePageHandler)
	mux.HandleFunc("/login",    LoginPageHandler)
	mux.HandleFunc("/register", RegisterPageHandler)
	mux.HandleFunc("/profile",  ProfilePageHandler)
	mux.HandleFunc("/logout",   LogoutPageHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
}
