package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/http/cgi"
	"regexp"
	"strings"
	"time"

	// драйвер mysql для работы с mariadb
	_ "github.com/go-sql-driver/mysql"
)

// entity из бд
type Application struct {
	Surname   string `json:"surname"`
	Name      string `json:"name"`
	Midname   string `json:"midname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Birthdate string `json:"bday"`
	Gender    int    `json:"gender"`
	Favlangs  []int  `json:"favlangs"`
	Bio       string `json:"bio"`
}

func main() {
	// авторизация по данным. пока без .env
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
	// откидываем ненужные для бэкенда данные из пути
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
		var langcount int
		err = db.QueryRow("SELECT COUNT(*) FROM languages;").Scan(&langcount)

		if err != nil {
			http.Error(w, "Ошибка обращения к базе данных", http.StatusInternalServerError)
			return
		}
		// валидация данных
		if err := validate(app, langcount); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		// используем транзакции, чтобы исключить проблему с несоотнесённостью applications с languages связью многие-ко-многим
		//
		// если же не сделать этого, получим ситуацию: заявка есть в applications, но будет отсутствовать хотя бы одна запись
		// в applications_languages
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, `{"error": "Transaction failed"}`, http.StatusInternalServerError)
			return
		}
		query := "INSERT INTO applications (surname, name,midname, phone_number, email, birth_date, gender, bio) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
		w.Header().Set("Content-Type", "application/json")
		res, err := tx.Exec(
			query,
			// защитил сайт от xss атак, проведя sanitize через html.EscapeString

			// если этого не сделать: <script>вредоносный код</script> выполнится в полях без проверки
			html.EscapeString(app.Surname), html.EscapeString(app.Name), html.EscapeString(app.Midname), app.Phone, app.Email, app.Birthdate, app.Gender, html.EscapeString(app.Bio),
		)
		if err != nil {
			tx.Rollback()
			log.Printf("Insert error: %v", err)
			http.Error(w, `{"error": "Failed to save application"}`, http.StatusInternalServerError)
			return
		}
		// id добавленной заявки в бд
		lastId, _ := res.LastInsertId()
		langQuery := `INSERT INTO application_language (application_id, language_id) VALUES (?, ?)`
		for _, langID := range app.Favlangs {
			_, err := tx.Exec(langQuery, lastId, langID)
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

// валидация перед вводом в бд
func validate(app Application, langcount int) error {
	if strings.TrimSpace(app.Surname) == "" || len(app.Surname) > 128 {
		return errors.New("Фамилия обязательна и не должна превышать 128 символов")
	}
	if strings.TrimSpace(app.Name) == "" || len(app.Name) > 128 {
		return errors.New("Имя обязательно и не должно превышать 128 символов")
	}
	if strings.TrimSpace(app.Midname) != "" && len(app.Midname) > 128 {
		return errors.New("Отчество не должно превышать 128 символов")
	}

	// используем регулярные выражения
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(strings.ToLower(app.Email)) {
		return errors.New("Некорректный формат email")
	}

	app.Phone = strings.TrimSpace(app.Phone)
	phoneRegex := regexp.MustCompile(`^(\+?[0-9\-\(\)\s]{10,20})$`)
	if !phoneRegex.MatchString(app.Phone) {
		return errors.New("Некорректный формат номера телефона")
	}
	if app.Birthdate == "" {
		return errors.New("дата рождения не указана")
	}
	birthDate, err := time.Parse("2006-01-02", app.Birthdate)
	if err != nil {
		return errors.New("Введите корректную дату. Формат YYYY-MM-DD")
	}
	if birthDate.After(time.Now()) {
		return errors.New("Дата не может быть в будущем")
	}
	if time.Since(birthDate).Hours() > 24*365*150 {
		return errors.New("Указан слишком большой возраст")
	}
	if len(app.Favlangs) == 0 {
		return errors.New("выберите хотя бы один язык программирования")
	}
	for _, favlang := range app.Favlangs {
		if favlang < 0 || favlang > langcount {
			return errors.New("Выберите корректный любимый язык программирования")
		}
	}

	if strings.TrimSpace(app.Bio) == "" {
		return errors.New("заполните поле 'О себе'")
	}
	return nil
}
