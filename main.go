package main

import (
	"fmt"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/home.html")
}

func Login(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/Login.html")
}

func StaticFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func imgFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/home", Home)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/static/", StaticFiles)
	http.HandleFunc("/img/", imgFiles)

	fmt.Print("http://localhost:8081/home")
	http.ListenAndServe(":8081", nil)
}
