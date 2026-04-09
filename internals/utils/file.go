package utils

import (
	"os"
	"strings"
)

type BlogDetails struct {
	Title string
	Date  string
}

func ReadBlogDetails(filePath string) (*BlogDetails, error) {
	var blogDetails BlogDetails
	rawContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Line by line content
	content := strings.Split(string(rawContent), "\n")
	inDetailsBlock := false

	for _, line := range content {
		if inDetailsBlock {
			parts := strings.Split(line, ":")
			if len(parts) < 2 {
				continue
			}

			switch parts[0] {
			case "title":
				blogDetails.Title = parts[1]
			case "date":
				blogDetails.Date = parts[1]
			default:
				continue
			}

		}

		if line == "---" && !inDetailsBlock {
			inDetailsBlock = true
		} else if line == "---" && inDetailsBlock {
			break
		}
	}

	return &blogDetails, nil
}
