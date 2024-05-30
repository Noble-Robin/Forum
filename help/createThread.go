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
