package posts

import (
	"fmt"
)

func verifyCategories(categories []string) (bool, error) {
	if len(categories) == 0 {
		return false, fmt.Errorf("please choose at least one category")
	}

	defaults := []string{"All", "Programming", "Cybersecurity", "Gadgets & Hardware", "Web Development"}

	for _, cat := range categories {
		valid := false
		for _, allowed := range defaults {
			if cat == allowed {
				valid = true
				break
			}
		}
		if !valid {
			return false, fmt.Errorf("invalid category: %s", cat)
		}
	}

	return true, nil
}
