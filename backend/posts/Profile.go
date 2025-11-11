package posts

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	Errorhandel "forum/backend/Errors"
	"forum/backend/auth"
	"forum/backend/database"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	MyProfile, err := template.ParseFiles("./frontend/templates/Profile.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	MyProfileData := database.AllPageData(r, "Profile")
	Username, err := auth.GetUsernameFromCookie(r, "session_token")
	if err != nil && err != sql.ErrNoRows && fmt.Sprintf("%v", err) != "http: named cookie not present" {
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	MyProfileData.Username = Username

	MyProfile.Execute(w, MyProfileData)
}
