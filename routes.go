package main

import (
	"net/http"
)

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/",         HomePage)
	mux.HandleFunc("/login",    LoginPage)
	mux.HandleFunc("/register", RegisterPage)
	mux.HandleFunc("/profile",  ProfilePage)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
}
