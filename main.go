package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db    *sql.DB
	store = sessions.NewCookieStore([]byte("super-secret-key"))
)

type User struct {
	Username       string
	IsLoggedIn     bool
	ProfilePicture string
}

type CategoryData struct {
	// Add fields as necessary
}

type Forum struct {
	ID          int
	Title       string
	Description string
}

func Home(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	user := User{
		IsLoggedIn: session.Values["authenticated"] == true,
		Username:   fmt.Sprintf("%v", session.Values["username"]),
	}

	tmpl, err := template.ParseFiles("tmpl/home.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, user)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		loginHandler(w, r)
		return
	}
	http.ServeFile(w, r, "tmpl/Login.html")
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createUser(w, r)
		return
	}
	http.ServeFile(w, r, "tmpl/register.html")
}

func ErrorPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/error.html")
}

func Connect(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/connect.html")
}

func StaticFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func ImgFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func Feed(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/feed.html")
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "C:/Users/JENGO/Forum/sqlite/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/home", Home)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/error", ErrorPage)
	http.HandleFunc("/feed", Feed)
	http.HandleFunc("/connect", Connect)
	http.HandleFunc("/logout", Deconnect)
	http.HandleFunc("/static/", StaticFiles)
	http.HandleFunc("/img/", ImgFiles)

	db, err := sql.Open("sqlite3", "C:/Users/JENGO/Forum/sqlite/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//files := []string{"User.sql", "thread.sql", "post.sql"}
	//for _, file := range files {
	//	sqlFile, err := ioutil.ReadFile("sqlite/" + file)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	_, err = db.Exec(string(sqlFile))
	//	if err != nil {
	//	log.Fatal(err)
	//	}
	//	log.Printf("File %s executed successfully", file)
	//}

	fmt.Println("Server started at http://localhost:8081/home")
	http.ListenAndServe(":8081", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if connectDB(username, password) {
		session, _ := store.Get(r, "session-name")
		session.Values["authenticated"] = true
		session.Values["username"] = username
		session.Save(r, w)
		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		http.Redirect(w, r, "/error", http.StatusFound)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if !verifDB(username, email) {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}

	_, err := db.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", username, password, email)
	if err != nil {
		http.Error(w, "Error inserting user into the database", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func verifDB(username, email string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? OR email = ?", username, email).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}

	return count == 0
}

func connectDB(username, password string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND password = ?", username, password).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}

	return count > 0
}

func Deconnect(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/home", http.StatusFound)
}
