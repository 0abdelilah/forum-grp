package posts

import (
	"fmt"
	"html/template"
	"net/http"

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
	MyProfile.Execute(w, MyProfileData)
}
