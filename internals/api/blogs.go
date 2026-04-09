package api

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/DevitoDbug/portfolio/internals/utils"
	"github.com/DevitoDbug/portfolio/internals/web/pages"
	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
)

func (a *Api) BlogsHandler(w http.ResponseWriter, r *http.Request) {
	blogNames := []pages.Blog{}
	rootBlogDir := "./internals/blogs/"

	files, err := os.ReadDir(rootBlogDir)
	if err != nil {
		// TODO: create error page
		fmt.Printf("%v\n", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		blogInfo, err := utils.ReadBlogDetails(rootBlogDir + fileInfo.Name())
		if err != nil {
			continue
		}

		blogNames = append(blogNames, pages.Blog{
			Title:     blogInfo.Title,
			CreatedAt: blogInfo.Date,
			FileName:  fileInfo.Name(),
		})
	}

	_ = pages.BlogsPage(rootBlogDir, blogNames).Render(r.Context(), w)
}

func (a *Api) BlogHandler(w http.ResponseWriter, r *http.Request) {
	blogRoot := "./internals/blogs/"
	fileName := chi.URLParam(r, "name")
	filePath := blogRoot + fileName

	// Checking if blog is in the ./internals/blogs
	file, err := os.ReadFile(filePath)
	if err != nil {
		_ = pages.Error("Could not retrieve blog")
	}

	// Blog info
	blogDetails, err := utils.ReadBlogDetails(filePath)
	if err != nil {
		_ = pages.Error("Could not retrieve blog")
	}

	var blogBody string
	rawFileString := strings.SplitN(string(file), "---", 3)
	if len(rawFileString) == 3 {
		blogBody = rawFileString[2]
	} else {
		blogBody = string(file)
	}

	// Insert title and date as h1 and h2 respectively
	blogBody = fmt.Sprintf("# %v\n## %v\n%v", blogDetails.Title, blogDetails.Date, blogBody)

	var buffer bytes.Buffer
	parser := goldmark.New()
	err = parser.Convert([]byte(blogBody), &buffer)
	if err != nil {
		_ = pages.Error("Could not retrieve blog")
	}

	_ = pages.BlogPage(buffer.String()).Render(r.Context(), w)
}
