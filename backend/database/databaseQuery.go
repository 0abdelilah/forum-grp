package database

import (
	"forum/backend/models"
	"log"
	"net/http"
	"strconv"
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

	switch handle {
	case "HomeData":
		posts := GetPosts(categories)
		return models.PageData{AllPosts: posts}

	case "postContent":
		post := GetPostDetails(postId)
		return models.PageData{PostContent: post}

	default:
		return models.PageData{}
	}
}
