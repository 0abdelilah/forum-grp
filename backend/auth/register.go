package auth

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	"forum/backend/database"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandlerGet(w http.ResponseWriter, r *http.Request) {
	tmpt, err := template.ParseFiles("./frontend/templates/register.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	tmpt.Execute(w, nil)
}

func RegisterHandlerPost(w http.ResponseWriter, r *http.Request) {
	tmpt, err := template.ParseFiles("./frontend/templates/register.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	Email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmpassword")

	err = validateValues(Email, username, password, confirmPassword)
	if err != nil {
		tmpt.Execute(w, struct{ Error string }{Error: strings.Title(err.Error())})
		fmt.Println(err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tmpt.Execute(w, struct{ Error string }{Error: "Internal server error, try again later"})
		log.Println("Error hashing password:", err)
		return
	}

	_, err = database.Db.Exec(`
	INSERT INTO users(email, username, password_hash)
	VALUES(?, ?, ?)`,
		Email, username, hashedPassword,
	)
	if err != nil {
		tmpt.Execute(w, struct{ Error string }{Error: "Internal server error, try again later"})
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func validateValues(email, username, password, confirmPassword string) error {
	email = strings.TrimSpace(email)
	// --- Check if email exists ---
	var existingEmail string
	err := database.Db.QueryRow(`SELECT email FROM users WHERE email = ?`, email).Scan(&existingEmail)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("email database error: %v\n", err)
			return fmt.Errorf("internal server error, try later")
		}
	} else {
		return fmt.Errorf("this email is already used")
	}

	// --- Validate email format ---
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched := regexp.MustCompile(regex).MatchString(email); !matched {
		return fmt.Errorf("invalid email format")
	}

	// --- Validate username ---
	if len(username) < 4 || len(username) > 20 {
		return fmt.Errorf("username must be between 4 and 20 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(username) {
		return fmt.Errorf("The username must contain only letters, numbers, hyphens and underscores.")
	}

	// --- Check if username exists ---
	var existingUsername string
	err = database.Db.QueryRow(`SELECT username FROM users WHERE username = ?`, username).Scan(&existingUsername)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("username database error: %v\n", err)
			return fmt.Errorf("internal server error, try later")
		}
	} else {
		return fmt.Errorf("this username is already used")
	}

	// --- Validate password ---
	if len(password) < 8 || len(password) > 60 {
		return fmt.Errorf("password must be between 8 and 60 characters")
	}
	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}
	var hasLetter, hasNumber bool
	for _, c := range password {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		} else if c >= '0' && c <= '9' {
			hasNumber = true
		}
	}
	if !hasLetter || !hasNumber {
		return fmt.Errorf("password must contain at least one letter and one number")
	}

	return nil
}
