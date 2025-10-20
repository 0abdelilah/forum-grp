package filters

import "forum/backend/models"

// func showfelters()
func FelterbyCategory(m models.PageData, value string) []models.Post {
	var newPosts []models.Post
	for i := 0; i < len(m.AllPosts); i++ {
		if contains(m.AllPosts[i].Categories, value) {
			newPosts = append(newPosts, m.AllPosts[i])
		}
	}
	return newPosts
}
func contains(elems []models.Category, v string) bool {
	for _, s := range elems {
		if v == s.Category {
			return true
		}
	}
	return false
}
