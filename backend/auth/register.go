package auth

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	Errorhandel "forum/backend/Errors"
	"forum/backend/database"
	"forum/backend/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		Errorhandel.Errordirect(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	tmpt, err := template.ParseFiles("./frontend/templates/register.html")
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tmpt.Execute(w, nil)
}

func RegisterHandlerPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Errorhandel.Errordirect(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	tmpt, err := template.ParseFiles("./frontend/templates/register.html")
	if err != nil {
		log.Printf("template parse error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var input models.LoginData
	input.Email = r.FormValue("email")
	input.Username = r.FormValue("username")
	input.Password = r.FormValue("password")
	confirmPassword := r.FormValue("confirmpassword")

	if err := validateValues(input.Email, input.Username, input.Password, confirmPassword); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		tmpt.Execute(w, models.ErrorData{Error: err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("password hash error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		tmpt.Execute(w, models.ErrorData{Error: "Internal server error, try again later"})
		return
	}

	_, err = database.Db.Exec(`
		INSERT INTO users(email, username, password_hash)
		VALUES(?, ?, ?)`,
		input.Email, input.Username, hashedPassword,
	)
	if err != nil {
		log.Printf("DB insert error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		tmpt.Execute(w, models.ErrorData{Error: "Internal server error, try again later"})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func validateValues(email, username, password, confirmPassword string) error {
	email = strings.TrimSpace(email)

	var existingEmail string
	err := database.Db.QueryRow(`SELECT email FROM users WHERE email = ?`, email).Scan(&existingEmail)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("email DB error: %v", err)
			return fmt.Errorf("internal server error, try later")
		}
	} else {
		return fmt.Errorf("this email is already used")
	}

	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched := regexp.MustCompile(regex).MatchString(email); !matched {
		return fmt.Errorf("invalid email format")
	}

	if len(username) < 4 || len(username) > 20 {
		return fmt.Errorf("username must be between 4 and 20 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(username) {
		return fmt.Errorf("username can only use letters, numbers, - and _")
	}

	var existingUsername string
	err = database.Db.QueryRow(`SELECT username FROM users WHERE username = ?`, username).Scan(&existingUsername)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("username DB error: %v", err)
			return fmt.Errorf("internal server error, try later")
		}
	} else {
		return fmt.Errorf("this username is already used")
	}

	if len(password) < 8 || len(password) > 60 {
		return fmt.Errorf("password must be between 8 and 60 characters")
	}
	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	var hasLetter, hasNumber bool
	for _, c := range password {
		switch {
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z'):
			hasLetter = true
		case c >= '0' && c <= '9':
			hasNumber = true
		}
	}
	if !hasLetter || !hasNumber {
		return fmt.Errorf("password must contain at least one letter and one number")
	}

	return nil
}
