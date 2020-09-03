package main

import (
        "context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
 	"go.mongodb.org/mongo-driver/bson/primitive"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gorilla/mux"
)


// Client mongo Db
var collection *mongo.Collection
var ctx = context.TODO()

func init() {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    collection = client.Database("posts").Collection("post")


}

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

type Post struct {
        Titulo    string             `bson:"titulo"`
        SubTitulo string             `bson:"sub_titulo"`
        Conteudo  string             `bson:"conteudo"`
        Fotos     string             `bson:"fotos"`
        Autor     string             `bson:"autor"`
        Data      time.Time 	     `bson:"data"`
        ID primitive.ObjectID `bson:"permalink"`
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
	r.HandleFunc("/post", GetPost).Methods("GET")

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

	fmt.Print("blog")
	posts,_ := getAll()
	for _, v := range posts {
          fmt.Print(v.Titulo)
    	}
	data := GetPosts()
	blog := template.Must(template.ParseFiles("template/blog.html"))
	blog.Execute(w, data)
}

// Post page
func GetPost(w http.ResponseWriter, r *http.Request) {
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

var pp = &Post{
                Titulo:    "LOREM IPSUM  LOREM IPSUM LOREM IPSUM LOREM IPSUM ",
                SubTitulo: "LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM ",
                Conteudo:  "LOREM IPSUM LOREM IPSUM LOREIPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSU ",
                Fotos:     "static/img/portfolio/fullsize/1.jpg",
                Autor:     "Comuna da Catarina",
		ID: primitive.NewObjectID(),
		Data: time.Now()}

	err := createPost(pp)
	if err != nil {
          log.Fatal(err)
	}

	post := template.Must(template.ParseFiles("template/post.html"))
	post.Execute(w, p)
}

func createPost(post *Post) error {
    _, err := collection.InsertOne(ctx, post)
  return err
}


func getAll() ([]*Post, error) {
  // passing bson.D{{}} matches all documents in the collection
    filter := bson.D{{}}
    return filterPosts(filter)
}

func filterPosts(filter interface{}) ([]*Post, error) {
    // A slice of tasks for storing the decoded documents
    var posts []*Post

    cur, err := collection.Find(ctx, filter)
    if err != nil {
        return posts, err
    }

    for cur.Next(ctx) {
        var t Post
        err := cur.Decode(&t)
        if err != nil {
            return posts, err
        }

        posts = append(posts, &t)
    }

    if err := cur.Err(); err != nil {
        return posts, err
    }

  // once exhausted, close the cursor
    cur.Close(ctx)

    if len(posts) == 0 {
        return posts, mongo.ErrNoDocuments
    }

    return posts, nil
}
