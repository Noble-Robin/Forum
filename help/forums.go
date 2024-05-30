func Category(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "C:/Users/JENGO/Forum/sqlite/data.db")
	if err != nil {
		http.Error(w, "Erreur lors de l'ouverture de la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, description FROM forums")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des forums", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	forums := []Forum{}
	for rows.Next() {
		var forum Forum
		if err := rows.Scan(&forum.ID, &forum.Title, &forum.Description); err != nil {
			http.Error(w, "Erreur lors de la lecture des forums", http.StatusInternalServerError)
			return
		}
		forums = append(forums, forum)
	}

	tmpl := template.Must(template.ParseFiles("tmpl/forums.html"))
	tmpl.Execute(w, forums)
}

type Categorie struct {
	ID          int
	Title       string
	Description string
}
