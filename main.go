package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Templates
const (
	IndexTemplate     = "template/index.html"
	BlogTemplate      = "template/blog.html"
	NotFoundTemplate  = "template/404.html"
	PostTemplate      = "template/post.html"
	ManifestoTemplate = "template/manifesto.html"
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

type Post struct {
	Titulo    string             `bson:"titulo"`
	SubTitulo string             `bson:"sub_titulo"`
	Conteudo  string             `bson:"conteudo"`
	Fotos     string             `bson:"fotos"`
	Autor     string             `bson:"autor"`
	Data      time.Time          `bson:"data"`
	ID        primitive.ObjectID `bson:"permalink"`
}

type Posts struct {
	Posts []*Post
}

func main() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("assets/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/manifesto", Manifesto).Methods("GET")
	r.HandleFunc("/blog", Blog).Methods("GET")
	r.HandleFunc("/post/{permalink}", GetPost).Methods("GET")

	http.ListenAndServe(":9991", r)
}

// Manifesto page
func Manifesto(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles(ManifestoTemplate)).Execute(w, struct{ Success bool }{true})
}

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles(IndexTemplate)).Execute(w, struct{ Success bool }{true})
}

// Blog page
func Blog(w http.ResponseWriter, r *http.Request) {
	posts, _ := getAll()
	for _, v := range posts {
		fmt.Print(v.Titulo)
	}

	data := Posts{
		Posts: posts,
	}

	blog := template.Must(template.ParseFiles(BlogTemplate))
	blog.Execute(w, data)
}

// Post page
func GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	permalink := vars["permalink"]

	data, err := getByID(permalink)
	// var p = &Post{
	// 	Titulo:    "LOREM IPSUM  LOREM IPSUM LOREM IPSUM LOREM IPSUM ",
	// 	SubTitulo: "LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM ",
	// 	Conteudo:  "LOREM IPSUM LOREM IPSUM LOREIPSUM LOREM IPSUM LOREM IPSUM LOREM IPSUM LOREM IPSU ",
	// 	Fotos:     "static/img/portfolio/fullsize/1.jpg",
	// 	Autor:     "Comuna da Catarina",
	// 	ID:        primitive.NewObjectID(),
	// 	Data:      time.Now()}

	// err := createPost(pp)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if err != nil {
		template.Must(template.ParseFiles(NotFoundTemplate))
	}
	post := template.Must(template.ParseFiles(PostTemplate))
	post.Execute(w, data)
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

func getByID(permalink string) ([]*Post, error) {
	filter := bson.D{
		primitive.E{Key: "Object_ID", Value: permalink},
	}

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
