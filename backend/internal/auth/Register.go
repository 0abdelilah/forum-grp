package handlers

import (
	"log"
	"net/http"
	"regexp"

	"forum/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "frontend/templates/register.html")
		return
	case "POST":
		Email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmpassword")
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		if password != confirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}
		if len(password) < 8 {
			http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		if matched := regexp.MustCompile(regex).MatchString(Email); !matched {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		_, err = database.Db.Exec(`
		INSERT INTO users(email, username, password)
		VALUES(?, ?, ?)`,
			Email, username, hashedPassword,
		)
		if err != nil {
			log.Println("Error inserting user:", err)
			http.Error(w, "Internal server error f", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
