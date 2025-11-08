package posts

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	Errorhandel "forum/backend/Errors"
	"forum/backend/auth"
	"forum/backend/database"
	"forum/backend/home"
)

func SeePostdetail(w http.ResponseWriter, r *http.Request) {
	PostsTemplete, err := template.ParseFiles("./frontend/templates/post-detail.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	PageData := database.AllPageData(r, "postContent")
	if PageData.PostContent.Id == 0 {
		Errorhandel.Errordirect(w, "Page not Found", http.StatusNotFound)
		return
	}

	PageData.Username, _ = auth.GetUsernameFromCookie(r, "session_token")

	PostsTemplete.Execute(w, PageData)
}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	username, err := auth.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Failed to parse form")
		return
	}

	categories := r.Form["categories"]
	title := r.FormValue("title")
	content := r.FormValue("content")

	if len(strings.Trim(title, " ")) < 1 || len(title) > 90 {
		home.HomePageError(w, r, "Title must be between 1 and 90 characters" ,http.StatusBadRequest)
		return
	}

	if len(strings.Trim(content, " ")) < 1 || len(content) > 300 {
		home.HomePageError(w, r, "Content must be between 1 and 300 characters",http.StatusBadRequest)
		return
	}

	// test if working
	if valid, err := verifyCategories(categories); !valid {
		home.HomePageError(w, r, err.Error(),http.StatusBadRequest)
		return
	}

	err = InsertPost(username, strings.Trim(title, " "), strings.Trim(content, " "), categories)
	if err != nil {
		home.HomePageError(w, r, "Internal Server error, try later",http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func PostPageError(w http.ResponseWriter, r *http.Request, Error string) {
	tmpl, err := template.ParseFiles("./frontend/templates/post-detail.html")
	if err != nil {
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get the normal page data
	PageData := database.AllPageData(r, "postContent")

	PageData.Username, _ = auth.GetUsernameFromCookie(r, "session_token")
	// Add error
	PageData.Error = Error

	// Execute template
	if err := tmpl.Execute(w, PageData); err != nil {
		log.Printf("template execution error: %v", err)
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
	}
}

func InsertPost(username, title, content string, categories []string) error {
	fmt.Println(len(title))
	res, err := database.Db.Exec(`
		INSERT INTO posts (username, title, content)
		VALUES (?, ?, ?)
	`, username, title, content)
	if err != nil {
		return err
	}

	postID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// Add all categories
	for _, catName := range categories {
		catID, err := getCategoryID(catName)
		if err != nil {
			fmt.Println("Error getting category ID:", err)
			continue
		}
		fmt.Println("Post ID:", postID, "Category ID:", catID)

		err = addPostCategory(int(postID), catID) // working
		if err != nil {
			fmt.Println("Error adding category to post:", err)
		}
	}

	return nil
}

func PostDelete(w http.ResponseWriter, r *http.Request) {
	username, err := auth.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Failed to parse form")
		return
	}
	PostId, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/Delete/"))
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	parentpath := r.FormValue("path")
	// id:=r.FormValue("id")
	// fmt.Println("the id",id)
	// fmt.Println("parentpath:",parentpath)
	// fmt.Println("the Id", PostId)
	// fmt.Println("The name", username)
	err = Deletepost(PostId)
	if err != nil {
		fmt.Println("there an error:", err)
	}
	fmt.Println(parentpath)
	if parentpath == "PostDeleteDetail" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/Profile/%v", username), http.StatusSeeOther)
}

func Deletepost(Postid int) error {
	_, err := database.Db.Exec("delete from comments where id =?", Postid)
	if err != nil {
		return (err)
	}
	_, err = database.Db.Exec("Delete from  posts WHERE id = ?", Postid)
	if err != nil {
		return err
	}
	_, err = database.Db.Exec("delete from post_categories where post_id = ?", Postid)
	if err != nil {
		return err
	}
	_, err = database.Db.Exec("delete from likes where target_id =?", Postid)
	if err != nil {
		return err
	}
	return nil
}
