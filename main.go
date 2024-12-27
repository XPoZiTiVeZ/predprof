package main

import (
	"html/template"
	_ "html/template"
	"net/http"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./templates/index.jinja")
}

func main() {
	http.HandleFunc("/", IndexPage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8080", nil)
}
