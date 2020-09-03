package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// P
type P struct {
	Titulo    string `json:"titulo"`
	SubTitulo string `json:"sub_titulo"`
	Conteudo  string `json:"conteudo"`
	Fotos     string `json:"fotos"`
	Autor     string `json:"autor"`
	Data      string `json:"data"`
	Permalink string `json:"permalink"`
}

// Posts
type Posts struct {
	Posts []P
}

// Find posts
func FindPost() []P {

	jsonFile, err := os.Open("posts.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened posts.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var p []P

	json.Unmarshal(byteValue, &p)

	return p
}

// Get posts
func GetPosts() Posts {
	data := Posts{
		Posts: FindPost(),
	}
	return data
}

func main() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("assets/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/manifesto", Manifesto).Methods("GET")
	r.HandleFunc("/blog", Blog).Methods("GET")
	r.HandleFunc("/post", Post).Methods("GET")

	http.ListenAndServe(":9991", r)
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
	data := GetPosts()
	blog := template.Must(template.ParseFiles("template/blog.html"))
	blog.Execute(w, data)
}

// Post page
func Post(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// permalink := vars["permalink"]
	var p = P{
		Titulo:    "LOREM IPSUM  LOREM IPSUM LOREM IPSUM LOREM IPSUM ",
		SubTitulo: "LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM ",
		Conteudo:  "LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM ",
		Fotos:     "static/img/portfolio/fullsize/1.jpg",
		Autor:     "Comuna da Catarina",
		Data:      "30/06/2020",
		Permalink: "poder"}

	post := template.Must(template.ParseFiles("template/post.html"))
	post.Execute(w, p)
}
