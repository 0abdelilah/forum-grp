package databasecreate

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"forum/backend/models"
)

func AllPageData(r *http.Request, handle string) models.PageData {
	postStr := r.URL.Query().Get("postid")
	postId, err := strconv.Atoi(postStr)
	if err != nil && postStr != "" {
		log.Printf("invalid postid: %v", err)
		return models.PageData{}
	}

	r.ParseForm()
	categories := r.Form["categories"]

	username := r.PathValue("username")
	fmt.Println(username)
	switch handle {
	case "HomeData":
		posts := GetPosts(categories)
		return models.PageData{AllPosts: posts}

	case "postContent":
		post := GetPostDetails(postId)
		return models.PageData{PostContent: post}
	case "Profile":
		Profile := GetProfile(username)
		fmt.Println(Profile)
		return models.PageData{Profile: Profile}
	default:
		return models.PageData{}
	}
}
