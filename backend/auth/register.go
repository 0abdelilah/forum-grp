package auth

import (
	"fmt"
	"forum/backend/database"
	"html/template"
	"log"
	"net/http"
	"regexp"

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
	if username == "" {
		tmpt.Execute(w, struct{ Error string }{Error: "Username is required"})
		return
	}

	if password != confirmPassword {
		tmpt.Execute(w, struct{ Error string }{Error: "Passwords do not match"})
		return
	}
	if len(password) < 8 {
		tmpt.Execute(w, struct{ Error string }{Error: "Password must be at least 8 characters long"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tmpt.Execute(w, struct{ Error string }{Error: "Internal server error, try again later"})
		log.Println("Error hashing password:", err)
		return
	}
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched := regexp.MustCompile(regex).MatchString(Email); !matched {
		tmpt.Execute(w, struct{ Error string }{Error: "Invalid email format"})
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
