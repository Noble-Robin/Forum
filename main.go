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

func main() {
	http.HandleFunc("/home", Home)
	http.HandleFunc("/login", Login)

	fmt.Print("http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}
