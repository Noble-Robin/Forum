package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/home.html")
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		loginHandler(w, r)
		return
	}
	http.ServeFile(w, r, "tmpl/Login.html")
}

func Error1(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/error.html")
}

func Connect(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/connect.html")
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
	http.HandleFunc("/error", Error1)
	http.HandleFunc("/connect", Connect)
	http.HandleFunc("/static/", StaticFiles)
	http.HandleFunc("/img/", imgFiles)

	db, err := sql.Open("sqlite3", "C:/Users/JENGO/Forum/sqlite/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//files := []string{"User.sql", "conversation.sql", "contact.sql", "message.sql", "group.sql"}
	//for _, file := range files {
	//	sqlFile, err := ioutil.ReadFile("sqlite/" + file)
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	//	_, err = db.Exec(string(sqlFile))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("File %s executed successfully", file)
	//}

	fmt.Print("http://localhost:8081/home")
	http.ListenAndServe(":8081", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if isValidUser(username, password) {
		http.Redirect(w, r, "/connect", http.StatusFound)
		fmt.Print("yes")
	} else {
		http.Redirect(w, r, "/error", http.StatusFound)
	}
}

func isValidUser(username, password string) bool {
	return (username == "admin" && password == "password")
}
