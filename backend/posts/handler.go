package posts

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

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
		home.PageNotFound(w)
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

	if len(title) < 1 || len(content) > 90 {
		home.HomePageError(w, r, "Title must be between 1 and 90 characters")
		return
	}

	if len(content) < 1 || len(content) > 300 {
		home.HomePageError(w, r, "Content must be between 1 and 300 characters")
		return
	}

	// test if working
	if valid, err := verifyCategories(categories); !valid {
		home.HomePageError(w, r, err.Error())
		return
	}

	err = InsertPost(username, title, content, categories)
	if err != nil {
		home.HomePageError(w, r, "Internal Server error, try later")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func InsertPost(username, title, content string, categories []string) error {
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
	parentpath:=r.FormValue("path")
	// id:=r.FormValue("id")
	// fmt.Println("the id",id)
	// fmt.Println("parentpath:",parentpath)
	// fmt.Println("the Id", PostId)
	// fmt.Println("The name", username)
	err = Deletepost(username, PostId)
	if err != nil {
		fmt.Println("there an error:", err)
	}
	if parentpath=="PostDeteleDetail"{
	   http.Redirect(w,r,"/",http.StatusSeeOther)
	}
	http.Redirect(w,r,fmt.Sprintf("/Profile/%v",username),http.StatusSeeOther)
}

func Deletepost(usarname string, Postid int) error {
	_, err := database.Db.Exec("Delete from  posts WHERE id = ? AND username = ?", Postid, usarname)
	if err != nil {
		return err
	}
	return nil
}
