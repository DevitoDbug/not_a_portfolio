package api

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

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

	// Checking if blog is in the ./internals/blogs
	file, err := os.ReadFile(blogRoot + fileName)
	if err != nil {
		_ = pages.Error("Could not retrieve blog")
	}

	var buffer bytes.Buffer
	parser := goldmark.New()
	err = parser.Convert(file, &buffer)
	if err != nil {
		_ = pages.Error("Could not retrieve blog")
	}

	_ = pages.BlogPage(buffer.String()).Render(r.Context(), w)
}
