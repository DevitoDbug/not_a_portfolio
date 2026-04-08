package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DevitoDbug/portfolio/internals/web/pages"
)

func (a *Api) BlogsHandler(w http.ResponseWriter, r *http.Request) {
	blogNames := []string{}
	files, err := os.ReadDir("./internals/blogs")
	if err != nil {
		// TODO: create error page
		fmt.Printf("%v\n", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		blogNames = append(blogNames, file.Name())
	}

	_ = pages.BlogsPage(blogNames).Render(r.Context(), w)
}

func (a *Api) BlogHandler(w http.ResponseWriter, r *http.Request) {
}
