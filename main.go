package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	db       *sql.DB
	sessions = map[string]string{}
)

type User struct {
	ID             int
	Username       string
	Name           string
	Email          string
	Role           UserRole
	IsLoggedIn     bool
	ProfilePicture string
}

type UserRole int

const (
	RoleGuest UserRole = iota
	RoleUser
	RoleModerator
	RoleAdministrator
)

type Categorie struct {
	ID          int
	Title       string
	Description string
	Post        int
	Threads     []Thread
}


type Thread struct {
	ID            int
	Title         string
	CategoryTitle string
	UserUsername  string
	CreatedAt     time.Time
	Posts         []Post
}
type Post struct {
	ID        int
	Username  string
	Content   string
	CreatedAt time.Time
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
	user := getUserFromSession(r)

	categories, err := getCategories()
	if err != nil {
		log.Printf("Error retrieving categories: %v", err)
		http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
		return
	}

	// Initialiser une carte de threads par catégorie
	categoriesMap := make(map[string][]Thread)

	// Pour chaque catégorie, récupérer les threads associés
	for _, cat := range categories {
		categoryTitle := cat.Title

		// Récupérer les threads de la catégorie depuis la base de données
		rows, err := db.Query(`
            SELECT t.id, t.title, t.categorie_title, t.user_username, t.created_at 
            FROM threads t 
            WHERE t.categorie_title = ?
        `, categoryTitle)
		if err != nil {
			log.Printf("Error querying threads for category %s: %v", categoryTitle, err)
			http.Error(w, "Error retrieving threads", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var threads []Thread
		for rows.Next() {
			var thread Thread
			if err := rows.Scan(&thread.ID, &thread.Title, &thread.CategoryTitle, &thread.UserUsername, &thread.CreatedAt); err != nil {
				log.Printf("Error scanning threads for category %s: %v", categoryTitle, err)
				http.Error(w, "Error reading threads", http.StatusInternalServerError)
				return
			}

			// Récupérer les posts pour chaque thread
			postsRows, err := db.Query(`
                SELECT p.id, u.username, p.content, p.created_at 
                FROM posts p 
                JOIN users u ON p.user_id = u.id 
                WHERE p.thread_id = ?
            `, thread.ID)
			if err != nil {
				log.Printf("Error querying posts for thread %d: %v", thread.ID, err)
				http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
				return
			}
			defer postsRows.Close()

			var posts []Post
			for postsRows.Next() {
				var post Post
				if err := postsRows.Scan(&post.ID, &post.Username, &post.Content, &post.CreatedAt); err != nil {
					log.Printf("Error scanning posts for thread %d: %v", thread.ID, err)
					http.Error(w, "Error reading posts", http.StatusInternalServerError)
					return
				}
				posts = append(posts, post)
			}
			thread.Posts = posts

			threads = append(threads, thread)
		}

		categoriesMap[categoryTitle] = threads
	}

	var categoriesData []Categorie
	for _, cat := range categories {
		categoriesData = append(categoriesData, Categorie{
			ID:          cat.ID,
			Title:       cat.Title,
			Description: cat.Description,
			Threads:     categoriesMap[cat.Title],
		})
	}

	data := struct {
		User       User
		Categories []Categorie
	}{
		User:       user,
		Categories: categoriesData,
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
		return
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {
	user := getUserFromSession(r)
	
	if !user.IsLoggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
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


	renderTemplate(w, "tmpl/profile.html", data)
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
	http.HandleFunc("/create-thread", CreateThread)
	http.HandleFunc("/forums", forums)
	http.HandleFunc("/thread", ct)
	http.HandleFunc("/create-post", CreatePost)
	http.HandleFunc("/create-category", CreateCategories)
	http.HandleFunc("/profile", Profile)
	http.HandleFunc("/update-profile", UpdateProfile)
	http.HandleFunc("/delete_thread", DeleteThread)

	user := User{
		ID:             1,
		Username:       "john_doe",
		Name:           "John Doe",
		Email:          "john.doe@example.com",
		IsLoggedIn:     true,
		ProfilePicture: "/path/to/profile_picture.jpg",
	}
	role := RoleUser
	if user.Username == "admin" {
		role = RoleAdministrator
	} else if user.Username == "moderator" {
		role = RoleModerator
	} else if user.Username == "" {
		role = RoleGuest
	}

	fmt.Printf("User: %s, Role: %d\n", user.Username, role)

	files := []string{"User.sql", "thread.sql", "post.sql", "Categorie.sql", "update.sql","report.sql"}
	for _, file := range files {
		sqlFile, err := ioutil.ReadFile("sqlite/" + file)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(string(sqlFile))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("File %s executed successfully", file)
	}

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

	hashedPassword, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, name, email, password) VALUES (?, ?, ?, ?)", username, name, email, hashedPassword)
	if err != nil {
		http.Error(w, "Error inserting user into the database", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func connectDB(username, password string) bool {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		log.Println(err)
		return false
	}

	return CheckPasswordHash(password, hashedPassword)
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

func CreateCategories(w http.ResponseWriter, r *http.Request) {
	user := getUserFromSession(r)
	if !user.IsLoggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

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

	title := r.FormValue("title")
	description := r.FormValue("description")

	_, err = db.Exec("INSERT INTO Categories (title, description, user_username) VALUES (?, ?, ?)", title, description, username)
	if err != nil {
		log.Printf("Error creating category: %v", err)
		http.Error(w, "Error creating category", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/thread", http.StatusFound)
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
		log.Printf("Error retrieving user ID: %v", err)
		http.Error(w, "Error retrieving user ID", http.StatusInternalServerError)
		return
	}

	categoryTitle := r.FormValue("categorie_title")
	title := r.FormValue("title")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, "Error beginning transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO threads (categorie_title, title, user_username) VALUES (?, ?, ?)", categoryTitle, title, username)
	if err != nil {
		tx.Rollback()
		log.Printf("Error creating thread: %v", err)
		http.Error(w, "Error creating thread", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE categories SET post = post + 1 WHERE title = ?", categoryTitle)
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating post count: %v", err)
		http.Error(w, "Error updating post count", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/thread", http.StatusFound)
}

func DeleteThread(w http.ResponseWriter, r *http.Request) {

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

	threadID := r.FormValue("thread_id")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, "Error beginning transaction", http.StatusInternalServerError)
		return
	}

	var dbUsername string
	err = tx.QueryRow("SELECT user_username FROM threads WHERE id = ?", threadID).Scan(&dbUsername)
	if err != nil {
		tx.Rollback()
		log.Printf("Error retrieving thread owner: %v", err)
		http.Error(w, "Error retrieving thread owner", http.StatusInternalServerError)
		return
	}

	if dbUsername != username {
		tx.Rollback()
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err = tx.Exec("DELETE FROM threads WHERE id = ?", threadID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting thread: %v", err)
		http.Error(w, "Error deleting thread", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("DELETE FROM comments WHERE thread_id = ?", threadID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting comments: %v", err)
		http.Error(w, "Error deleting comments", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/thread", http.StatusFound)
}

func ReportThread(w http.ResponseWriter, r *http.Request) {
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

	threadID := r.FormValue("thread_id")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, "Error beginning transaction", http.StatusInternalServerError)
		return
	}

	var dbUsername string
	err = tx.QueryRow("SELECT user_username FROM threads WHERE id = ?", threadID).Scan(&dbUsername)
	if err != nil {
		tx.Rollback()
		log.Printf("Error retrieving thread owner: %v", err)
		http.Error(w, "Error retrieving thread owner", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO reports (reporter_username, thread_id) VALUES (?, ?)", username, threadID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting report: %v", err)
		http.Error(w, "Error reporting thread", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/thread", http.StatusFound)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	user := getUserFromSession(r)
	if !user.IsLoggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	threadID := r.FormValue("thread_id")
	content := r.FormValue("content")

	var categoryTitle string
	err := db.QueryRow("SELECT categorie_title FROM threads WHERE id = ?", threadID).Scan(&categoryTitle)
	if err != nil {
		log.Printf("Error retrieving category title: %v", err)
		http.Error(w, "Error retrieving category title", http.StatusInternalServerError)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, "Error beginning transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO posts (thread_id, user_id, content) VALUES (?, ?, ?)", threadID, user.ID, content)
	if err != nil {
		tx.Rollback()
		log.Printf("Error creating post: %v", err)
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE categories SET post = post + 1 WHERE title = ?", categoryTitle)
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating post count: %v", err)
		http.Error(w, "Error updating post count", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/thread?thread_id=%s", threadID), http.StatusFound)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	user := getUserFromSession(r)
	if !user.IsLoggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	postID := r.FormValue("post_id")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, "Error beginning transaction", http.StatusInternalServerError)
		return
	}

	var threadID string
	err = tx.QueryRow("SELECT thread_id FROM posts WHERE id = ? AND user_id = ?", postID, user.ID).Scan(&threadID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error retrieving post information: %v", err)
		http.Error(w, "Error retrieving post information", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", postID, user.ID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting post: %v", err)
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE categories SET post = post - 1 WHERE title = (SELECT categorie_title FROM threads WHERE id = ?)", threadID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating post count: %v", err)
		http.Error(w, "Error updating post count", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/thread?thread_id=%s", threadID), http.StatusFound)
}

func ReportPost(w http.ResponseWriter, r *http.Request) {
	user := getUserFromSession(r)
	if !user.IsLoggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	postID := r.FormValue("post_id")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, "Error beginning transaction", http.StatusInternalServerError)
		return
	}

	var threadID string
	err = tx.QueryRow("SELECT thread_id FROM posts WHERE id = ?", postID).Scan(&threadID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error retrieving post information: %v", err)
		http.Error(w, "Error retrieving post information", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO reports (reporter_username, post_id) VALUES (?, ?)", user.Username, postID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting report: %v", err)
		http.Error(w, "Error reporting post", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/thread?thread_id=%s", threadID), http.StatusFound)
}

func getUserFromSession(r *http.Request) User {
	sessionID, err := r.Cookie("session_id")
	user := User{}

	if err == nil {
		username, ok := sessions[sessionID.Value]
		if ok {
			err := db.QueryRow("SELECT id, username, email,name FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Email, &user.Name)
			if err != nil {
				log.Printf("Error querying user: %v", err)
			} else {
				user.IsLoggedIn = true
			}
		}
	}

	return user
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := getUserFromSession(r)
	if !user.IsLoggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method == http.MethodPost {
		newUsername := r.FormValue("username")
		newProfilePicture := r.FormValue("profile_picture")

		_, err := db.Exec("UPDATE users SET username = ?, profile_picture = ? WHERE id = ?", newUsername, newProfilePicture, user.ID)
		if err != nil {
			log.Printf("Error updating user: %v", err)
			http.Error(w, "Error updating profile", http.StatusInternalServerError)
			return
		}

		sessions[newUsername] = sessions[user.Username]
		delete(sessions, user.Username)

		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	data := struct {
		User User
	}{
		User: user,
	}

	renderTemplate(w, "tmpl/profile.html", data)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func getCategories() ([]Categorie, error) {
	rows, err := db.Query("SELECT id, title, description FROM categories")
	if err != nil {
		return nil, fmt.Errorf("error querying categories: %v", err)
	}
	defer rows.Close()

	var categories []Categorie
	for rows.Next() {
		var cat Categorie
		if err := rows.Scan(&cat.ID, &cat.Title, &cat.Description); err != nil {
			return nil, fmt.Errorf("error scanning categories: %v", err)
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over categories: %v", err)
	}

	return categories, nil
}
