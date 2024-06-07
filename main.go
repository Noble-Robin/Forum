package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db       *sql.DB
	sessions = map[string]string{} // simple in-memory session storage
)

type User struct {
	Username       string
	IsLoggedIn     bool
	ProfilePicture string
}

type Categorie struct {
	ID          int
	Title       string
	Description string
	View        int
}

type Thread struct {
	ID            int
	Title         string
	CategoryTitle string
	UserUsername  string
	CreatedAt     time.Time
}
type Post struct {
	ID        int
	Username  string
	Content   string
	CreatedAt string
}

func Home(w http.ResponseWriter, r *http.Request) {
	user := getUserFromSession(r)

	categories, err := getCategories()
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
		return
	}

	data := struct {
		User       User
		Categories []Categorie
	}{
		User:       user,
		Categories: categories,
	}

	renderTemplate(w, "tmpl/home.html", data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		loginHandler(w, r)
		return
	}
	http.ServeFile(w, r, "tmpl/Login.html")
}

func forums(w http.ResponseWriter, r *http.Request) {
	user := getUserFromSession(r)

	categories, err := getCategories()
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
		return
	}

	data := struct {
		User       User
		Categories []Categorie
	}{
		User:       user,
		Categories: categories,
	}

	renderTemplate(w, "tmpl/forums.html", data)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createUser(w, r)
		return
	}
	http.ServeFile(w, r, "tmpl/register.html")
}

func ct(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		log.Printf("Error getting session cookie: %v", err)
	}

	user := User{}

	if sessionID != nil {
		username, ok := sessions[sessionID.Value]
		if ok {
			user = User{
				IsLoggedIn: true,
				Username:   username,
			}
		}
	}

	rows, err := db.Query(`
		SELECT t.id, t.title, t.categorie_title, t.user_username, t.created_at 
		FROM threads t 
		LEFT JOIN users u ON t.user_username = u.username 
		LEFT JOIN categories c ON t.categorie_title = c.title`)
	if err != nil {
		log.Printf("Error querying threads: %v", err)
		http.Error(w, "Error retrieving threads", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	threads := []Thread{}

	for rows.Next() {
		var thread Thread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.CategoryTitle, &thread.UserUsername, &thread.CreatedAt); err != nil {
			log.Printf("Error scanning threads: %v", err)
			http.Error(w, "Error reading threads", http.StatusInternalServerError)
			return
		}
		threads = append(threads, thread)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		http.Error(w, "Error reading threads", http.StatusInternalServerError)
		return
	}

	data := struct {
		User    User
		Threads []Thread
	}{
		User:    user,
		Threads: threads,
	}

	tmpl, err := template.ParseFiles("tmpl/thread.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
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

// func Category(w http.ResponseWriter, r *http.Request) { //les users ne peuvent pas créer de catégorie
// 	sessionID, err := r.Cookie("session_id")
// 	user := User{}

// 	if err == nil {
// 		username, ok := sessions[sessionID.Value]
// 		if ok {
// 			user = User{
// 				IsLoggedIn: true,
// 				Username:   username,
// 			}
// 		}
// 	}

// 	rows, err := db.Query("SELECT id, title, description, view FROM Categories ORDER BY view DESC")
// 	if err != nil {
// 		log.Printf("Error querying categories: %v", err)
// 		http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	categories := []Categorie{}
// 	for rows.Next() {
// 		var categorie Categorie
// 		if err := rows.Scan(&categorie.ID, &categorie.Title, &categorie.Description, &categorie.View); err != nil {
// 			log.Printf("Error scanning category: %v", err)
// 			http.Error(w, "Error reading category", http.StatusInternalServerError)
// 			return
// 		}
// 		categories = append(categories, categorie)
// 	}

// 	if err = rows.Err(); err != nil {
// 		log.Printf("Error iterating over rows: %v", err)
// 		http.Error(w, "Error reading category", http.StatusInternalServerError)
// 		return
// 	}

// 	data := struct {
// 		User       User
// 		Categories []Categorie
// 	}{
// 		User:       user,
// 		Categories: categories,
// 	}

// 	tmpl, err := template.ParseFiles("tmpl/category.html")
// 	if err != nil {
// 		log.Printf("Error parsing template: %v", err)
// 		http.Error(w, "Error parsing template", http.StatusInternalServerError)
// 		return
// 	}

// 	err = tmpl.Execute(w, data)
// 	if err != nil {
// 		log.Printf("Error executing template: %v", err)
// 		http.Error(w, "Error executing template", http.StatusInternalServerError)
// 	}

// }

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./sqlite/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/home", Home)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/error", ErrorPage)
	http.HandleFunc("/connect", Connect)
	http.HandleFunc("/logout", Deconnect)
	http.HandleFunc("/static/", StaticFiles)
	http.HandleFunc("/img/", ImgFiles)
	// http.HandleFunc("/category/", Category)
	http.HandleFunc("/threads", Threads)
	http.HandleFunc("/posts", Posts)
	http.HandleFunc("/create-thread", CreateThread)
	http.HandleFunc("/forums", forums)
	http.HandleFunc("/thread", ct)

	// files := []string{"User.sql", "thread.sql", "post.sql", "Categorie.sql"}
	// for _, file := range files {
	// 	sqlFile, err := ioutil.ReadFile("sqlite/" + file)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	_, err = db.Exec(string(sqlFile))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Printf("File %s executed successfully", file)
	// }

	fmt.Println("Server started at http://localhost:8081/home")
	http.ListenAndServe(":8081", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if connectDB(username, password) {
		sessionID := fmt.Sprintf("%d", rand.Int())
		sessions[sessionID] = username

		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: sessionID,
			Path:  "/",
		})

		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		http.Redirect(w, r, "/error", http.StatusFound)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	name := r.FormValue("name")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if !verifDB(username, email) {
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}

	_, err := db.Exec("INSERT INTO users (username, name, email, password) VALUES (?, ?, ?, ?)", username, name, email, password)
	if err != nil {
		http.Error(w, "Error inserting user into the database", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {

	threadID := r.FormValue("thread_id")
	content := r.FormValue("content")
	userID := r.FormValue("user_id")

	db, err := sql.Open("sqlite3", "./sqlite/data.db")
	if err != nil {
		http.Error(w, "Erreur lors de l'ouverture de la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO posts (thread_id, user_id, content) VALUES (?, ?, ?)", threadID, userID, content)
	if err != nil {
		http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts?thread_id=%s", threadID), http.StatusFound)

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
	cookie, err := r.Cookie("session_id")
	if err == nil {
		delete(sessions, cookie.Value)
		http.SetCookie(w, &http.Cookie{
			Name:   "session_id",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func Threads(w http.ResponseWriter, r *http.Request) {
	categorieTitle := r.URL.Query().Get("categorie_title")
	if categorieTitle == "" {
		http.Error(w, "Missing category title", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT id, title, categorie_title, user_username FROM threads WHERE categorie_title = ?", categorieTitle)
	if err != nil {
		log.Printf("Error querying threads: %v", err)
		http.Error(w, "Error retrieving threads", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	threads := []Thread{}
	for rows.Next() {
		var thread Thread
		if err := rows.Scan(&thread.ID, &thread.Title, &thread.CategoryTitle, &thread.UserUsername); err != nil {
			log.Printf("Error scanning threads: %v", err)
			http.Error(w, "Error reading threads", http.StatusInternalServerError)
			return
		}
		threads = append(threads, thread)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over threads: %v", err)
		http.Error(w, "Error reading threads", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("tmpl/thread.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Threads []Thread
	}{
		Threads: threads,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func Posts(w http.ResponseWriter, r *http.Request) {
	threadID := r.URL.Query().Get("thread_id")

	rows, err := db.Query("SELECT p.id, u.username, p.content, p.created_at FROM posts p JOIN users u ON p.user_id = u.id WHERE thread_id = ?", threadID)
	if err != nil {
		http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Content, &post.CreatedAt); err != nil {
			http.Error(w, "Error reading posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	tmpl := template.Must(template.ParseFiles("tmpl/posts.html"))
	tmpl.Execute(w, posts)
}

func CreateThread(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	username, ok := sessions[sessionID.Value]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	var userID string
	err = db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'ID utilisateur: %v", err)
		http.Error(w, "Erreur lors de la récupération de l'ID utilisateur", http.StatusInternalServerError)
		return
	}

	categorieTitle := r.FormValue("categorie_title")
	title := r.FormValue("title")

	_, err = db.Exec("INSERT INTO threads (categorie_title, title, user_username) VALUES (?, ?, ?)", categorieTitle, title, username)
	if err != nil {
		log.Printf("Erreur lors de la création du thread: %v", err)
		http.Error(w, "Erreur lors de la création du thread", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, ("/thread"), http.StatusFound)
}

func getUserFromSession(r *http.Request) User {
	sessionID, err := r.Cookie("session_id")
	user := User{}

	if err == nil {
		username, ok := sessions[sessionID.Value]
		if ok {
			user = User{
				IsLoggedIn: true,
				Username:   username,
			}
		}
	}

	return user
}

func getCategories() ([]Categorie, error) {
	rows, err := db.Query("SELECT id, title, description, view FROM Categories ORDER BY view DESC")
	if err != nil {
		return nil, fmt.Errorf("Error querying categories: %v", err)
	}
	defer rows.Close()

	categories := []Categorie{}
	for rows.Next() {
		var categorie Categorie
		if err := rows.Scan(&categorie.ID, &categorie.Title, &categorie.Description, &categorie.View); err != nil {
			return nil, fmt.Errorf("Error scanning category: %v", err)
		}
		categories = append(categories, categorie)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}

	return categories, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
