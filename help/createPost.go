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
