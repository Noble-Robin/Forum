func Threads(w http.ResponseWriter, r *http.Request) {
	forumID := r.URL.Query().Get("forum_id")

	db, err := sql.Open("sqlite3", "C:/Users/JENGO/Forum/sqlite/data.db")
	if err != nil {
		http.Error(w, "Erreur lors de l'ouverture de la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title FROM threads WHERE categorie_id = ?", forumID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des threads", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	threads := []Thread{}
	for rows.Next() {
		var thread Thread
		if err := rows.Scan(&thread.ID, &thread.Title); err != nil {
			http.Error(w, "Erreur lors de la lecture des threads", http.StatusInternalServerError)
			return
		}
		threads = append(threads, thread)
	}

	tmpl := template.Must(template.ParseFiles("tmpl/threads.html"))
	tmpl.Execute(w, threads)
}

type Thread struct {
	ID    int
	Title string
}
