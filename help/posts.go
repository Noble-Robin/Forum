func Posts(w http.ResponseWriter, r *http.Request) {
	threadID := r.URL.Query().Get("thread_id")

	db, err := sql.Open("sqlite3", "C:/Users/JENGO/Forum/sqlite/data.db")
	if err != nil {
		http.Error(w, "Erreur lors de l'ouverture de la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT p.id, u.username, p.content, p.created_at FROM posts p JOIN users u ON p.user_id = u.id WHERE thread_id = ?", threadID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Content, &post.CreatedAt); err != nil {
			http.Error(w, "Erreur lors de la lecture des posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	tmpl := template.Must(template.ParseFiles("tmpl/posts.html"))
	tmpl.Execute(w, posts)
}

type Post struct {
	ID        int
	Username  string
	Content   string
	CreatedAt string
}
