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

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createUser(w, r)
		return
	}
	http.ServeFile(w, r, "tmpl/register.html")
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
	http.HandleFunc("/register", register)
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
	if !connect_db(username, password) {
		http.Redirect(w, r, "/home", http.StatusFound)
		fmt.Print("yes")
	} else {
		http.Redirect(w, r, "/error", http.StatusFound)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if !verif_db(username, email) {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}

	db, err := sql.Open("sqlite3", "C:/Users/JENGO/Forum/sqlite/data.db")
	if err != nil {
		http.Error(w, "Erreur lors de l'ouverture de la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	insertQuery := "INSERT INTO users (username, password, email) VALUES (?, ?, ?)"
	_, err = db.Exec(insertQuery, username, password, email)
	if err != nil {
		http.Error(w, "Erreur lors de l'insertion de l'utilisateur dans la base de données", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Utilisateur créé avec succès: %s\n", username)
	http.Redirect(w, r, "/connect", http.StatusFound)
}
func verif_db(username, email string) bool {
	db, err := sql.Open("sqlite3", "C:/Users/JENGO/Forum/sqlite/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	var count int
	err = db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}

	if count > 0 {
		fmt.Printf("L'utilisateur %s ou l'email %s existe déjà\n", username, email)
		return false
	}

	return true
}

func connect_db(username, password string) bool {
	db, err := sql.Open("sqlite3", "C:/Users/JENGO/Forum/sqlite/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "SELECT COUNT(*) FROM users WHERE username = ? AND password = ?"
	var count int
	err = db.QueryRow(query, username, password).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}

	if count > 0 {
		fmt.Printf("L'utilisateur %s est connecter\n", username)
		return false
	}

	return true
}
