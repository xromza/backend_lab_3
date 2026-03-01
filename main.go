package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cgi"

	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	Surname   string `json:"surname"`
	Name      string `json:"name"`
	Midname   string `json:"midname"`
	Phone     string `json:"tel"`
	Email     string `json:"email"`
	Birthdate string `json:"bday"`
	Gender    int    `json:"gender"`
	Favlangs  []int  `json:"favlangs"`
	Bio       string `json:"bio"`
}

func main() {

	dsn := "u82186:1169903@tcp(localhost:3306)/u82186"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, "Это главная страница CGI")
	})

	mux.HandleFunc("/save", saveHandler(db))

	cgi.Serve(http.StripPrefix("/backend_lab_3/backend.cgi", mux))
}

func saveHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var app Application
		err := json.NewDecoder(r.Body).Decode(&app)
		if err != nil {
			http.Error(w, "Ошибка парсинга JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, `{"error": "Transaction failed"}`, http.StatusInternalServerError)
			return
		}
		query := "INSERT INTO applications (surname, name,midname, phone, email, bday, gender, bio) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
		w.Header().Set("Content-Type", "application/json")
		res, err := tx.Exec(
			query,
			app.Surname, app.Name, app.Midname, app.Phone, app.Email, app.Birthdate, app.Gender, app.Bio,
		)
		if err != nil {
			tx.Rollback()
			log.Printf("Insert error: %v", err)
			http.Error(w, `{"error": "Failed to save application"}`, http.StatusInternalServerError)
			return
		}
		lastId, _ := res.LastInsertId()
		langQuery := `INSERT INTO application_language (application_id, language_id) VALUES (?, ?)`
		for _, langID := range app.Favlangs {
			_, err := db.Exec(langQuery, lastId, langID)
			if err != nil {
				tx.Rollback()
				http.Error(w, `{"error": "Failed to save languages"}`, http.StatusInternalServerError)
				return
			}
		}
		if err := tx.Commit(); err != nil {
			http.Error(w, `{"error": "Commit failed"}`, http.StatusInternalServerError)
			return
		}
		response := map[string]string{
			"message": fmt.Sprintf("Заявка #%d успешно сохранена", lastId),
			"status":  "success",
		}

		json.NewEncoder(w).Encode(response)
	}
}
