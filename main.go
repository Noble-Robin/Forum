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

type Categorie struct {
	ID          int
	Title       string
	Description string
	View        int
}

type Thread struct {
	ID    int
	Title string
}

type Post struct {
	ID        int
	Username  string
	Content   string
	CreatedAt string
}

func Home(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	user := User{
		IsLoggedIn: session.Values["authenticated"] == true,
		Username:   fmt.Sprintf("%v", session.Values["username"]),
	}

	rows, err := db.Query("SELECT id, title, description, view FROM Categories ORDER BY view DESC")
	if err != nil {
		log.Printf("Error querying categories: %v", err)
		http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	categories := []Categorie{}
	for rows.Next() {
		var categorie Categorie
		if err := rows.Scan(&categorie.ID, &categorie.Title, &categorie.Description, &categorie.View); err != nil {
			log.Printf("Error scanning category: %v", err)
			http.Error(w, "Error reading category", http.StatusInternalServerError)
			return
		}
		categories = append(categories, categorie)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		http.Error(w, "Error reading category", http.StatusInternalServerError)
		return
	}

	data := struct {
		User       User
		Categories []Categorie
	}{
		User:       user,
		Categories: categories,
	}

	tmpl, err := template.ParseFiles("tmpl/home.html")
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

func main() {
	var err error
	db, err = sql.Open("sqlite3", "C:\\Users\\JENGO\\Forum\\sqlite\\data.db")
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
	// http.HandleFunc("/category", Category)
	http.HandleFunc("/threads", Threads)
	http.HandleFunc("/posts", Posts)

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

// func Category(w http.ResponseWriter, r *http.Request) {
//     rows, err := db.Query("SELECT id, title, description, view FROM Categories")
//     if err != nil {
//         http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
//         return
//     }
//     defer rows.Close()

//     categories := []Categorie{}
//     for rows.Next() {
//         var categorie Categorie
//         if err := rows.Scan(&categorie.ID, &categorie.Title, &categorie.Description, &categorie.View); err != nil {
//             http.Error(w, "Error reading category", http.StatusInternalServerError)
//             return
//         }
//         categories = append(categories, categorie)
//     }

//     http.Redirect(w, r, "/home", http.StatusFound)
// }

func Threads(w http.ResponseWriter, r *http.Request) {
	categorieID := r.URL.Query().Get("categorie_id")

	rows, err := db.Query("SELECT id, title FROM threads WHERE category_id = ?", categorieID)
	if err != nil {
		http.Error(w, "Error retrieving threads", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	threads := []Thread{}
	for rows.Next() {
		var thread Thread
		if err := rows.Scan(&thread.ID, &thread.Title); err != nil {
			http.Error(w, "Error reading threads", http.StatusInternalServerError)
			return
		}
		threads = append(threads, thread)
	}

	tmpl := template.Must(template.ParseFiles("tmpl/threads.html"))
	tmpl.Execute(w, threads)
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
	if r.Method == http.MethodPost {
		forumID := r.FormValue("forum_id")
		title := r.FormValue("title")
		userID := r.FormValue("user_id") // Assurez-vous que l'utilisateur est connecté et récupérez son ID

		db, err := sql.Open("sqlite3", "C:/Users/JENGO/Forum/sqlite/data.db")
		if err != nil {
			http.Error(w, "Erreur lors de l'ouverture de la base de données", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		_, err = db.Exec("INSERT INTO threads (forum_id, title, user_id) VALUES (?, ?, ?)", forumID, title, userID)
		if err != nil {
			http.Error(w, "Erreur lors de la création du thread", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/threads?forum_id=%s", forumID), http.StatusFound)
	} else {
		http.ServeFile(w, r, "tmpl/create_thread.html")
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		threadID := r.FormValue("thread_id")
		content := r.FormValue("content")
		userID := r.FormValue("user_id")

		db, err := sql.Open("sqlite3", "C:/Users/JENGO/Forum/sqlite/data.db")
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
	} else {
		http.ServeFile(w, r, "tmpl/create_post.html")
	}
}
