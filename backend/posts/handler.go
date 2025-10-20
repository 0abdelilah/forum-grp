package posts

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"forum/backend/database"
)

func LoadPostsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not used anymore")
}

func SeePostdetail(w http.ResponseWriter, r *http.Request) {
	PostsTemplete, err := template.ParseFiles("./frontend/templates/post-detail.html")
	if err != nil {
		log.Fatal(err)
	}
	if r.URL.Path != "/post-detail" {
		http.Error(w, "Page Not Found ", http.StatusNotFound)
		return
	}
	PostDatacontent := database.AllPageData(r, "postContent")
	// fmt.Println("Data from Post",PostDatacontent)
	PostsTemplete.Execute(w, PostDatacontent)
}
