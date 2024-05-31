func Category(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, description, view FROM Categories")
	if err != nil {
		http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	categories := []Categorie{}
	for rows.Next() {
		var categorie Categorie
		if err := rows.Scan(&categorie.ID, &categorie.Title, &categorie.Description, &categorie.View); err != nil {
			http.Error(w, "Error reading category", http.StatusInternalServerError)
			return
		}
		categories = append(categories, categorie)
	}

	http.Redirect(w, r, "/home", http.StatusFound)
}

type Categorie struct {
	ID          int
	Title       string
	Description string
	View        int
}
