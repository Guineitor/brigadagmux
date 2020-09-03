package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/manifesto", Manifesto).Methods("GET")
	r.HandleFunc("/blog/{page}", Blog).Methods("GET")
	r.HandleFunc("/post/{permalink}", Post).Methods("GET")

	http.ListenAndServe(":9990", r)
}

// Manifesto page
func Manifesto(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("template/manifesto.html")).Execute(w, struct{ Success bool }{true})
}

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("template/index.html")).Execute(w, struct{ Success bool }{true})

}

// Blog page
func Blog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	fmt.Fprintf(w, "Blog: %s", page)
}

// Post page
func Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	permalink := vars["permalink"]
	fmt.Fprintf(w, "Post: %s", permalink)
}
